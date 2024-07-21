package events

import (
	"context"
	"gitbeam/core"
	"gitbeam/events/topics"
	"gitbeam/models"
	"gitbeam/store"
	"gitbeam/utils"
	"github.com/sirupsen/logrus"
	"time"
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

		var repo models.Repo
		_ = utils.UnPack(event.Data(), &repo)
		ctx := context.Background()

		startTimeString, ok := repo.Meta["startTime"].(string)
		var startTime time.Time
		if ok {
			tt, _ := time.Parse(time.RFC3339, startTimeString)
			startTime = tt // Fall back to when the repo was created as the point of mirroring.
		} else {
			startTime = repo.TimeCreated
		}

		e.logger.WithFields(logrus.Fields{
			"repo":      repo,
			"startTime": startTime,
		}).Infoln("started beaming repository commits with the following payload")

		return e.service.FetchAndSaveCommits(ctx, &models.OwnerAndRepoName{
			OwnerName: repo.Owner,
			RepoName:  repo.Name,
		}, startTime)
	})
}
