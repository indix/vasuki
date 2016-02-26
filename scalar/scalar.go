package scalar

import (
	"fmt"
	"math"

	"github.com/ashwanthkumar/go-gocd"
	"github.com/ashwanthkumar/vasuki/executor"
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
	multierror.Append(resultErr, err)
	idleAgents, err := s.IdleAgents() // supply
	multierror.Append(resultErr, err)
	if resultErr.ErrorOrNil() != nil {
		return resultErr.ErrorOrNil()
	}

	demand := len(pendingJobs)
	supply := len(idleAgents)
	if demand > supply {
		diff := demand - supply
		config := s.config()
		instancesToScaleUp := int(math.Ceil(float64(diff) / 2))
		fmt.Printf("We need to invoke the Executor#ScaleUp for %d agents with Env=%v, Resources=%v\n", instancesToScaleUp, config.Env, config.Resources)
		fmt.Printf("We need to invoke the Executor#ScaleUp for %d agents with Env=%v, Resources=%v\n", diff, config.Env, config.Resources)
		executor.DefaultExecutor.ScaleUp(instancesToScaleUp)
	} else if supply > demand {
		diff := supply - demand
		config := s.config()
		instancesToScaleDown := int(math.Ceil(float64(diff) / 2))
		fmt.Printf("We need to invoke the Executor#ScaleDown for %d agents with Env=%v, Resources=%v\n", instancesToScaleDown, config.Env, config.Resources)
		fmt.Printf("We need to invoke the Executor#ScaleDown for %d agents with Env=%v, Resources=%v\n", diff, config.Env, config.Resources)
		executor.DefaultExecutor.ScaleDown(instancesToScaleDown)
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
	}
	return []*gocd.Agent{}, err
}
