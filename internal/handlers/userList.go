package handlers

import (
	"github.com/codegangsta/martini-contrib/render"
	"github.com/tiburon-777/OTUS_HighLoad/internal/application"
	"github.com/tiburon-777/OTUS_HighLoad/internal/auth"
	"net/http"
	"time"
)

func GetUserList(app application.App, r render.Render) {
	doc := make(map[string]interface{})
	doc["UsersFound"] = 0
	var tmp int
	if err := app.DBMaster.QueryRow(`SELECT COUNT(*) FROM users`).Scan(&tmp); err != nil {
		err500("can't get total of user profiles from DB: ", err, r)
	}
	doc["UsersTotal"] = tmp
	r.HTML(200, "list", doc)
}

func PostUserList(app application.App, user auth.User, r render.Render, req *http.Request) {
	postName := req.FormValue("name")
	postSurname := req.FormValue("surname")
	doc := make(map[string]interface{})
	doc["user"] = user.(*auth.UserModel)
	var users []auth.UserModel
	var tmp auth.UserModel
	var tmpTime string
	var results, err = app.DBMaster.Query(`SELECT
			users.id as id,
			users.name as name,
			users.surname as surname,
			users.birthdate as birthdate,
			users.gender as gender,
			users.city as city
		FROM
			users
		WHERE
			NOT users.id=?     
			AND users.id NOT IN ( 
				SELECT
					relations.friendId
				FROM
					relations
				WHERE
					relations.userId=?)
			AND ( users.Name LIKE concat(?, '%') AND users.Surname LIKE concat(?, '%') )`,
		user.(*auth.UserModel).Id,
		user.(*auth.UserModel).Id,
		postName,
		postSurname,
	)
	if err != nil || results == nil {
		err500("can't get user list from DB: ", err, r)
	}
	defer results.Close()
	for results.Next() {
		err = results.Scan(&tmp.Id, &tmp.Name, &tmp.Surname, &tmpTime, &tmp.Gender, &tmp.City)
		if err != nil {
			err500("can't scan result from DB: ", err, r)
		}
		tmp.BirthDate = str2Time(tmpTime, r)
		tmp.YearsOld = int(time.Since(tmp.BirthDate).Hours() / 8760)
		users = append(users, tmp)
		if len(users) >= 100 {
			doc["msg"] = "( Too much rows in result. We will display only the first 100. )"
			break
		}
	}
	doc["table"] = users
	doc["UsersFound"] = len(users)
	var uTotal int
	if err := app.DBMaster.QueryRow(`SELECT COUNT(*) FROM users`).Scan(&uTotal); err != nil {
		err500("can't get total of user profiles from DB: ", err, r)
	}
	doc["UsersTotal"] = uTotal
	r.HTML(200, "list", doc)
}

func PostUserSearch(app application.App, r render.Render, req *http.Request) {
	db := app.DBMaster
	if app.Config.DSN.Slave1!="" {
		db = app.DBSlave1
	}

	postName := req.FormValue("name")
	postSurname := req.FormValue("surname")
	doc := make(map[string]interface{})
	var users []auth.UserModel
	var tmp auth.UserModel
	var tmpTime string
	var results, err = db.Query(`SELECT
			users.id as id,
			users.name as name,
			users.surname as surname,
			users.birthdate as birthdate,
			users.gender as gender,
			users.city as city
		FROM
			users
		WHERE
		  	( users.Name LIKE concat(?, '%') AND users.Surname LIKE concat(?, '%') )`,
		postName,
		postSurname,
	)
	if err != nil || results == nil {
		err500("can't get user list from DB: ", err, r)
	}
	defer results.Close()
	for results.Next() {
		err = results.Scan(&tmp.Id, &tmp.Name, &tmp.Surname, &tmpTime, &tmp.Gender, &tmp.City)
		if err != nil {
			err500("can't scan result from DB: ", err, r)
		}
		tmp.BirthDate = str2Time(tmpTime, r)
		tmp.YearsOld = int(time.Since(tmp.BirthDate).Hours() / 8760)
		users = append(users, tmp)
		if len(users) >= 100 {
			doc["msg"] = "( Too much rows in result. We will display only the first 100. )"
			break
		}
	}
	doc["table"] = users
	doc["UsersFound"] = len(users)
	var uTotal int
	if err := db.QueryRow(`SELECT COUNT(*) FROM users`).Scan(&uTotal); err != nil {
		err500("can't get total of user profiles from DB: ", err, r)
	}
	doc["UsersTotal"] = uTotal
	r.HTML(200, "list", doc)
}

