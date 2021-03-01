package handlers

import (
	"github.com/codegangsta/martini-contrib/render"
	"github.com/tiburon-777/OTUS_HighLoad/internal/application"
	"github.com/tiburon-777/OTUS_HighLoad/internal/auth"
	"net/http"
	"time"
)

func GetAddPost(app application.App, r render.Render) {
	r.HTML(200, "postadd", nil)
}

func PostAddPost(app application.App, user auth.User, r render.Render, req *http.Request) {
	postSubj := req.FormValue("subj")
	postBody := req.FormValue("body")
	var results, err = app.DBMaster.Query(`INSERT INTO posts (Author, Created, Subject, Body) VALUES (
        ?, ?, ?, ?)`,
        user.(*auth.UserModel).Id,
        time.Now().Format("2006-01-02 15:04:05"),
		postSubj,
		postBody,
	)
	if err != nil || results == nil {
		err500("can't add new post: ", err, r)
	}
	r.Redirect("/",302)
}
