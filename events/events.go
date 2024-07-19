package events

import (
	"gitbeam/core"
	"gitbeam/events/topics"
	"gitbeam/store"
	"github.com/sirupsen/logrus"
)

type EventHandlers struct {
	logger        *logrus.Logger
	service       *core.GitBeamService
	eventStore    store.EventStore
	subscriptions []func() error
}

func NewEventHandler(
	eventStore store.EventStore,
	logger *logrus.Logger,
	service *core.GitBeamService,
) EventHandlers {
	return EventHandlers{
		logger:     logger.WithField("module", "EventHandler").Logger,
		service:    service,
		eventStore: eventStore,
	}
}

func (e EventHandlers) Listen() {
	useLogger := e.logger.WithField("methodName", "Listen")
	e.subscriptions = append(
		e.subscriptions,
		e.handleRepoCreated,
	)

	for _, sub := range e.subscriptions {
		if err := sub(); err != nil {
			useLogger.WithError(err).Fatal("failed to mount subscription")
		}
	}

	<-make(chan bool)
}

// handleRepoCreated will receive data that a repo has been created.
// it will get the first 100 commits, and register a cron job that will continually check for new commits every 2minutes. (for this take-home exercise, the cron job will run every 2minutes so we can see changes faster )
func (e EventHandlers) handleRepoCreated() error {
	return e.eventStore.Subscribe(topics.RepoCreated, func(event store.Event) error {
		e.logger.Infof("received event on %s", topics.RepoCreated)
		//var repo models.Repo
		//if err := utils.UnPack(event.Data(), &repo); err != nil {
		//	e.logger.WithError(err).Errorf("failed to decode event payload for %s", event.Topic())
		//}

		// TODO: Implement cron job and

		return nil
	})
}
