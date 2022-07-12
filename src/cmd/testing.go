package cmd

import (
	"github.com/rafalmnich/edge-iqcontrols-app/internal/config"
)

// ResetContainer resets service container for the testing purposes.
func ResetContainer() {
	services = &serviceContainer{}
}

// SetConfigService allows to inject config service into service container for the testing purposes.
func SetConfigService(cfgSrv *config.Service) {
	services.configService = cfgSrv
}
