package app

import (
	"fmt"
	"strings"

	"github.com/futurehomeno/cliffhanger/app"
	"github.com/futurehomeno/cliffhanger/lifecycle"
	"github.com/futurehomeno/cliffhanger/manifest"
	log "github.com/sirupsen/logrus"

	"github.com/rafalmnich/edge-iqcontrols-app/internal/config"
	"github.com/rafalmnich/edge-iqcontrols-app/internal/reporter"
)

// Application represents the main application.
type Application interface {
	app.App
	app.InitializableApp
}

// application is a private implementation of the main application service.
type application struct {
	cfgSrv         *config.Service
	appLifecycle   *lifecycle.Lifecycle
	manifestLoader manifest.Loader
	deviceReporter reporter.Device
}

// New creates new instance of an application.
func New(
	cfgSrv *config.Service,
	appLifecycle *lifecycle.Lifecycle,
	manifestLoader manifest.Loader,
	deviceReporter reporter.Device,
) Application {
	return &application{
		cfgSrv:         cfgSrv,
		appLifecycle:   appLifecycle,
		manifestLoader: manifestLoader,
		deviceReporter: deviceReporter,
	}
}

// Initialize performs all required initialization before starting the application.
// todo: It should make inclusion reports for not stored devices.
func (a application) Initialize() error {
	for _, d := range a.cfgSrv.GetDevices() {
		if !strings.Contains(d.ServiceName, "sensor") && d.VariableName == "" {
			log.Errorf("application: device %s has no variable name, device is NOT included", d.Name)

			continue
		}

		// todo: make inclusion report for device later with checking for existence.
		// err := a.deviceReporter.InclusionReport(d)
		// if err != nil {
		// 	log.WithError(err).Errorf("application: failed to report inclusion for device %s", d.Name)
		//
		// 	return err
		// }
	}

	return nil
}

// GetManifest returns the manifest object based on current application state and configuration.
func (a application) GetManifest() (*manifest.Manifest, error) {
	appManifest, err := a.manifestLoader.Load()
	if err != nil {
		log.WithError(err).Error("application: failed to load the template")

		return nil, fmt.Errorf("failed to load the template")
	}

	// TODO: You may want to manipulate the manifest depending on current application state or available configuration.
	//  Good examples include modifying list of available devices or dynamic options based on API calls and application lifecycle.

	return appManifest, nil
}

// Configure performs update of the application state based on the provided configuration.
func (a application) Configure(model interface{}) error {
	cfg, ok := model.(*config.Config)
	if !ok {
		log.Errorf("application: invalid config received, should be of %T type, received %T instead", cfg, model)

		return fmt.Errorf("received an invalid configuration")
	}

	// TODO: You may want persist here specific configuration settings provided by the user or act upon them.
	//  Good examples include adding or removing devices from an adapter.

	return nil
}

// Uninstall performs all required cleaning up before uninstalling the application.
func (a application) Uninstall() error {
	if err := a.cfgSrv.Reset(); err != nil {
		log.WithError(err).Errorf("application: failed to reset configuration")

		return fmt.Errorf("failed to reset configuration")
	}

	a.appLifecycle.SetAppState(lifecycle.AppStateNotConfigured, nil)
	a.appLifecycle.SetConfigState(lifecycle.ConfigStateNotConfigured)
	a.appLifecycle.SetConnectionState(lifecycle.ConnStateDisconnected)
	a.appLifecycle.SetAuthState(lifecycle.AuthStateNotAuthenticated)

	return nil
}
