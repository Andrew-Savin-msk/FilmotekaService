package apiserver

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"

	"github.com/Andrew-Savin-msk/filmoteka-service/backend/internal/config"
	"github.com/Andrew-Savin-msk/filmoteka-service/backend/internal/store"
	"github.com/Andrew-Savin-msk/filmoteka-service/backend/internal/store/pgstore"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"

	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

var (
	ErrDbTypeUnknown = fmt.Errorf("unknown database name")
)

func Start(cfg *config.Config) error {

	st, err := loadStore(cfg.DbPath, cfg.DbType)
	if err != nil {
		return fmt.Errorf("unable to init db. ended with error: %s", err)
	}
	defer st.Close()

	srv := newServer(st, cfg)

	srv.logger.Infof("Api started work")
	err = http.ListenAndServe(cfg.Port, srv.mux)
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

func loadStore(url, sType string) (store.Store, error) {
	switch strings.ToLower(sType) {
	case "postgres", "psql", "pg4":
		return loadPg(url)
	}
	return nil, ErrDbTypeUnknown
}

func loadPg(url string) (store.Store, error) {
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
