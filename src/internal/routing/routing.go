package routing

import (
	"github.com/futurehomeno/cliffhanger/app"
	cliffConfig "github.com/futurehomeno/cliffhanger/config"
	"github.com/futurehomeno/cliffhanger/lifecycle"
	"github.com/futurehomeno/cliffhanger/router"

	"github.com/rafalmnich/edge-iqcontrols-app/internal/config"
	"github.com/rafalmnich/edge-iqcontrols-app/internal/reporter"
)

const (
	ResourceName = "iqcontrols" // ResourceName is the name of the application.
	ServiceName  = "iqcontrols" // ServiceName is the name of the main service of the application.
)

// New creates a new routing table with all message handlers and their voters.
func New(
	cfgSrv *config.Service,
	appLifecycle *lifecycle.Lifecycle,
	configurationLocker router.MessageHandlerLocker,
	application app.App,
	reporter reporter.Rest,
) []*router.Routing {
	return router.Combine(
		[]*router.Routing{
			cliffConfig.RouteCmdLogSetLevel(ServiceName, cfgSrv.SetLogLevel),
			cliffConfig.RouteCmdLogGetLevel(ServiceName, cfgSrv.GetLogLevel),
			routeDeviceCommands(reporter),
		},
		app.RouteApp(ServiceName, appLifecycle, cfgSrv, config.Factory, configurationLocker, application),
	)
}
