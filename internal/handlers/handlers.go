package handlers

import (
	"fmt"
	"github.com/codegangsta/martini-contrib/render"
	"github.com/codegangsta/martini-contrib/sessions"
	"github.com/tiburon-777/OTUS_HighLoad/internal/application"
	"github.com/tiburon-777/OTUS_HighLoad/internal/auth"
	"net"
	"net/http"
)

func GetHome(r render.Render) {
	doc := map[string]interface{}{
		"PageTitle":  "Вы имеете доступ к проектам",
	}
	r.HTML(200, "index", doc)
}

func GetSigned(r render.Render) {
	doc := map[string]interface{}{
		"PageTitle":  "page not exists",
	}
	r.HTML(200, "signin", doc)
}

func PostSigned(app application.App, r render.Render) {
	r.Redirect(net.JoinHostPort(app.Config.Server.Address, app.Config.Server.Port)+"/login")
}


func GetUserList(r render.Render) {
	doc := map[string]interface{}{
		"PageTitle":  "page not exists",
	}
	r.HTML(200, "list", doc)
}

func PostLogin(app application.App, session sessions.Session, postedUser auth.UserModel, r render.Render, req *http.Request) {
	user := auth.UserModel{}
	query := fmt.Sprintf("SELECT * FROM users WHERE username=\"%s\" and password =\"%s\"", postedUser.Username, postedUser.Password)
	err := app.DB.QueryRow(query).Scan(&user.Id, &user.Username, &user.Password)

	if err != nil || user.Id==0 {
		r.Redirect(auth.RedirectUrl)
		return
	} else {
		err := auth.AuthenticateSession(session, &user)
		if err != nil {
			r.JSON(500, err)
		}

		params := req.URL.Query()
		redirect := params.Get(auth.RedirectParam)
		r.Redirect(redirect)
		return
	}
}