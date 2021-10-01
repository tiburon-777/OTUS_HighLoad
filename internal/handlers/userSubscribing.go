package handlers

import (
	"net/http"
	"strconv"

	"github.com/codegangsta/martini-contrib/render"

	"github.com/tiburon-777/OTUS_HighLoad/internal/application"
	"github.com/tiburon-777/OTUS_HighLoad/internal/auth"
)

func GetSubscribe(app application.App, r render.Render, user auth.User, req *http.Request) {
	sid, ok := req.URL.Query()["id"]
	if !ok {
		err500("can't parce URL query", nil, r)
	}
	did, err := strconv.Atoi(sid[0])
	if err != nil {
		err500("can't convert URL query value: ", err, r)
	}
	_, err = app.DBMaster.Exec(`REPLACE INTO relations (userId, friendId) values (?, ?)`, user.(*auth.UserModel).ID, did)
	if err != nil {
		err500("can't create relation in DB: ", err, r)
	}
	_, err = app.DBMaster.Exec(`REPLACE INTO relations (userId, friendId) values (?, ?)`, did, user.(*auth.UserModel).ID)
	if err != nil {
		err500("can't create relation in DB: ", err, r)
	}
	r.Redirect("/list")
}

func GetUnSubscribe(app application.App, r render.Render, user auth.User, req *http.Request) {
	sid, ok := req.URL.Query()["id"]
	if !ok {
		err500("can't parce URL query", nil, r)
	}
	did, err := strconv.Atoi(sid[0])
	if err != nil {
		err500("can't convert URL query value: ", err, r)
	}
	_, err = app.DBMaster.Exec(`DELETE FROM relations WHERE (userId,friendId) IN ((?, ?),(?, ?))`, user.(*auth.UserModel).ID, did, did, user.(*auth.UserModel).ID)
	if err != nil {
		err500("can't remove relation from DB: ", err, r)
	}
	r.Redirect("/")

}
