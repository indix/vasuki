package scalar

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigMatch(t *testing.T) {
	config := Config{
		Env:       []string{"FT", "Staging"},
		Resources: []string{"FT"},
	}

	assert.Equal(t, true, config.match("FT", []string{"FT"}))
	assert.Equal(t, false, config.match("Production", []string{"Production"}))
	assert.Equal(t, false, config.match("FT", []string{"Firefox"}))
}

func TestConfigMatchWhenNoEnvironmentOrResources(t *testing.T) {
	config := Config{
		Env:       []string{},
		Resources: []string{},
	}

	assert.Equal(t, true, config.match("", []string{}))
}

func TestConfigMatchWhenOnlyNoEnvironment(t *testing.T) {
	config := Config{
		Env:       []string{},
		Resources: []string{"packer", "terraform"},
	}

	assert.Equal(t, true, config.match("", []string{"packer"}))
	assert.Equal(t, false, config.match("", []string{"docker"}))
}
