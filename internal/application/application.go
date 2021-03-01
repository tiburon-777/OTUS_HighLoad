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
	DBMaster     *sql.DB
	DBSlave1     *sql.DB
	DBSlave2     *sql.DB
}

func New(configFile, envPrefix string) (app App, err error) {
	app.Config, err = configure(configFile, envPrefix)
	if err != nil {
		return App{}, fmt.Errorf("can't apply config: %w\n", err)
	}

	app.DBMaster, err = sql.Open("mysql", app.Config.DSN.User+":"+app.Config.DSN.Pass+"@tcp("+app.Config.DSN.Master+":"+app.Config.DSN.Port+")/"+app.Config.DSN.Base+"?charset=utf8&collation=utf8_unicode_ci")
	if err != nil {
		return App{}, err
	}
	if err = dbInit(app.DBMaster); err != nil {
		return App{}, err
	}
	if app.Config.DSN.Slave1 != "" {
		app.DBSlave1, err = sql.Open("mysql", app.Config.DSN.User+":"+app.Config.DSN.Pass+"@tcp("+app.Config.DSN.Slave1+":"+app.Config.DSN.Port+")/"+app.Config.DSN.Base+"?charset=utf8&collation=utf8_unicode_ci")
		if err != nil {
			return App{}, err
		}
	}

	if app.Config.DSN.Slave2 != "" {
		app.DBSlave2, err = sql.Open("mysql", app.Config.DSN.User+":"+app.Config.DSN.Pass+"@tcp("+app.Config.DSN.Slave2+":"+app.Config.DSN.Port+")/"+app.Config.DSN.Base+"?charset=utf8&collation=utf8_unicode_ci")
		if err != nil {
			return App{}, err
		}
	}
	return app, nil
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
		) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE utf8_unicode_ci`); err != nil {
		return err
	}
	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS relations (
		userId int(11) DEFAULT NULL,
		friendId int(11) DEFAULT NULL
	) ENGINE=InnoDB DEFAULT CHARSET=utf8`); err != nil {
		return err
	}
	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS posts (
		Id INT(11) NOT NULL AUTO_INCREMENT,
		Author INT(11) NULL DEFAULT NULL,
		Created TIMESTAMP NULL DEFAULT NULL,
		Subject VARCHAR(50) NULL DEFAULT NULL,
		Body MEDIUMTEXT NULL DEFAULT NULL,
		PRIMARY KEY (Id) USING BTREE,
		INDEX AuthorID (author) USING BTREE,
		CONSTRAINT AuthorID FOREIGN KEY (author) REFERENCES app.users (Id) ON UPDATE RESTRICT ON DELETE RESTRICT
	) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE utf8_unicode_ci`); err != nil {
		return err
	}
	log.Println("All tables exists")
	return nil
}
