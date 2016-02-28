package scalar

import (
	"math"

	"github.com/ashwanthkumar/go-gocd"
	"github.com/ashwanthkumar/vasuki/executor"
	"github.com/ashwanthkumar/vasuki/utils/logging"
	"github.com/ashwanthkumar/vasuki/utils/sets"
	"github.com/hashicorp/go-multierror"
)

// Scalar - Wrapper that decides the agent resources for scaling up or down
type Scalar interface {
	// Instance of the config for the Scalar
	config() *Config
	// Instance of the GoCD Client
	client() *gocd.Client
	// Execute the Scalar
	Execute() error
	// Compute the demand of the GoCD sever
	Demand() (int, error)
	// Compute the supply of agents to GoCD Server
	Supply() (int, error)
}

// SimpleScalar implementation
type SimpleScalar struct {
	_config *Config
	_client *gocd.Client
}

// NewSimpleScalar - Creates a new scalar.SimpleScalar instance
func NewSimpleScalar(env []string, resources []string,
	maxAgents int,
	client *gocd.Client) (Scalar, error) {
	return &SimpleScalar{
		_config: NewConfig(env, resources, maxAgents),
		_client: client,
	}, nil
}

func (s *SimpleScalar) config() *Config {
	return s._config
}

func (s *SimpleScalar) client() *gocd.Client {
	return s._client
}

// Demand in GoCD Server based on ScheduledJobs + Agents that're building
func (s *SimpleScalar) Demand() (int, error) {
	var resultErr *multierror.Error
	pendingJobs, err := s.ScheduledJobs() // demand - from Job Queue
	resultErr = updateErrors(resultErr, err)
	buildingAgents, err := s.BuildingAgents() // demand - from from Agent Queu
	resultErr = updateErrors(resultErr, err)

	demand := len(pendingJobs) + len(buildingAgents)
	return demand, resultErr.ErrorOrNil()
}

// Supply in GoCD Server based on Idle agents + DefaultExecutor's ManagedAgents
func (s *SimpleScalar) Supply() (int, error) {
	var resultErr *multierror.Error
	idleAgentIds, err := s.IdleAgents() // supply - from GoCD Server
	resultErr = updateErrors(resultErr, err)
	executorReportedAgentIds, err := executor.DefaultExecutor.ManagedAgents() // supply - from Executor instance
	resultErr = updateErrors(resultErr, err)
	if resultErr.ErrorOrNil() != nil {
		return 0, resultErr.ErrorOrNil()
	}

	supplyAgents := sets.FromSlice(executorReportedAgentIds).Union(sets.FromSlice(idleAgentIds))

	return supplyAgents.Size(), resultErr.ErrorOrNil()
}

