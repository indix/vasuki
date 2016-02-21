package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ashwanthkumar/go-gocd"
	"github.com/ashwanthkumar/vasuki/scalar"
	"github.com/spf13/cobra"
)

var goServerPort int
var goServerHost string
var env []string
var resources []string
var maxAgents int
var autoRegisterKey string
var username string
var password string
var dockerImage string
var dockerEndpoint string

var vasukiCommand = &cobra.Command{
	Use:   "vasuki",
	Short: "Scale GoCD Agents on demand",
	Long:  `Scale GoCD Agents on demand`,
	Run: func(cmd *cobra.Command, args []string) {
		ServerHost := fmt.Sprintf("http://%s:%d", goServerHost, goServerPort)
		scalar := scalar.NewSimpleScalar(env, resources, gocd.New(ServerHost, username, password))
		doWork(scalar, cmd)
		c := time.Tick(30 * time.Second)
		for {
			select {
			case <-c:
				doWork(scalar, cmd)
			}
		}
	},
}

func doWork(scalar *scalar.SimpleScalar, cmd *cobra.Command) {
	err := scalar.Execute()
	if err != nil {
		log.Printf("[Error] %s", err.Error())
		cmd.Help()
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
	vasukiCommand.PersistentFlags().StringVar(&username, "server-username", "admin", "Username to connect to Go Server")
	vasukiCommand.PersistentFlags().StringVar(&password, "server-password", "badger", "Password of the User to connect to Go Server")

	// docker related flags
	vasukiCommand.PersistentFlags().StringVar(&dockerImage, "docker-image", "travix/go-agent", "Docker image used for spinning up the agent")
	vasukiCommand.PersistentFlags().StringVar(&dockerEndpoint, "docker-endpoint", "unix:///var/run/docker.sock", "Docker endpoint to connect to")
}
