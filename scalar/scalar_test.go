package scalar

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	TestEnv       = []string{"Test"}
	TestResources = []string{"TestResource"}
	TestMaxAgents = 3
)

func TestComputeScaleUp(t *testing.T) {
	scalar, err := NewSimpleScalar(TestEnv, TestResources, TestMaxAgents, nil)
	assert.NoError(t, err)

	instances, _ := scalar.ComputeScaleUp(5, 0)
	assert.Equal(t, 3, instances)
	instances, _ = scalar.ComputeScaleUp(5, 5)
	assert.Equal(t, 0, instances)

	instances, _ = scalar.ComputeScaleUp(5, 2)
	assert.Equal(t, 1, instances)
	assert.Equal(t, TestMaxAgents, instances+2)

	instances, _ = scalar.ComputeScaleUp(5, 1)
	assert.Equal(t, 2, instances)
	assert.Equal(t, TestMaxAgents, instances+1)
}

func TestComputeScaleDown(t *testing.T) {
	scalar, err := NewSimpleScalar(TestEnv, TestResources, TestMaxAgents, nil)
	assert.NoError(t, err)

	instances, _ := scalar.ComputeScaleDown(0, 5, 1)
	assert.Equal(t, 1, instances)
	instances, _ = scalar.ComputeScaleDown(0, 5, 3)
	assert.Equal(t, 3, instances)
	instances, _ = scalar.ComputeScaleDown(0, 2, 2)
	assert.Equal(t, 1, instances)
}