// Execute - Entry point of the Scalar
func (s *SimpleScalar) Execute() error {
	var resultErr *multierror.Error

	demand, err := s.Demand()
	resultErr = updateErrors(resultErr, err)
	supply, err := s.Supply()
	resultErr = updateErrors(resultErr, err)
	if resultErr.ErrorOrNil() != nil {
		return resultErr.ErrorOrNil()
	}

	logging.Log.Debugf("Jobs in Queue (aka) Demand=%d", demand)
	logging.Log.Debugf("Reporting Agents (aka) Supply=%d", supply)
	if demand > supply {
		diff := demand - supply
		config := s.config()
		instancesToScaleUp := int(math.Ceil(float64(diff) / 2))
		if supply >= config.MaxAgents {
			logging.Log.Infof("Found demand with Env=%v, Resources=%v, but we already have %d / %d max agents. Not scaling up.", config.Env, config.Resources, supply, config.MaxAgents)
		} else {
			if supply+instancesToScaleUp > config.MaxAgents {
				instancesToScaleUp = config.MaxAgents - supply
			}
			logging.Log.Infof("Found demand with Env=%v, Resources=%v, scaling up by %d instances.\n", config.Env, config.Resources, instancesToScaleUp)
			err = executor.DefaultExecutor.ScaleUp(instancesToScaleUp)
			resultErr = updateErrors(resultErr, err)
		}
	} else if supply > demand {
		diff := supply - demand
		config := s.config()
		idleAgentIds, err := s.IdleAgents()
		instancesToScaleDown := int(math.Min(math.Ceil(float64(diff)/2), float64(len(idleAgentIds))))

		if len(idleAgentIds) > 0 {
			logging.Log.Infof("Found excess supply for Env=%v, Resources=%v. # of Idle Agents = %d.", config.Env, config.Resources, len(idleAgentIds))
			logging.Log.Infof("# of Agents Scaling down = %d", instancesToScaleDown)
			agentsToKill := idleAgentIds[0:instancesToScaleDown]
			for _, agentID := range agentsToKill {
				logging.Log.Infof("Disabling the agent %s on Go Server\n", agentID)
				err = s.client().DisableAgent(agentID)
				resultErr = updateErrors(resultErr, err)
				logging.Log.Debugf("Checking if the disabled agent %s has started building", agentID)
				agent, err := s.client().GetAgent(agentID)
				resultErr = updateErrors(resultErr, err)
				if agent.BuildState != "Building" {
					logging.Log.Debugf("Disabled agent %s is in %s state so deleting it", agentID, agent.BuildState)
					logging.Log.Infof("Deleting the agent %s on Go Server\n", agentID)
					err = s.client().DeleteAgent(agentID)
					resultErr = updateErrors(resultErr, err)
				} else {
					// Agent has started building after we disabled it, enabling it back
					logging.Log.Noticef("Agent %s has started building after it was disabled, enabling it back", agentID)
					err = s.client().EnableAgent(agentID)
					resultErr = updateErrors(resultErr, err)
				}
			}
			err = executor.DefaultExecutor.ScaleDown(agentsToKill)
			resultErr = updateErrors(resultErr, err)
		} else {
			logging.Log.Infof("All agents are busy. Waiting for them to complete work.")
		}
	} else if supply == 0 && demand == 0 {
		logging.Log.Info("No demand / supply was found.")
	} else {
		// When all the demand is scheduledJobs then upscaled agent's haven't registered with
		// Go server yet. Happens when the time to bootstrap (downloading agent-launcher and agent-plugins)
		// is more than our polling frequency
		logging.Log.Infof("All agents are busy. Waiting for them to complete work.")
	}

	return resultErr.ErrorOrNil()
}

// ScheduledJobs - Get array of ScheduledJob that match our environment, resource combination
func (s *SimpleScalar) ScheduledJobs() ([]*gocd.ScheduledJob, error) {
	config := s.config()
	jobs, err := s.client().GetScheduledJobs()
	if err == nil {
		filteredJobs := jobs[:0]
		for _, job := range jobs {
			if config.matchJob(job.Environment, job.Resources()) {
				filteredJobs = append(filteredJobs, job)
			}
		}
		return filteredJobs, nil
	}
	return []*gocd.ScheduledJob{}, err
}

// IdleAgents - Get array of Agents that match our environment, resource combination
func (s *SimpleScalar) IdleAgents() ([]string, error) {
	config := s.config()
	agents, err := s.client().GetAllAgents()
	var idleAgents []string
	if err == nil {
		for _, agent := range agents {
			if config.matchAgent(agent.Env, agent.Resources) &&
				(agent.BuildState == "Idle" || agent.BuildState == "Unknown") &&
				(agent.AgentState == "Idle" || agent.AgentState == "Unknown") {
				idleAgents = append(idleAgents, agent.UUID)
			}
		}

		return idleAgents, nil
	}
	return idleAgents, err
}

// BuildingAgents - Get array of Agents that match our environment, resource combination
func (s *SimpleScalar) BuildingAgents() ([]*gocd.Agent, error) {
	config := s.config()
	agents, err := s.client().GetAllAgents()
	if err == nil {
		filteredAgents := agents[:0]
		for _, agent := range agents {
			if config.matchAgent(agent.Env, agent.Resources) &&
				(agent.BuildState == "Building") &&
				(agent.AgentState == "Building") {
				filteredAgents = append(filteredAgents, agent)
			}
		}

		return filteredAgents, nil
	}
	return []*gocd.Agent{}, err
}

func updateErrors(resultErr *multierror.Error, err error) *multierror.Error {
	if err != nil {
		resultErr = multierror.Append(resultErr, err)
	}

	return resultErr
}
