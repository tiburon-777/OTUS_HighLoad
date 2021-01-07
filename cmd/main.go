package main

import (
	"github.com/codegangsta/martini"
	"github.com/codegangsta/martini-contrib/binding"
	"github.com/codegangsta/martini-contrib/render"
	"github.com/codegangsta/martini-contrib/sessionauth"
	"github.com/codegangsta/martini-contrib/sessions"
	_ "github.com/go-sql-driver/mysql"
	"github.com/tiburon-777/OTUS_HighLoad/internal/application"
	"github.com/tiburon-777/OTUS_HighLoad/internal/handlers"
	"github.com/tiburon-777/OTUS_HighLoad/internal/models"
	"log"
	"net"
	"net/http"
	"os"
	"time"
)

func init() {
	http.DefaultClient.Timeout = time.Second * 30
}

func main() {
	log.Println("Starting...")
	app, err := application.New("application.conf", "APP")
	if err != nil{
		panic(err.Error())
	}

	m := martini.Classic()

	m.Map(log.New(os.Stdout, "[app]", log.Lshortfile))
	m.Map(app)
	m.Use(sessions.Sessions("app", sessions.NewCookieStore([]byte("BfyfgIyngIOUgmOIUgt87thrg5RHn78b"))))
	m.Use(sessionauth.SessionUser(models.GenerateAnonymousUser))
	m.Use(render.Renderer(render.Options{
		Directory: "templates",
		Extensions: []string{".tmpl"},
	}))

	sessionauth.RedirectUrl = "/login"
	sessionauth.RedirectParam = "next"

	m.Get("/404", func(r render.Render) {
		r.HTML(200, "404", nil)
	})
	m.Get("/login", func(r render.Render) {
		r.HTML(200, "login", nil)
	})
	m.Post("/login", binding.Bind(models.UserModel{}), handlers.PostLogin)

	m.Get("/logout", sessionauth.LoginRequired, func(session sessions.Session, user sessionauth.User, r render.Render) {
		sessionauth.Logout(session, user)
		r.Redirect("/")
	})

	// Регистрация пользователя, после которой нас перебрасывает на страницу логина
	m.Get("/signup", handlers.GetSigned)
	m.Post("/signup", handlers.PostSigned)

	//Анкета текущего пользователя
	m.Get("/", sessionauth.LoginRequired, handlers.GetHome)

	m.Get("/list", sessionauth.LoginRequired, handlers.GetUserList)

	m.NotFound(func(r render.Render) {
		r.HTML(404, "404", nil)
	})
	if err := http.ListenAndServe(net.JoinHostPort(app.Config.Server.Address, app.Config.Server.Port), m); err != nil {
		log.Fatalln(err.Error())
	}
}
