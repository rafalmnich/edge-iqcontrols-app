package config

import (
	"sync"
	"time"

	"github.com/futurehomeno/cliffhanger/config"
	"github.com/futurehomeno/cliffhanger/storage"
)

type (
	udp struct {
		Port string `json:"port"`
	}

	mass struct {
		Lights  string `json:"lights"`
		Heating string `json:"heating"`
	}

	// Device represents a device configuration for mapping and value conversion.
	Device struct {
		Name         string                 `json:"name"`
		Address      int64                  `json:"address"`
		ServiceName  string                 `json:"serviceName"`
		Config       map[string]interface{} `json:"config"`
		VariableName string                 `json:"variableName"`
		MsgType      string                 `json:"msgType"`
		ValueType    string                 `json:"valueType"`
		Mass         string                 `json:"mass"`
	}
)

// Config is a model containing all application configuration settings.
type Config struct {
	config.Default

	// User configurable settings
	UDP     udp      `json:"udp"`
	Mass    mass     `json:"mass"`
	Devices []Device `json:"devices"`
}

// New creates new instance of a configuration object.
func New(workDir string) *Config {
	return &Config{
		Default: config.NewDefault(workDir),
	}
}

// NewConfigService creates a new configuration service.
func NewConfigService(workDir string) *Service {
	return &Service{
		Storage: config.NewStorage(New(workDir), workDir),
		lock:    &sync.RWMutex{},
	}
}

// Service is a configuration service responsible for:
// - providing concurrency safe access to settings,
// - persistence of settings.
type Service struct {
	storage.Storage
	lock *sync.RWMutex
}

// GetLogLevel allows to safely access a configuration setting.
func (cs *Service) GetLogLevel() string {
	cs.lock.RLock()
	defer cs.lock.RUnlock()

	return cs.Storage.Model().(*Config).LogLevel
}

// SetLogLevel allows to safely set and persist a configuration setting.
func (cs *Service) SetLogLevel(value string) error {
	cs.lock.Lock()
	defer cs.lock.Unlock()

	cs.Storage.Model().(*Config).ConfiguredAt = time.Now().Format(time.RFC3339)
	cs.Storage.Model().(*Config).LogLevel = value

	return cs.Storage.Save()
}

// GetUDPPort allows to safely access udp port from configuration.
func (cs *Service) GetUDPPort() string {
	cs.lock.RLock()
	defer cs.lock.RUnlock()

	return cs.Storage.Model().(*Config).UDP.Port
}

// SetUDPPort allows to safely set and persist udp port in configuration file.
func (cs *Service) SetUDPPort(value string) error {
	cs.lock.Lock()
	defer cs.lock.Unlock()

	cs.Storage.Model().(*Config).UDP.Port = value

	return cs.Storage.Save()
}

// GetHTTPHost allows to safely access mass host from configuration.
func (cs *Service) GetHTTPHost() string {
	cs.lock.RLock()
	defer cs.lock.RUnlock()

	return cs.Storage.Model().(*Config).Mass.Lights
}

// SetHTTPHost allows to safely set and persist mass host in configuration file.
func (cs *Service) SetHTTPHost(value string) error {
	cs.lock.Lock()
	defer cs.lock.Unlock()

	cs.Storage.Model().(*Config).Mass.Lights = value

	return cs.Storage.Save()
}

// GetDevices allows to safely access devices from configuration.
func (cs *Service) GetDevices() []Device {
	cs.lock.RLock()
	defer cs.lock.RUnlock()

	return cs.Storage.Model().(*Config).Devices
}

// SetDevices allows to safely set and persist devices in configuration file.
func (cs *Service) SetDevices(value []Device) error {
	cs.lock.Lock()
	defer cs.lock.Unlock()

	cs.Storage.Model().(*Config).Devices = value

	return cs.Storage.Save()
}

// Factory is a factory method returning the configuration object without default settings.
func Factory() interface{} {
	return &Config{}
}
