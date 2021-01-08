package handlers

import (
	"encoding/base64"
	"fmt"
	"github.com/codegangsta/martini-contrib/render"
	"github.com/codegangsta/martini-contrib/sessions"
	"github.com/tiburon-777/OTUS_HighLoad/internal/application"
	"github.com/tiburon-777/OTUS_HighLoad/internal/auth"
	"log"
	"net/http"
	"time"
)

func GetHome(r render.Render, user auth.User) {
	r.HTML(200, "index",  user)
}

func GetSigned(r render.Render) {
	doc := map[string]interface{}{
		"PageTitle":  "page not exists",
	}
	r.HTML(200, "signup", doc)
}

func PostSigned(app application.App, session sessions.Session, postedUser auth.UserModel, r render.Render, req *http.Request) {
	t, err := time.Parse("2006-1-2", postedUser.FormBirthDate)
	if err != nil {
		e := fmt.Errorf("can't parce date: %w", err)
		log.Println(e)
		doc := map[string]interface{}{
			"Error": e,
		}
		r.HTML(500, "500", doc)
	}
	query := fmt.Sprintf(`INSERT INTO users (username, password, name, surname, birthdate, gender, city, interests)
							values ("%s", "%s", "%s", "%s", "%s", "%s", "%s", "%s")`,
		postedUser.Username,
		base64.StdEncoding.EncodeToString([]byte(postedUser.Username + ":" + postedUser.Password)),
		postedUser.Name,
		postedUser.Surname,
		t.Format("2006-01-02 15:04:05"),
		postedUser.Gender,
		postedUser.City,
		postedUser.Interests,
	)
	_, err = app.DB.Exec(query)
	if err != nil {
		e := fmt.Errorf("can't create account in DB: %w", err)
		log.Println(e)
		doc := map[string]interface{}{
			"Error": e,
		}
		r.HTML(500, "500", doc)
	}
	r.Redirect("/login")
}

func GetUserList(r render.Render) {
	doc := map[string]interface{}{
		"PageTitle":  "page not exists",
	}
	r.HTML(200, "list", doc)
}

func PostLogin(app application.App, session sessions.Session, postedUser auth.UserModel, r render.Render, req *http.Request) {
	hash :=  base64.StdEncoding.EncodeToString([]byte(postedUser.Username + ":" + postedUser.Password))
	user := auth.UserModel{}
	query := fmt.Sprintf("SELECT id FROM users WHERE username=\"%s\" and password =\"%s\"", postedUser.Username, hash)
	err := app.DB.QueryRow(query).Scan(&user.Id)

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