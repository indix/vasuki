package scalar

import "github.com/deckarep/golang-set"

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
	env := mapset.NewSetFromSlice(stringSliceToInterfaceSlice(c.Env))

	allResources := mapset.NewSetFromSlice(stringSliceToInterfaceSlice(c.Resources))
	requiredResources := mapset.NewSetFromSlice(stringSliceToInterfaceSlice(jobResources))

	return env.Contains(jobEnv) && allResources.IsSuperset(requiredResources)
}

func (c *Config) matchAgent(agentEnv []string, resources []string) bool {
	env := mapset.NewSetFromSlice(stringSliceToInterfaceSlice(c.Env))
	agentEnvSet := mapset.NewSetFromSlice(stringSliceToInterfaceSlice(agentEnv))
	allResources := mapset.NewSetFromSlice(stringSliceToInterfaceSlice(c.Resources))
	agentResource := mapset.NewSetFromSlice(stringSliceToInterfaceSlice(resources))

	return env.Equal(agentEnvSet) && allResources.Equal(agentResource)
}
