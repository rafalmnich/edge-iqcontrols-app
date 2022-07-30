package routing

import (
	"github.com/futurehomeno/cliffhanger/app"
	cliffConfig "github.com/futurehomeno/cliffhanger/config"
	"github.com/futurehomeno/cliffhanger/lifecycle"
	"github.com/futurehomeno/cliffhanger/router"
	"github.com/futurehomeno/fimpgo"
	log "github.com/sirupsen/logrus"

	"github.com/rafalmnich/edge-iqcontrols-app/internal/config"
	"github.com/rafalmnich/edge-iqcontrols-app/internal/reporter"
)

const (
	// ResourceName is a name of the resource.
	ResourceName = "iqcontrols" // ResourceName is the name of the application.
	// ServiceName is a name of the service.
	ServiceName = "iqcontrols" // ServiceName is the name of the main service of the application.
)

// New creates a new routing table with all message handlers and their voters.
func New(
	cfgSrv *config.Service,
	appLifecycle *lifecycle.Lifecycle,
	configurationLocker router.MessageHandlerLocker,
	application app.App,
	deviceReporter reporter.Rest,
	modesReporter reporter.Rest,
) []*router.Routing {
	return router.Combine(
		[]*router.Routing{
			cliffConfig.RouteCmdLogSetLevel(ServiceName, cfgSrv.SetLogLevel),
			cliffConfig.RouteCmdLogGetLevel(ServiceName, cfgSrv.GetLogLevel),
			routeDeviceCommands(deviceReporter),
			routeHomeModeEvents(modesReporter),
		},
		app.RouteApp(ServiceName, appLifecycle, cfgSrv, config.Factory, configurationLocker, application),
	)
}

func routeDeviceCommands(reporter reporter.Rest) *router.Routing {
	return router.NewRouting(
		router.MessageHandlerFn(func(message *fimpgo.Message) (reply *fimpgo.Message) {
			err := reporter.Report(message)
			if err != nil {
				log.WithError(err).Error("failed to report message")
			}

			return nil
		}),
		forDeviceCommands(),
	)
}

func routeHomeModeEvents(modesReporter reporter.Rest) *router.Routing {
	return router.NewRouting(
		router.MessageHandlerFn(func(message *fimpgo.Message) (reply *fimpgo.Message) {
			err := modesReporter.Report(message)
			if err != nil {
				log.WithError(err).Error("failed to report message")
			}

			return nil
		}),
		router.ForServiceAndType("vinculum", "evt.pd7.notify"),
	)
}

func forDeviceCommands() router.MessageVoter {
	return router.MessageVoterFn(func(m *fimpgo.Message) bool {
		return m.Addr.ResourceName == ServiceName && m.Addr.MsgType == "cmd" && m.Addr.ResourceType == "dev"
	})
}
