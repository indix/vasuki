package docker

import (
	"fmt"
	"strings"

	"github.com/fsouza/go-dockerclient"
	"github.com/hashicorp/go-multierror"
	"github.com/ind9/vasuki/executor"
	"github.com/ind9/vasuki/utils/logging"
	"github.com/satori/go.uuid"
)

// Executor - Docker based Executor implementation
type Executor struct {
	config       *executor.Config
	dockerClient *docker.Client
	dockerImage  string
}

// Init - Initialize this Executor instance
func (e *Executor) Init(config *executor.Config) (err error) {
	e.config = config
	dockerEndpoint := config.Additional["DOCKER_ENDPOINT"]
	e.dockerImage = config.Additional["DOCKER_IMAGE"]
	if config.Additional["DOCKER_FROM_ENV"] == "true" {
		e.dockerClient, err = docker.NewClientFromEnv()
		return err
	}
	e.dockerClient, err = docker.NewClient(dockerEndpoint)
	return err
}

// ScaleUp - Initiate a scaleUp activity among the agents that are managed by this executor instance
func (e *Executor) ScaleUp(instances int) (err error) {
	logging.Log.Infof("Scaling up %d agents via Docker", instances)
	containerLabels := make(map[string]string, 0)
	containerLabels["ENV"] = strings.Join(e.config.Env, ",")
	containerLabels["RESOURCES"] = strings.Join(e.config.Resources, ",")
	containerLabels["VASUKI_MANAGED"] = "true" // watermark to find the containers we spun
	var resultErr *multierror.Error
	for count := 0; count < instances; count++ {
		agentID := uuid.NewV4()
		containerLabels["GO_AGENT_UUID"] = agentID.String()

		config := &docker.Config{
			Image: e.dockerImage,
			Env: []string{
				fmt.Sprintf("GO_SERVER=%s", e.config.ServerHost),
				fmt.Sprintf("GO_SERVER_PORT=%d", e.config.ServerPort),
				fmt.Sprintf("AGENT_ENVIRONMENTS=%s", strings.Join(e.config.Env, ",")),
				fmt.Sprintf("AGENT_RESOURCES=%s", strings.Join(e.config.Resources, ",")),
				fmt.Sprintf("AGENT_KEY=%s", e.config.AutoRegisterKey),
				fmt.Sprintf("AGENT_GUID=%s", agentID.String()),
			},
			Labels: containerLabels,
		}
		opts := docker.CreateContainerOptions{
			Config: config,
		}
		container, err := e.dockerClient.CreateContainer(opts)
		resultErr = updateErrors(resultErr, err)
		if err == nil {
			err = e.dockerClient.StartContainer(container.ID, nil)
			resultErr = updateErrors(resultErr, err)
		}
		logging.Log.Debugf("Started agent container %s", agentID.String())
	}

	return err
}

// ScaleDown - Initiate a scaledown activity among the agents that are managed by this executor instance
func (e *Executor) ScaleDown(agentsToKill []string) (err error) {
	var resultErr *multierror.Error
	for _, agentID := range agentsToKill {
		containerID, err := e.findContainerIDFor(agentID)
		resultErr = updateErrors(resultErr, err)
		if err == nil {
			opts := docker.KillContainerOptions{
				ID: *containerID,
			}
			err := e.dockerClient.KillContainer(opts)
			logging.Log.Infof("Terminating agent %s created via Docker", agentID)
			resultErr = updateErrors(resultErr, err)
		}
	}
	return resultErr.ErrorOrNil()
}

func (e *Executor) findContainerIDFor(agentID string) (*string, error) {
	var resultErr *multierror.Error
	containerFilters := make(map[string][]string)
	containerFilters["label"] = []string{
		"VASUKI_MANAGED=true", // watermark
		fmt.Sprintf("ENV=%s", strings.Join(e.config.Env, ",")),
		fmt.Sprintf("RESOURCES=%s", strings.Join(e.config.Resources, ",")),
		fmt.Sprintf("GO_AGENT_UUID=%s", agentID),
	}
	opts := docker.ListContainersOptions{
		Filters: containerFilters,
	}
	containers, err := e.dockerClient.ListContainers(opts)
	resultErr = updateErrors(resultErr, err)
	if len(containers) < 1 {
		err = fmt.Errorf("Container for agent id=%s not found", agentID)
		resultErr = updateErrors(resultErr, err)
	}

	if resultErr.ErrorOrNil() != nil {
		return nil, resultErr.ErrorOrNil()
	}

	return &(containers[0].ID), nil
}

// ManagedAgents - List of UUIDs of the agents that are managed through this executor instance
func (e *Executor) ManagedAgents() ([]string, error) {
	containerFilters := make(map[string][]string)
	containerFilters["label"] = []string{
		"VASUKI_MANAGED=true", // watermark
		fmt.Sprintf("ENV=%s", strings.Join(e.config.Env, ",")),
		fmt.Sprintf("RESOURCES=%s", strings.Join(e.config.Resources, ",")),
	}
	opts := docker.ListContainersOptions{
		Filters: containerFilters,
	}
	containers, err := e.dockerClient.ListContainers(opts)
	var agentIds []string
	for _, container := range containers {
		agentIds = append(agentIds, container.Labels["GO_AGENT_UUID"])
	}
	return agentIds, err
}

func init() {
	executor.DefaultExecutor = &Executor{}
}

func updateErrors(resultErr *multierror.Error, err error) *multierror.Error {
	if err != nil {
		resultErr = multierror.Append(resultErr, err)
	}

	return resultErr
}
