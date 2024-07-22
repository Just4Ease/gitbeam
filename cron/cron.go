package cron

import (
	"context"
	"gitbeam/core"
	"gitbeam/models"
	"gitbeam/repository"
	"github.com/go-co-op/gocron/v2"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

type Service struct {
	logger    *logrus.Logger
	service   *core.GitBeamService
	scheduler gocron.Scheduler
	cronStore repository.CronServiceStore
}

func NewCronService(cronStore repository.CronServiceStore, service *core.GitBeamService, logger *logrus.Logger) *Service {
	scheduler, err := gocron.NewScheduler()
	if err != nil {
		logger.WithError(err).Fatal("Failed to start cron service")
	}

	return &Service{
		logger:    logger.WithField("moduleName", "CronService").Logger,
		service:   service,
		scheduler: scheduler,
		cronStore: cronStore,
	}
}

func (s Service) Start() {
	// Schedule the task to run every 10 minutes
	_, err := s.scheduler.NewJob(
		gocron.DurationJob(time.Minute*10),
		gocron.NewTask(s.executeEntriesInCronTracker),
	)
	if err != nil {
		s.logger.WithError(err).Fatal("Failed to start job")
	}

	s.scheduler.Start() // This is non-blocking.
	<-make(chan bool)   // use this to block and hold the cron service.
}

func (s Service) executeEntriesInCronTracker() {
	useLogger := s.logger.WithField("methodName", "executeEntriesInCronTracker")
	useLogger.Info("Started fetching and saving commits")
	ctx := context.Background()
	trackers, _ := s.cronStore.ListCronTrackers(ctx)
	wg := &sync.WaitGroup{}
	for _, tracker := range trackers {
		wg.Add(1)
		go func(name models.OwnerAndRepoName) {
			defer wg.Done()
			_ = s.service.FetchAndSaveCommits(ctx, models.ListCommitFilter{
				OwnerAndRepoName: name,
			})
		}(models.OwnerAndRepoName{
			OwnerName: tracker.OwnerName,
			RepoName:  tracker.RepoName,
		})
	}
	wg.Wait()
	useLogger.Info("Finished fetching and saving commits")
}
