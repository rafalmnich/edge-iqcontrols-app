package routing

import (
	"github.com/futurehomeno/cliffhanger/router"
	"github.com/futurehomeno/fimpgo"
	log "github.com/sirupsen/logrus"

	"github.com/rafalmnich/edge-iqcontrols-app/internal/reporter"
)

func routeDeviceCommands(
	reporter reporter.Rest,
) *router.Routing {
	return router.NewRouting(
		router.MessageHandlerFn(func(message *fimpgo.Message) (reply *fimpgo.Message) {
			err := reporter.Report(message)
			if err != nil {
				log.WithError(err).Error("failed to report message")
			}

			return nil
		}),
	)
}
