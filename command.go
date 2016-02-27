package main

import (
	"fmt"
	"os"
	"time"

	"github.com/ashwanthkumar/go-gocd"
	"github.com/ashwanthkumar/vasuki/executor"
	_ "github.com/ashwanthkumar/vasuki/executor/docker"
	"github.com/ashwanthkumar/vasuki/scalar"
	"github.com/ashwanthkumar/vasuki/utils/logging"
	"github.com/spf13/cobra"
)

var goServerPort int
var goServerHost string
var pollInterval time.Duration
var env []string
var resources []string
var maxAgents int
var autoRegisterKey string
var username string
var password string

// docker settings
var dockerImage string
var dockerEndpoint string
var dockerSettingsFromEnv bool

// misc
var verboseMode bool

var vasukiCommand = &cobra.Command{
	Use:   "vasuki",
	Short: "Scale GoCD Agents on demand",
	Long:  `Scale GoCD Agents on demand`,
	Run: func(cmd *cobra.Command, args []string) {
		logging.EnableDebug(verboseMode)
		logging.Log.Infof("Starting Vasuki instance with Env=%v, Resources=%v", env, resources)
		ServerHost := fmt.Sprintf("http://%s:%d", goServerHost, goServerPort)

		executorAdditionalConfig := make(map[string]string)
		executorAdditionalConfig["DOCKER_IMAGE"] = dockerImage
		executorAdditionalConfig["DOCKER_ENDPOINT"] = dockerEndpoint
		executorAdditionalConfig["DOCKER_FROM_ENV"] = fmt.Sprintf("%t", dockerSettingsFromEnv)

		executorConfig := &executor.Config{
			ServerHost:      goServerHost,
			ServerPort:      goServerPort,
			AutoRegisterKey: autoRegisterKey,
			Env:             env,
			Resources:       resources,
			Additional:      executorAdditionalConfig,
		}
		executor.DefaultExecutor.Init(executorConfig)

		scalar, err := scalar.NewSimpleScalar(env, resources, gocd.New(ServerHost, username, password))
		if err != nil {
			handleError(cmd, err)
		}

		doWork(scalar, cmd)
		c := time.Tick(pollInterval)
		for {
			select {
			case <-c:
				doWork(scalar, cmd)
			}
		}
	},
}

func doWork(scalar scalar.Scalar, cmd *cobra.Command) {
	err := scalar.Execute()
	handleError(cmd, err)
}

func handleError(cmd *cobra.Command, err error) {
	if err != nil {
		logging.Log.Criticalf("[Error] %s", err.Error())
		os.Exit(1)
	}
}

func init() {
	// GoCD Agent related flags
	vasukiCommand.PersistentFlags().StringSliceVar(&env, "agent-env", []string{}, "List of environments for the go-agent")
	vasukiCommand.PersistentFlags().StringSliceVar(&resources, "agent-resources", []string{}, "List of resources for the go-agent")
	vasukiCommand.PersistentFlags().IntVar(&maxAgents, "agent-max-count", 1, "Maximum number of agents managed by this Vasuki instance")
	vasukiCommand.PersistentFlags().StringVar(&autoRegisterKey, "agent-auto-register-key", "123456ABCDEFG", "AutoRegisterKey for the agent to register to the GoCD Server")

	// GoCD Server related flags
	vasukiCommand.PersistentFlags().StringVar(&goServerHost, "server-host", "localhost", "Go Server Domain / IP Address")
	vasukiCommand.PersistentFlags().IntVar(&goServerPort, "server-port", 8153, "Go Server Port")
	vasukiCommand.PersistentFlags().StringVar(&username, "server-username", "", "Username to connect to Go Server")
	vasukiCommand.PersistentFlags().StringVar(&password, "server-password", "", "Password of the User to connect to Go Server")
	vasukiCommand.PersistentFlags().DurationVar(&pollInterval, "server-poll-interval", 30*time.Second, "Poll interval for new scheduled jobs")

	// docker related flags
	vasukiCommand.PersistentFlags().StringVar(&dockerImage, "docker-image", "ashwanthkumar/gocd-agent", "Docker image used for spinning up the agent")
	vasukiCommand.PersistentFlags().StringVar(&dockerEndpoint, "docker-endpoint", "unix:///var/run/docker.sock", "Docker endpoint to connect to")
	vasukiCommand.PersistentFlags().BoolVar(&dockerSettingsFromEnv, "docker-env", false, "Flag to pick up docker settings from Env. Useful when working with boot2docker / docker-machine")

	// misc flags
	vasukiCommand.PersistentFlags().BoolVar(&verboseMode, "verbose", false, "Enable verbose logging")
}
