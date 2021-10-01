package handlers

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/codegangsta/martini-contrib/render"
	"github.com/codegangsta/martini-contrib/sessions"
	"golang.org/x/crypto/bcrypt"
	// MySQL driver
	_ "github.com/go-sql-driver/mysql"

	"github.com/tiburon-777/OTUS_HighLoad/internal/application"
	"github.com/tiburon-777/OTUS_HighLoad/internal/auth"
)

type Post struct{
	ID		int `db:"Id"`
	Author	string `db:"Author"`
	Created	time.Time `db:"Created"`
	Subject	string `db:"Subject"`
	Body	string `db:"Body"`
}

func GetHome(app application.App, r render.Render, user auth.User) {
	h := user.(*auth.UserModel).BirthDate
	user.(*auth.UserModel).YearsOld = int(time.Since(h).Hours() / 8760)
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
			users JOIN relations
		WHERE
			relations.friendId=users.Id
			AND relations.userId=?
		GROUP BY users.Id`,
		user.(*auth.UserModel).ID)
	if err != nil || results == nil {
		err500("can't get user list from DB: ", err, r)
	}
	defer results.Close()
	for results.Next() {
		err = results.Scan(&tmp.ID, &tmp.Name, &tmp.Surname, &tmpTime, &tmp.Gender, &tmp.City)
		if err != nil {
			err500("can't scan result from DB: ", err, r)
		}
		tmp.BirthDate = str2Time(tmpTime, r)
		tmp.YearsOld = int(time.Since(tmp.BirthDate).Hours() / 8760)
		users = append(users, tmp)
	}
	doc["table"] = users

	var post Post
	var posts []Post
	results, err = app.DBMaster.Query(`SELECT Id, Created, Subject, Body FROM posts WHERE Author=? ORDER BY Created DESC;`, user.(*auth.UserModel).ID)
	if err != nil || results == nil {
		err500("can't get user list from DB: ", err, r)
	}
	defer results.Close()
	for results.Next() {
		err = results.Scan(&post.ID, &tmpTime, &post.Subject, &post.Body)
		if err != nil {
			err500("can't scan result from DB: ", err, r)
		}
		post.Created = str2Time(tmpTime, r)
		posts = append(posts, post)
	}
	doc["posts"] = posts

	r.HTML(200, "index", doc)
}

func GetSignup(r render.Render) {
	r.HTML(200, "signup", nil)
}

func PostSignup(app application.App, postedUser auth.UserModel, r render.Render) {
	if len(postedUser.Username) < 3 || len(postedUser.Password) < 3 {
		doc := map[string]interface{}{
			"msg": "Login and password must be longer then 3 chars",
		}
		r.HTML(200, "signup", doc)
		return
	}
	t, err := time.Parse("2006-1-2", postedUser.FormBirthDate)
	if err != nil {
		e := fmt.Errorf("can't parce date: %w", err)
		log.Println(e)
		doc := map[string]interface{}{
			"Error": e,
		}
		r.HTML(500, "500", doc)
	}
	pHash, err := bcrypt.GenerateFromPassword([]byte(postedUser.Password), bcrypt.DefaultCost)
	if err != nil {
		err500("can't generate password hash: ", err, r)
	}
	_, err = app.DBMaster.Exec(`INSERT INTO users (username, password, name, surname, birthdate, gender, city, interests)
							values (?, ?, ?, ?, ?, ?, ?, ?)`,
		postedUser.Username,
		pHash,
		postedUser.Name,
		postedUser.Surname,
		t.Format("2006-01-02 15:04:05"),
		postedUser.Gender,
		postedUser.City,
		postedUser.Interests,
	)
	if err != nil {
		err500("can't create account in DB: ", err, r)
	}
	r.Redirect("/login")
}

func PostLogin(app application.App, session sessions.Session, postedUser auth.UserModel, r render.Render, req *http.Request) {
	user := auth.UserModel{}
	err1 := app.DBMaster.QueryRow("SELECT id, password FROM users WHERE username=?", postedUser.Username).Scan(&user.ID, &user.Password)
	err2 := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(postedUser.Password))
	if err1 != nil || err2 != nil {
		doc := map[string]interface{}{
			"msg": "Wrong user or password. You may sign in.",
		}
		r.HTML(200, "login", doc)
		r.Redirect(auth.RedirectURL)
		return
	}
	err := auth.AuthenticateSession(session, &user)
	if err != nil {
		err500("can't auth session: ", err, r)
	}
	params := req.URL.Query()
	redirect := params.Get(auth.RedirectParam)
	r.Redirect(redirect)
	return
}

func str2Time(s string, r render.Render) time.Time {
	t, err := time.Parse("2006-01-02 15:04:05", s)
	if err != nil {
		err500("can't parce date: ", err, r)
	}
	return t
}

func err500(s string, err error, r render.Render) {
	e := fmt.Errorf("%s %w", s, err)
	log.Println(e)
	doc := map[string]interface{}{
		"Error": e,
	}
	r.HTML(500, "500", doc)
}
