package application

import (
	"database/sql"
	"fmt"
	"github.com/tiburon-777/OTUS_HighLoad/internal/models"
	"github.com/tiburon-777/modules/core/config"
	"log"
)

type App struct {
	Config *models.Configuration
	DB     *sql.DB
}

func New(configFile, envPrefix string) (App, error) {
	conf, err := configure(configFile, envPrefix)
	if err != nil {
		return App{}, fmt.Errorf("can't apply config: %w\n", err)
	}

	db, err := sql.Open("mysql", conf.DSN.User+":"+conf.DSN.Pass+"@tcp("+conf.DSN.Host+":"+conf.DSN.Port+")/"+conf.DSN.Base)
	if err != nil {
		return App{}, err
	}
	if err = dbInit(db); err != nil {
		return App{}, err
	}
	if err = dbFill(db); err != nil {
		return App{}, err
	}
	return App{Config: conf, DB: db}, nil
}

func configure(fileName string, envPrefix string) (*models.Configuration, error) {
	var conf models.Configuration
	s := config.New(&conf)
	if fileName != "" {
		log.Printf("try to apply config from file %s...\n", fileName)
		if err := s.SetFromFile(fileName); err != nil {
			return &models.Configuration{}, fmt.Errorf("can't apply config from file: %w", err)
		}
	}
	if envPrefix != "" {
		log.Println("try to apply config from environment...")
		if err := s.SetFromEnv(envPrefix); err != nil {
			return &models.Configuration{}, fmt.Errorf("can't apply envvars to config:%w", err)
		}
	}
	return &conf, nil
}

func dbInit(db *sql.DB) error {
	log.Println("Check DB tables consistency...")
	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS users (
			Id int(11) NOT NULL AUTO_INCREMENT,
			Username varchar(255) DEFAULT NULL,
			Password varchar(255) DEFAULT NULL,
			Name varchar(255) DEFAULT NULL,
			Surname varchar(255) DEFAULT NULL,
			BirthDate datetime DEFAULT NULL,
			Gender varchar(255) DEFAULT NULL,
			City varchar(255) DEFAULT NULL,
			Interests varchar(255) DEFAULT NULL,
			PRIMARY KEY (Id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8`); err != nil {
		return err
	}
	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS relations (
		userId int(11) DEFAULT NULL,
		friendId int(11) DEFAULT NULL
	) ENGINE=InnoDB DEFAULT CHARSET=utf8`); err != nil {
		return err
	}
	log.Println("All tables exists")
	return nil
}

func dbFill(db *sql.DB) error {
	log.Println("Try to generate rows and fill the DB...")

	log.Println("All tables exists")
	return nil
}
