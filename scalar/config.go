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
	env := sets.FromSlice(c.Env)

	allResources := sets.FromSlice(c.Resources)
	requiredResources := sets.FromSlice(jobResources)

	envMatch := true
	if jobEnv != "" {
		envMatch = env.Contains(jobEnv)
	}

	return envMatch && allResources.IsSupersetOf(requiredResources)
}

func (c *Config) matchAgent(agentEnv []string, resources []string) bool {
	env := sets.FromSlice(c.Env)
	agentEnvSet := sets.FromSlice(agentEnv)
	allResources := sets.FromSlice(c.Resources)
	agentResource := sets.FromSlice(resources)

	return env.Equal(agentEnvSet) && allResources.Equal(agentResource)
}
