package scalar

import (
	"fmt"
	"math"

	"github.com/ashwanthkumar/go-gocd"
	"github.com/ashwanthkumar/vasuki/executor"
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
	// ScheduledJobs matching the corresponding Env and Resources
	ScheduledJobs() ([]*gocd.ScheduledJob, error)
	// IdleAgents matching the corresponding Env and Resources
	IdleAgents() ([]*gocd.Agent, error)
}

// SimpleScalar implementation
type SimpleScalar struct {
	_config *Config
	_client *gocd.Client
}

// NewSimpleScalar - Creates a new scalar.SimpleScalar instance
func NewSimpleScalar(env []string, resources []string, client *gocd.Client) (Scalar, error) {
	return &SimpleScalar{
		_config: NewConfig(env, resources),
		_client: client,
	}, nil
}

func (s *SimpleScalar) config() *Config {
	return s._config
}

func (s *SimpleScalar) client() *gocd.Client {
	return s._client
}

// Execute - Entry point of the Scalar
func (s *SimpleScalar) Execute() error {
	var resultErr *multierror.Error
	pendingJobs, err := s.ScheduledJobs() // demand
	idleAgents, err := s.IdleAgents()     // supply - from GoCD Server
	updateErrors(resultErr, err)
	executorReportedAgentIds, err := executor.DefaultExecutor.ManagedAgents() // supply - from Executor instance
	updateErrors(resultErr, err)
	if resultErr.ErrorOrNil() != nil {
		return resultErr.ErrorOrNil()
	}
	var idleAgentIds []string
	for _, idleAgent := range idleAgents {
		idleAgentIds = append(idleAgentIds, idleAgent.UUID)
	}
	// fmt.Printf("Idle Agents =%v\n", idleAgentIds)
	supplyAgents := sets.FromSlice(executorReportedAgentIds).Union(sets.FromSlice(idleAgentIds))
	// fmt.Printf("Supply agents =%v\n", supplyAgents)

	demand := len(pendingJobs)
	supply := supplyAgents.Size()
	if demand > supply {
		diff := demand - supply
		config := s.config()
		instancesToScaleUp := int(math.Ceil(float64(diff) / 2))
		fmt.Printf("Found demand with Env=%v, Resources=%v, scaling up by %d instances.\n", config.Env, config.Resources, instancesToScaleUp)
		err = executor.DefaultExecutor.ScaleUp(instancesToScaleUp)
		updateErrors(resultErr, err)
	} else if supply > demand {
		diff := supply - demand
		config := s.config()
		instancesToScaleDown := int(math.Ceil(float64(diff) / 2))

		if len(idleAgentIds) >= instancesToScaleDown {
			fmt.Printf("Found excess supply for Env=%v, Resources=%v. Idle Agents is %d\n", config.Env, config.Resources, len(idleAgentIds))
			agentsToKill := idleAgentIds[0:instancesToScaleDown]
			err = executor.DefaultExecutor.ScaleDown(agentsToKill)
			for _, agentID := range agentsToKill {
				fmt.Printf("Disabling the agent %s on Go Server\n", agentID)
				err = s.client().DisableAgent(agentID)
				updateErrors(resultErr, err)
				fmt.Printf("Deleting the agent %s on Go Server\n", agentID)
				err = s.client().DeleteAgent(agentID)
				updateErrors(resultErr, err)
			}
			updateErrors(resultErr, err)
		}
	} else {
		fmt.Println("We're in Ideal world. Inner Peace.")
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
			if config.match(job.Environment, job.Resources()) {
				filteredJobs = append(filteredJobs, job)
			}
		}
		return filteredJobs, nil
	}
	return []*gocd.ScheduledJob{}, err
}

// IdleAgents - Get array of Agents that match our environment, resource combination
func (s *SimpleScalar) IdleAgents() ([]*gocd.Agent, error) {
	config := s.config()
	agents, err := s.client().GetAllAgents()
	if err == nil {
		filteredAgents := agents[:0]
		for _, agent := range agents {
			if config.matchAgent(agent.Env, agent.Resources) &&
				(agent.BuildState == "Idle" || agent.BuildState == "Unknown") {
				filteredAgents = append(filteredAgents, agent)
			}
		}

		return filteredAgents, nil
	}
	return []*gocd.Agent{}, err
}

func updateErrors(resultErr *multierror.Error, err error) {
	if err != nil {
		resultErr = multierror.Append(resultErr, err)
	}
}
