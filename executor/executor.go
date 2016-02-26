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
	ScaleDown(instances int) error
}

var DefaultExecutor Executor
