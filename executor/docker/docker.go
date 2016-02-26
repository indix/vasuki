package docker

import (
	"fmt"

	"github.com/ashwanthkumar/vasuki/executor"
	dockerclient "github.com/fsouza/go-dockerclient"
)

// Executor - Docker based Executor implementation
type Executor struct {
	config       *executor.Config
	dockerClient *dockerclient.Client
	dockerImage  string
}

func (e *Executor) Init(config *executor.Config) (err error) {
	e.config = config
	dockerEndpoint := config.Additional["DOCKER_ENDPOINT"]
	e.dockerImage = config.Additional["DOCKER_IMAGE"]
	e.dockerClient, err = dockerclient.NewClient(dockerEndpoint)
	return err
}

func (e *Executor) ScaleUp(instances int) (err error) {
	// TODO - Implement scaling up
	fmt.Printf("Scaling up %d agents via Docker\n", instances)
	return err
}

func (e *Executor) ScaleDown(instances int) (err error) {
	// TODO - Implement scaling down
	fmt.Printf("Scaling down %d agents created via Docker", instances)
	return err
}

func init() {
	executor.DefaultExecutor = &Executor{}
}
