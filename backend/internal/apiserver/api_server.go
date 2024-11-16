package apiserver

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"

	brockerclient "github.com/Andrew-Savin-msk/filmoteka-service/backend/internal/broker_client"
	kafkaclient "github.com/Andrew-Savin-msk/filmoteka-service/backend/internal/broker_client/kafka_client"
	rabbitclient "github.com/Andrew-Savin-msk/filmoteka-service/backend/internal/broker_client/rabbit_client"
	"github.com/Andrew-Savin-msk/filmoteka-service/backend/internal/config"
	"github.com/Andrew-Savin-msk/filmoteka-service/backend/internal/repostore"
	"github.com/Andrew-Savin-msk/filmoteka-service/backend/internal/repostore/pgstore"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"

	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

var (
	ErrDbTypeUnknown = fmt.Errorf("unknown database name")
	ErrBcTypeUnknown = fmt.Errorf("unknown type of broker client")
)

func Start(cfg *config.Config) error {

	st, err := loadStore(cfg.Db)
	if err != nil {
		return fmt.Errorf("unable to init db. ended with error: %s", err)
	}
	defer st.Close()

	log := setLog(cfg.Srv.LogLevel)

	bc, err := loadBrokerClient(cfg.Bc, log)
	if err != nil {
		return fmt.Errorf("unable to load broker client. ended with error: %s", err)
	}

	srv := newServer(st, bc, log, cfg)

	srv.logger.Infof("Api started work")
	err = http.ListenAndServe(cfg.Srv.Port, srv.mux)
	if err != nil {
		return fmt.Errorf("unable to start listening port. ended with error: %s", err)
	}
	return nil
}

func setLog(level string) *logrus.Logger {
	log := logrus.New()
	switch strings.ToLower(level) {
	case "debug":
		log.SetLevel(logrus.DebugLevel)
	case "error":
		log.SetLevel(logrus.ErrorLevel)
	case "info":
		log.SetLevel(logrus.InfoLevel)
	}
	fmt.Printf("logger set in level: %s \n", level)
	return log
}

func loadStore(cfg config.Database) (repostore.Store, error) {
	switch strings.ToLower(cfg.DbType) {
	case "postgres", "psql", "pg4":
		return loadPg("postgresql://" + cfg.User + ":" + cfg.Password + "@" + cfg.Host + ":5432/" + cfg.DbName + "?sslmode=disable")
	}
	return nil, ErrDbTypeUnknown
}

func loadPg(url string) (repostore.Store, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, fmt.Errorf("open: %v", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return pgstore.New(db), nil
}

func loadBrokerClient(cfg config.Broker, logger *logrus.Logger) (brockerclient.Client, error) {
	switch strings.ToLower(cfg.BrokerType) {
	case "rabbitmq", "rabbit_mq", "rabbit":
		return rabbitclient.New(cfg, logrus.NewEntry(logger))
	case "kafka", "apache-kafka", "mannaya":
		return kafkaclient.New(cfg, logrus.NewEntry(logger))
	}
	return nil, ErrBcTypeUnknown
}
