package api

import (
	"encoding/json"
	"gitbeam/core"
	"gitbeam/events"
	"gitbeam/models"
	"gitbeam/repository"
	"gitbeam/repository/sqlite"
	"gitbeam/store"
	"github.com/go-chi/chi/v5"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestListCommits(t *testing.T) {
	logger := logrus.New()
	service := setupService(logger)
	router := chi.NewMux()
	New(service, nil, logger).Routes(router)

	req, err := http.NewRequest(http.MethodGet, "/commits?ownerName=brave&repoName=brave-browser", nil)
	assert.Nil(t, err)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestGetRepo(t *testing.T) {
	logger := logrus.New()
	service := setupService(logger)
	router := chi.NewMux()
	New(service, nil, logger).Routes(router)

	req, err := http.NewRequest(http.MethodGet, "/repos/brave/brave-browser", nil)
	assert.Nil(t, err)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	var result models.Result
	assert.Nil(t, json.NewDecoder(rr.Body).Decode(&result))
	assert.Equal(t, true, result.Success)
}

func setupService(logger *logrus.Logger) *core.GitBeamService {
	var eventStore store.EventStore
	var dataStore repository.DataStore
	var err error
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logrus.InfoLevel)

	//Using SQLite as the mini persistent storage.
	//( in a real world system, this would be any production level or vendor managed db )
	if dataStore, err = sqlite.NewSqliteRepo("test_data.db"); err != nil {
		logger.WithError(err).Fatal("failed to initialize sqlite database repository")
	}

	// A channel based pub/sub messaging system.
	//( in a real world system, this would be apache-pulsar, kafka, nats.io or rabbitmq )
	eventStore = store.NewEventStore(logger)

	// If the dependencies were more than 3, I would use a variadic function to inject them.
	//Clarity is better here for this exercise.
	service := core.NewGitBeamService(logger, eventStore, dataStore, nil)

	// To handle event-based background activities. ( in a real world system, this would be apache-pulsar, kafka, nats.io or rabbitmq )
	go events.NewEventHandler(eventStore, logger, service).Listen()
	return service
}
