package tasks

import (
	"github.com/futurehomeno/cliffhanger/app"
	"github.com/futurehomeno/cliffhanger/lifecycle"
	"github.com/futurehomeno/cliffhanger/task"

	"github.com/rafalmnich/edge-iqcontrols-app/internal/listener"
)

// New creates a new definition of background tasks to be performed by the application.
func New(udpListener listener.Listener, appLifecycle *lifecycle.Lifecycle, application app.App) []*task.Task {
	return task.Combine(
		app.TaskApp(application, appLifecycle),
		listenerTask(udpListener),
	)
}

func listenerTask(listener listener.Listener) []*task.Task {
	return []*task.Task{
		task.New(listener.Start, 0),
	}
}
