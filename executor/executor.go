package executor

// Config - Static Configuration for Executor
type Config struct {
	ServerHost      string
	ServerPort      int
	AutoRegisterKey string
	Env             []string
	Resources       []string
	Additional      map[string]string // Additional configurations for Executor implementation
}

// Executor -
type Executor interface {
	Init(config *Config) error
	ScaleUp(instances int) error
	// currentAgents - List of UUIDs of the *idle agents* we know as of now
	ScaleDown(agentsToKill []string) error
	// Agents that are managed by this Executor instance.
	ManagedAgents() ([]string, error)
}

var Executors map[string]Executor

func init() {
	Executors = make(map[string]Executor)
}
