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

		startTimeString, ok := repo.Meta["fromDate"].(string)
		var fromDate models.Date
		if ok {
			ct, _ := models.Parse(startTimeString)
			//tt, _ := models.Date.(time.RFC3339, startTimeString)
			fromDate = *ct // Fall back to when the repo was created as the point of mirroring.
		} else {
			ct, _ := models.Parse(repo.TimeCreated.Format(time.DateTime))
			fromDate = *ct
		}

		toDateString, ok := repo.Meta["toDate"].(string)
		var toDate *models.Date
		if ok {
			ct, _ := models.Parse(toDateString)
			toDate = ct
		} else {
			toDate = nil
		}

		e.logger.WithFields(logrus.Fields{
			"repo":     repo,
			"fromDate": fromDate,
		}).Infoln("started beaming repository commits with the following payload")

		return e.service.FetchAndSaveCommits(ctx, models.ListCommitFilter{
			OwnerAndRepoName: models.OwnerAndRepoName{
				OwnerName: repo.Owner,
				RepoName:  repo.Name,
			},
			FromDate: &fromDate,
			ToDate:   toDate,
		})
	})
}
