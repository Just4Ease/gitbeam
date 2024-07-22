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
		e.handleCronTaskCreated,
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

		fromDate, _ := models.Parse(repo.TimeCreated.Format(time.DateTime))
		toDate, _ := models.Parse(time.Now().Format(time.DateOnly))

		return e.service.FetchAndSaveCommits(ctx, models.CommitFilters{
			OwnerAndRepoName: models.OwnerAndRepoName{
				OwnerName: repo.Owner,
				RepoName:  repo.Name,
			},
			FromDate: fromDate,
			ToDate:   toDate,
		})
	})
}

func (e EventHandlers) handleCronTaskCreated() error {
	return e.eventStore.Subscribe(topics.CronTaskCreated, func(event store.Event) error {
		e.logger.Infof("received event on %s", topics.CronTaskCreated)

		var repo models.Repo
		_ = utils.UnPack(event.Data(), &repo)
		ctx := context.Background()

		startTimeString, ok := repo.Meta["fromDate"].(string)
		var fromDate *models.Date
		if ok {
			fromDate, _ = models.Parse(startTimeString)
		} else {
			fromDate, _ = models.Parse(repo.TimeCreated.Format(time.DateOnly))
		}

		var toDate *models.Date
		toDateString, ok := repo.Meta["toDate"].(string)
		if ok {
			toDate, _ = models.Parse(toDateString)
		} else {
			toDate = nil
		}

		return e.service.FetchAndSaveCommits(ctx, models.CommitFilters{
			OwnerAndRepoName: models.OwnerAndRepoName{
				OwnerName: repo.Owner,
				RepoName:  repo.Name,
			},
			FromDate: fromDate,
			ToDate:   toDate,
		})
	})
}
