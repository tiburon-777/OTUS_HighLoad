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
	h := user.(*auth.UserModel).BirthDate
	user.(*auth.UserModel).YearsOld = int(time.Since(h).Hours()/8760)
	r.HTML(200, "index",  user)
}

func GetSignup(r render.Render) {
	r.HTML(200, "signup", nil)
}

func PostSignup(app application.App, postedUser auth.UserModel, r render.Render) {
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
		err500("can't create account in DB: ", err, r)
	}
	r.Redirect("/login")
}

func GetUserList(app application.App, user auth.User, r render.Render) {
	doc := make(map[string]interface{})
	doc["user"]=user.(*auth.UserModel)
	var users []auth.UserModel
	var tmp auth.UserModel
	var tmpTime string
	var results, err = app.DB.Query(`SELECT name, surname, birthdate, gender, city FROM users`)
	if err != nil || results==nil {
		err500("can't get user list from DB: ", err, r)
	}
	defer  results.Close()
	for results.Next() {
		err = results.Scan(&tmp.Name, &tmp.Surname, &tmpTime, &tmp.Gender, &tmp.City)
		if err != nil {
			err500("can't scan result from DB: ", err, r)
		}
		tmp.BirthDate = str2Time(tmpTime, r)
		tmp.YearsOld = int(time.Since(tmp.BirthDate).Hours()/8760)
		users = append(users,tmp)
	}
	doc["table"]=users
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
			err500("can't auth session: ", err, r)
		}
		params := req.URL.Query()
		redirect := params.Get(auth.RedirectParam)
		r.Redirect(redirect)
		return
	}
}

func str2Time(s string, r render.Render) time.Time {
	t, err := time.Parse("2006-01-02 15:04:05", s)
	if err != nil {
		err500("can't parce date: ", err, r)
	}
	return t
}

func err500(s string, err error, r render.Render) {
	e := fmt.Errorf("s% %w", s, err)
	log.Println(e)
	doc := map[string]interface{}{
		"Error": e,
	}
	r.HTML(500, "500", doc)
}
