package config

import (
	"os"
	"reflect"
	"runtime"
	"strconv"
)

var (
	defaults = map[string]map[string]string{
		"ApplicationName": map[string]string{
			"Env":     "STACKIFY_APPLICATION_NAME",
			"Default": "Go Application",
			"Type":    "string",
		},
		"EnvironmentName": map[string]string{
			"Env":     "STACKIFY_ENVIRONMENT_NAME",
			"Default": "Production",
			"Type":    "string",
		},
		"Debug": map[string]string{
			"Env":     "STACKIFY_DEBUG",
			"Default": strconv.FormatBool(false),
			"Type":    "bool",
		},
		"TransportType": map[string]string{
			"Env":     "STACKIFY_TRANSPORT",
			"Default": DefaultTransportType,
			"Type":    "string",
		},
		"LogPath": map[string]string{
			"Env":     "STACKIFY_TRANSPORT_LOG_PATH",
			"Default": "/usr/local/stackify/stackify-python-apm/log/",
			"Type":    "string",
		},
		"LogFileThresholdSize": map[string]string{
			"Env":     "STACKIFY_TRANSPORT_LOG_THRESHOLD_SIZE",
			"Default": strconv.FormatInt(DefaultLogFileThresholdSize, 10),
			"Type":    "int64",
		},
	}
)

type Config struct {
	ApplicationName      string
	EnvironmentName      string
	Debug                bool
	BaseDIR              string
	HostName             string
	OSType               string
	ProcessID            string
	TransportType        string
	LogPath              string
	LogFileThresholdSize int64
}

// Set Default Config Value or Environment Variable if available
func (c *Config) setConfigEnvironmentOrDefault() {
	var tempVal string
	for k, v := range defaults {
		if len(v["Env"]) > 0 {
			tempVal = os.Getenv(v["Env"])
		}
		if len(tempVal) == 0 {
			tempVal = v["Default"]
		}

		if v["Type"] == "string" {
			reflect.ValueOf(c).Elem().FieldByName(k).SetString(tempVal)
		} else if v["Type"] == "bool" {
			val, _ := strconv.ParseBool(tempVal)
			reflect.ValueOf(c).Elem().FieldByName(k).SetBool(val)
		} else if v["Type"] == "int64" {
			val, _ := strconv.ParseInt(tempVal, 10, 64)
			reflect.ValueOf(c).Elem().FieldByName(k).SetInt(val)
		}
	}
}

// Initialize and return config
func NewConfig(opts ...ConfigOptions) *Config {
	config := new(Config)
	config.setConfigEnvironmentOrDefault()

	// set with options
	for _, option := range opts {
		option.Apply(config)
	}

	// set working environments
	config.BaseDIR, _ = os.Getwd()
	config.HostName, _ = os.Hostname()
	config.OSType = runtime.GOOS
	config.ProcessID = strconv.Itoa(os.Getpid())

	return config
}

// ConfigOptions interface for configurable values
type ConfigOptions interface {
	Apply(*Config)
}

type applicationName string

func (o applicationName) Apply(config *Config) {
	config.ApplicationName = string(o)
}

func WithApplicationName(appName string) ConfigOptions {
	return applicationName(appName)
}

type environmentName string

func (o environmentName) Apply(config *Config) {
	config.EnvironmentName = string(o)
}

func WithEnvironmentName(envName string) ConfigOptions {
	return environmentName(envName)
}

type debug bool

func (d debug) Apply(config *Config) {
	config.Debug = bool(d)
}

func WithDebug(d bool) ConfigOptions {
	return debug(d)
}

type transportType string

func (tt transportType) Apply(config *Config) {
	config.TransportType = string(tt)
}

func WithTransportType(tt string) ConfigOptions {
	_, ok := TransportTypes[tt]
	if !ok {
		return transportType(DefaultTransportType)
	}
	return transportType(tt)
}

type logPath string

func (l logPath) Apply(config *Config) {
	config.LogPath = string(l)
}

func WithLogPath(l string) ConfigOptions {
	return logPath(l)
}

type logFileThresholdSize int64

func (l logFileThresholdSize) Apply(config *Config) {
	config.LogFileThresholdSize = int64(l)
}

func WithLogFileThresholdSize(l int64) ConfigOptions {
	return logFileThresholdSize(l)
}
