package scalar

import "github.com/ashwanthkumar/vasuki/utils/sets"

// Config - Holds scalar configurations
type Config struct {
	Env       []string
	Resources []string
}

// NewConfig - Creates a new scalar.Config instance
func NewConfig(env []string, resources []string) *Config {
	return &Config{
		Env:       env,
		Resources: resources,
	}
}

func (c *Config) match(jobEnv string, jobResources []string) bool {
	vasukiEnv := sets.FromSlice(c.Env)
	vasukiResource := sets.FromSlice(c.Resources)

	requiredResources := sets.FromSlice(jobResources)

	var envMatch bool

	if jobEnv == "" {
		// handle no environment job in a special way
		envMatch = vasukiEnv.Size() == 0
	} else {
		envMatch = vasukiEnv.Contains(jobEnv)
	}

	return envMatch && vasukiResource.IsSupersetOf(requiredResources)
}

func (c *Config) matchAgent(agentEnv []string, agentResource []string) bool {
	vasukiEnv := sets.FromSlice(c.Env)
	agentEnvSet := sets.FromSlice(agentEnv)
	vasukiResource := sets.FromSlice(c.Resources)
	agentResourceSet := sets.FromSlice(agentResource)

	return vasukiEnv.Equal(agentEnvSet) && vasukiResource.Equal(agentResourceSet)
}
