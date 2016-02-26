# Vasuki

Vasuki is a [GoCD](http://go.cd/) Agent autoscalar. It uses docker to bring up agents on demand and scale them back down. In GoCD you generally specify Environments and Resource tags to agents and GoCD Server then takes care of assigning each job to the agents that match the given criteria. With the advent of docker all these can be heavily simplified.

You just have to launch a Vasuki instance for an environment and resources with a docker image. It would periodically poll the GoCD Server for any active jobs waiting in Queue for these resources. If found, it would bring up agents (docker containers) matching the environment and resources.

You can run as many Vasuki instances in a machine for various environments / resources.

## Usage
```bash
$ vasuki \
  --server-host build.indix.tv \
  --server-port 8080 \
  --server-username admin \
  --server-password badger \
  --agent-auto-register-key "123456ABCDEFG" \
  --agent-env FT \
  --agent-resources FT \
  --agent-resources chrome \
  --agent-resources selenium
```

## Command line parameters
```bash
$ vasuki --help
Scale GoCD Agents on demand

Usage:
  vasuki [flags]

Flags:
      --agent-auto-register-key string   AutoRegisterKey for the agent to register to the GoCD Server (default "123456ABCDEFG")
      --agent-env value                  List of environments for the go-agent (default [])
      --agent-max-count int              Maximum number of agents managed by this Vasuki instance (default 1)
      --agent-resources value            List of resources for the go-agent (default [])
      --docker-endpoint string           Docker endpoint to connect to (default "unix:///var/run/docker.sock")
      --docker-image string              Docker image used for spinning up the agent (default "ashwanthkumar/gocd-agent")
      --server-host string               Go Server Domain / IP Address (default "localhost")
      --server-password string           Password of the User to connect to Go Server (default "badger")
      --server-port int                  Go Server Port (default 8153)
      --server-username string           Username to connect to Go Server (default "admin")
```
