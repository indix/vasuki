package scalar

import (
	"testing"

	"github.com/ashwanthkumar/go-gocd"
	gocdmocks "github.com/ashwanthkumar/go-gocd/mocks"
	"github.com/ind9/vasuki/executor"
	"github.com/ind9/vasuki/utils/logging"
	"github.com/stretchr/testify/assert"
)

var (
	TestEnv       = []string{"Test"}
	TestResources = []string{"TestResource"}
	TestMaxAgents = 3
)

func init() {
	logging.MuteLogs()
}

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

func TestExecuteForScaleUp(t *testing.T) {
	mockExecutor := new(executor.MockExecutor)
	executor.DefaultExecutor = mockExecutor
	mockExecutor.On("ScaleUp", 1).Return(nil)

	config := NewConfig([]string{"Test-Env"}, []string{"Test-Resource"}, 1)
	scalar := new(MockScalar)
	scalar.On("config").Return(config)
	scalar.On("Demand").Return(1, nil)
	scalar.On("Supply").Return(0, nil)
	scalar.On("ComputeScaleUp", 1, 0).Return(1, nil)

	Execute(scalar)

	scalar.AssertExpectations(t)
	mockExecutor.AssertExpectations(t)
}

func TestExecuteForScaleDown(t *testing.T) {
	client := new(gocdmocks.Client)
	client.On("DisableAgent", "kill-agent-id").Return(nil, nil)
	client.On("DeleteAgent", "kill-agent-id").Return(nil, nil)
	agentStatus := &gocd.Agent{
		BuildState: "Idle",
	}
	client.On("GetAgent", "kill-agent-id").Return(agentStatus, nil)
	mockExecutor := new(executor.MockExecutor)
	executor.DefaultExecutor = mockExecutor
	mockExecutor.On("ScaleDown", []string{"kill-agent-id"}).Return(nil)

	config := NewConfig([]string{"Test-Env"}, []string{"Test-Resource"}, 1)
	scalar := new(MockScalar)
	scalar.On("config").Return(config)
	scalar.On("client").Return(client)
	scalar.On("Demand").Return(0, nil)
	scalar.On("Supply").Return(1, nil)
	scalar.On("ComputeScaleDown", 0, 1, 1).Return(1, nil)
	scalar.On("IdleAgents").Return([]string{"kill-agent-id"}, nil)

	Execute(scalar)

	client.AssertExpectations(t)
	scalar.AssertExpectations(t)
	mockExecutor.AssertExpectations(t)
}

func TestSkipScaleDownIfAgentFoundBuildingAfterDisabling(t *testing.T) {
	client := new(gocdmocks.Client)
	client.On("DisableAgent", "kill-agent-id").Return(nil, nil)
	client.On("EnableAgent", "kill-agent-id").Return(nil, nil)
	agentStatus := &gocd.Agent{
		BuildState: "Building",
	}
	client.On("GetAgent", "kill-agent-id").Return(agentStatus, nil)
	mockExecutor := new(executor.MockExecutor)
	executor.DefaultExecutor = mockExecutor

	config := NewConfig([]string{"Test-Env"}, []string{"Test-Resource"}, 1)
	scalar := new(MockScalar)
	scalar.On("config").Return(config)
	scalar.On("client").Return(client)
	scalar.On("Demand").Return(0, nil)
	scalar.On("Supply").Return(1, nil)
	scalar.On("ComputeScaleDown", 0, 1, 1).Return(1, nil)
	scalar.On("IdleAgents").Return([]string{"kill-agent-id"}, nil)

	Execute(scalar)

	client.AssertExpectations(t)
	scalar.AssertExpectations(t)
	mockExecutor.AssertExpectations(t)
	mockExecutor.AssertNotCalled(t, "ScaleDown", []string{"kill-agent-id"})
	client.AssertNotCalled(t, "DeleteAgent", "kill-agent-id")
}
