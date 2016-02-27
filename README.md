# Vasuki

Vasuki is a [GoCD](http://go.cd/) agent autoscaler. It uses docker to bring up agents on demand and scale them back down. In GoCD you generally specify Environments and Resource tags to agents and GoCD Server then takes care of assigning each job to the agents that match the given criteria. With the advent of docker all these can be heavily simplified.

You just have to launch a Vasuki instance for an environment and resources with a docker image. It would periodically poll the GoCD Server for any active jobs waiting in Queue with these constraints. If found, it would bring up agents (docker containers) matching the environment and resources. Once found idle, it would kill the container and remove the agent from the server.

You can run as many Vasuki instances in a machine for various environments / resources.

## Features
- Auto scale environments with respective resources only on demand
- Completely stateless, preferred deployment is to start it as a deamon and put it behind a monit like process watch.

## Usage
```bash
$ vasuki \
  --server-host localhost \
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
```
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
      --docker-env                       Flag to pick up docker settings from Env. Useful when working with boot2docker / docker-machine
      --docker-image string              Docker image used for spinning up the agent (default "ashwanthkumar/gocd-agent")
      --server-host string               Go Server Domain / IP Address (default "localhost")
      --server-password string           Password of the User to connect to Go Server
      --server-poll-interval duration    Poll interval for new scheduled jobs (default 30s)
      --server-port int                  Go Server Port (default 8153)
      --server-username string           Username to connect to Go Server
      --verbose                          Enable verbose logging
```

## How does Vasuki work?
1. Query for [active](https://api.go.cd/current/#get-all-agents) + [queued](https://api.go.cd/current/#get-scheduled-jobs) builds. This is Demand.
2. Query for all active agents + list of containers managed by the executor implementation. We then take a union of both. This is Supply.
3. If Demand > Supply, do scale up using the executor implementation
4. If Demand < Supply, do scale down using the executor implementation

## Known Issue
When scaling down, there might a job which is stuck because it got assigned to an agent but Vasuki had just deleted it. This is because Vasuku doesn't get the latest status from [Agents Endpoint](https://api.go.cd/current/#get-all-agents) even after the agent is assigned a job. More details can be found on this [GoCD Dev mail list](https://groups.google.com/d/msg/go-cd-dev/tWmV0Rw9sJM/cz_qe4LcAQAJ) message.

Solution - Decrease the value of `cruise.reschedule.hung.builds.interval` property in GoCD server. This enables faster detection of hung jobs and reschedules them. Start your GoCD server with the following environment variable `GO_SERVER_SYSTEM_PROPERTIES="-Dcruise.reschedule.hung.builds.interval=60000"` (60 seconds).

`cruise.reschedule.hung.builds.interval` by default is [5 minutes](https://github.com/gocd/gocd/blob/master/server/properties/src/cruise.properties#L26).

## FAQs
### Why my tasks take long to start now?
Vasuki polls your GoCD server to find active jobs in queue matching these resources. Hence it's a factor of `--server-poll-interval` flag that you pass. Remember, if you choose a very low value, it'll create unnecessarily load on the server instance.

### For multiple environments and resources should I launch multiple Vasuki instances?
Yes and No. A single Vasuki can manage multiple environments and resources, but currently only 1 docker image per instance. So if you've 2 environments (say FT and UAT) and both are identical with respect to Go Agent, then yes a single Vasuki instance would be fine. Else you might want to launch a separate instance for each environment.
