package cmd

import (
	"net"
	"net/http"
	"time"

	cliffApp "github.com/futurehomeno/cliffhanger/app"
	"github.com/futurehomeno/cliffhanger/bootstrap"
	"github.com/futurehomeno/cliffhanger/hub"
	"github.com/futurehomeno/cliffhanger/lifecycle"
	"github.com/futurehomeno/cliffhanger/manifest"
	"github.com/futurehomeno/cliffhanger/router"
	"github.com/futurehomeno/cliffhanger/task"
	"github.com/futurehomeno/fimpgo"
	log "github.com/sirupsen/logrus"

	"github.com/rafalmnich/edge-iqcontrols-app/internal/app"
	"github.com/rafalmnich/edge-iqcontrols-app/internal/config"
	"github.com/rafalmnich/edge-iqcontrols-app/internal/decoder"
	"github.com/rafalmnich/edge-iqcontrols-app/internal/listener"
	"github.com/rafalmnich/edge-iqcontrols-app/internal/reporter"
	"github.com/rafalmnich/edge-iqcontrols-app/internal/routing"
	"github.com/rafalmnich/edge-iqcontrols-app/internal/tasks"
	"github.com/rafalmnich/edge-iqcontrols-app/internal/transformer"
)

// services is a container for services that are common dependencies.
var services = &serviceContainer{}

// serviceContainer is a type representing a dependency injection container to be used during bootstrap of the application.
type serviceContainer struct {
	configService *config.Service
	hubInfo       *hub.Info
	mqtt          *fimpgo.MqttTransport
	lifecycle     *lifecycle.Lifecycle

	application         cliffApp.App
	deviceReporter      reporter.Device
	configurationLocker router.MessageHandlerLocker
	manifestLoader      manifest.Loader
	// TODO: You may add any additional dependency that has to be injected, e.g.: API client.
}

// getConfigService initiates a configuration service and loads the config.
func getConfigService() *config.Service {
	if services.configService == nil {
		services.configService = config.NewConfigService(
			bootstrap.GetWorkingDirectory(),
		)

		err := services.configService.Load()
		if err != nil {
			log.WithError(err).Fatal("failed to load configuration")
		}
	}

	return services.configService
}

// getInfo retrieves hub info.
// TODO: You may remove this method if you do not need hub information in your application.
func getInfo(cfg *config.Config) *hub.Info {
	if services.hubInfo == nil {
		var err error

		services.hubInfo, err = hub.LoadInfo(cfg.InfoFile)
		if err != nil {
			log.WithError(err).Fatal("failed to load hub info")
		}
	}

	return services.hubInfo
}

// getMQTT creates or returns existing MQTT broker service.
func getMQTT(cfg *config.Config) *fimpgo.MqttTransport {
	if services.mqtt == nil {
		services.mqtt = fimpgo.NewMqttTransport(cfg.MQTTServerURI, cfg.MQTTClientIDPrefix, cfg.MQTTUsername, cfg.MQTTPassword, true, 1, 1)
		services.mqtt.SetDefaultSource(routing.ResourceName)
		services.mqtt.SetAutoDecompression(true)
	}

	return services.mqtt
}

// getLifecycle creates or returns existing lifecycle service.
func getLifecycle(_ *config.Config) *lifecycle.Lifecycle {
	if services.lifecycle == nil {
		services.lifecycle = lifecycle.New()
	}

	return services.lifecycle
}

// getManifestLoader creates or returns existing manifest loader service.
func getManifestLoader(cfg *config.Config) manifest.Loader {
	if services.manifestLoader == nil {
		services.manifestLoader = manifest.NewLoader(cfg.WorkDir)
	}

	return services.manifestLoader
}

// getConfigurationLocker creates or returns existing configuration locker.
func getConfigurationLocker(_ *config.Config) router.MessageHandlerLocker {
	if services.configurationLocker == nil {
		services.configurationLocker = router.NewMessageHandlerLocker()
	}

	return services.configurationLocker
}

// getApplication creates or returns existing application.
func getApplication(cfg *config.Config) cliffApp.App {
	if services.application == nil {
		services.application = app.New(getConfigService(), getLifecycle(cfg), getManifestLoader(cfg), getDeviceReporter(cfg))
	}

	return services.application
}

func getDeviceReporter(cfg *config.Config) reporter.Device {
	if services.deviceReporter == nil {
		services.deviceReporter = reporter.NewDevice(getMQTT(cfg))
	}

	return services.deviceReporter
}

// newRouting creates new set of routing.
func newRouting(cfg *config.Config) []*router.Routing {
	return routing.New(
		getConfigService(),
		getLifecycle(cfg),
		getConfigurationLocker(cfg),
		getApplication(cfg),
		getRestReporter(cfg),
	)
}

// newTasks creates new set of tasks.
func newTasks(cfg *config.Config) []*task.Task {
	return tasks.New(getUDPListener(cfg), getLifecycle(cfg), getApplication(cfg))
}

func getUDPListener(cfg *config.Config) listener.Listener {
	return listener.NewUDP(getConnection(cfg), listener.DecoderFunc(decoder.AddressAndValue), getFimpReporter(cfg))
}

func getFimpReporter(cfg *config.Config) reporter.Reporter {
	return reporter.NewFimp(getMQTT(cfg), transformer.NewFimp(cfg.Devices))
}

func getRestReporter(cfg *config.Config) reporter.Rest {
	return reporter.NewRest(reporter.NewRestPublisher(cfg.HTTP.Host, &http.Client{Timeout: time.Second * 5}), transformer.NewDevice(cfg.Devices))
}

func getConnection(cfg *config.Config) listener.Connection {
	s, err := net.ResolveUDPAddr("udp4", ":"+cfg.UDP.Port)
	if err != nil {
		log.Errorf("error resolving Listener address: %s", err)

		return nil
	}

	c, err := net.ListenUDP("udp4", s)
	if err != nil {
		log.Fatalf("error listening Listener address: %s", err)

		return nil
	}

	return c
}
