package config_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stackify/stackify-go-apm/config"
)

func TestDefaultConfig(t *testing.T) {
	c := config.NewConfig()

	assert.Equal(t, c.ApplicationName, "Go Application")
	assert.Equal(t, c.EnvironmentName, "Production")
	assert.Equal(t, c.Debug, false)
	assert.Equal(t, c.TransportType, "default")
	assert.Equal(t, c.LogPath, "/usr/local/stackify/stackify-python-apm/log/")
	assert.Equal(t, c.LogFileThresholdSize, int64(50000000))
	assert.NotEmpty(t, c.BaseDIR)
	assert.NotEmpty(t, c.HostName)
	assert.NotEmpty(t, c.OSType)
	assert.NotEmpty(t, c.ProcessID)
}

func TestConfigOptions(t *testing.T) {
	c := config.NewConfig(
		config.WithApplicationName("TestName"),
		config.WithEnvironmentName("TestEnv"),
		config.WithDebug(true),
		config.WithLogPath("/"),
		config.WithLogFileThresholdSize(100),
		config.WithTransportType("default"),
	)

	assert.Equal(t, c.ApplicationName, "TestName")
	assert.Equal(t, c.EnvironmentName, "TestEnv")
	assert.Equal(t, c.Debug, true)
	assert.Equal(t, c.TransportType, "default")
	assert.Equal(t, c.LogPath, "/")
	assert.Equal(t, c.LogFileThresholdSize, int64(100))
	assert.NotEmpty(t, c.BaseDIR)
	assert.NotEmpty(t, c.HostName)
	assert.NotEmpty(t, c.OSType)
	assert.NotEmpty(t, c.ProcessID)
}

func TestConfigOptionInvalidTransportType(t *testing.T) {
	c := config.NewConfig(
		config.WithTransportType("invalid"),
	)

	assert.Equal(t, c.TransportType, "default")
}
