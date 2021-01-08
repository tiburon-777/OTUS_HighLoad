package application

import (
	"database/sql"
	"fmt"
	"github.com/tiburon-777/OTUS_HighLoad/internal/models"
	"github.com/tiburon-777/modules/core/config"
)

type App struct {
	Config *models.Configuration
	DB *sql.DB
}

func New(configFile, envPrefix string) (App, error) {
	conf, err := configure(configFile, envPrefix)
	if err != nil{
		return App{}, fmt.Errorf("can't apply config: %w\n",err)
	}

	db, err := sql.Open("mysql", conf.DSN.User+":"+conf.DSN.Pass+"@tcp("+conf.DSN.Host+":"+conf.DSN.Port+")/"+conf.DSN.Base)
	if err != nil {
		panic(err.Error())
	}

	return App{Config: conf, DB: db}, nil
}

func configure(fileName string, envPrefix string) (*models.Configuration,error) {
	var conf models.Configuration
	s := config.New(&conf)
	if fileName != "" {
		fmt.Printf("try to apply config from file %s...\n", fileName)
		if err := s.SetFromFile(fileName); err != nil {
			return &models.Configuration{}, fmt.Errorf("can't apply config from file: %w", err)
		}
	}
	if envPrefix != "" {
		fmt.Printf("try to apply config from environment...\n")
		if err := s.SetFromEnv(envPrefix); err != nil {
			return &models.Configuration{}, fmt.Errorf("can't apply envvars to config:%w", err)
		}
	}
	return &conf, nil
}