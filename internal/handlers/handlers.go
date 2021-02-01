package handlers

import (
	"fmt"
	"github.com/codegangsta/martini-contrib/render"
	"github.com/codegangsta/martini-contrib/sessions"
	_ "github.com/go-sql-driver/mysql"
	"github.com/tiburon-777/OTUS_HighLoad/internal/application"
	"github.com/tiburon-777/OTUS_HighLoad/internal/auth"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"strconv"
	"time"
)

func GetHome(app application.App, r render.Render, user auth.User) {
	h := user.(*auth.UserModel).BirthDate
	user.(*auth.UserModel).YearsOld = int(time.Since(h).Hours() / 8760)
	doc := make(map[string]interface{})
	doc["user"] = user.(*auth.UserModel)
	var users []auth.UserModel
	var tmp auth.UserModel
	var tmpTime string
	var results, err = app.DB.Query(`SELECT
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
		user.(*auth.UserModel).Id)
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
	}
	doc["table"] = users

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
	_, err = app.DB.Exec(`INSERT INTO users (username, password, name, surname, birthdate, gender, city, interests)
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

func GetUserList(app application.App, r render.Render) {
	doc := make(map[string]interface{})
	doc["UsersFound"] = 0
	var tmp int
	if err := app.DB.QueryRow(`SELECT COUNT(*) FROM users`).Scan(&tmp); err != nil {
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
	var results, err = app.DB.Query(`SELECT
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
	if err := app.DB.QueryRow(`SELECT COUNT(*) FROM users`).Scan(&uTotal); err != nil {
		err500("can't get total of user profiles from DB: ", err, r)
	}
	doc["UsersTotal"] = uTotal
	r.HTML(200, "list", doc)
}

func PostUserSearch(app application.App, r render.Render, req *http.Request) {
	postName := req.FormValue("name")
	postSurname := req.FormValue("surname")
	doc := make(map[string]interface{})
	var users []auth.UserModel
	var tmp auth.UserModel
	var tmpTime string
	var results, err = app.DB.Query(`SELECT
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
	if err := app.DB.QueryRow(`SELECT COUNT(*) FROM users`).Scan(&uTotal); err != nil {
		err500("can't get total of user profiles from DB: ", err, r)
	}
	doc["UsersTotal"] = uTotal
	r.HTML(200, "list", doc)
}

func PostLogin(app application.App, session sessions.Session, postedUser auth.UserModel, r render.Render, req *http.Request) {
	user := auth.UserModel{}
	err1 := app.DB.QueryRow("SELECT id, password FROM users WHERE username=?", postedUser.Username).Scan(&user.Id, &user.Password)
	err2 := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(postedUser.Password))
	if err1 != nil || err2 != nil {
		doc := map[string]interface{}{
			"msg": "Wrong user or password. You may sign in.",
		}
		r.HTML(200, "login", doc)
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

func GetSubscribe(app application.App, r render.Render, user auth.User, req *http.Request) {
	sid, ok := req.URL.Query()["id"]
	if !ok {
		err500("can't parce URL query", nil, r)
	}
	did, err := strconv.Atoi(sid[0])
	if err != nil {
		err500("can't convert URL query value: ", err, r)
	}
	_, err = app.DB.Exec(`REPLACE INTO relations (userId, friendId) values (?, ?)`, user.(*auth.UserModel).Id, did)
	if err != nil {
		err500("can't create relation in DB: ", err, r)
	}
	_, err = app.DB.Exec(`REPLACE INTO relations (userId, friendId) values (?, ?)`, did, user.(*auth.UserModel).Id)
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
	_, err = app.DB.Exec(`DELETE FROM relations WHERE (userId,friendId) IN ((?, ?),(?, ?))`, user.(*auth.UserModel).Id, did, did, user.(*auth.UserModel).Id)
	if err != nil {
		err500("can't remove relation from DB: ", err, r)
	}
	r.Redirect("/")

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
