package main

import (
	"fmt"
	"github.com/codegangsta/martini"
	"github.com/codegangsta/martini-contrib/binding"
	"github.com/codegangsta/martini-contrib/render"
	"github.com/codegangsta/martini-contrib/sessions"
	_ "github.com/go-sql-driver/mysql"
	"github.com/tiburon-777/OTUS_HighLoad/internal/application"
	"github.com/tiburon-777/OTUS_HighLoad/internal/auth"
	"github.com/tiburon-777/OTUS_HighLoad/internal/handlers"
	"github.com/tiburon-777/OTUS_HighLoad/pkg/dataset"
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
	m := martini.Classic()
	app, err := application.New("", "APP")
	if err != nil {
		log.Fatal(fmt.Errorf("can't build app: %w", err).Error())
	}
	go dataset.FillDB(app.DBMaster, 1000000)

	m.Map(log.New(os.Stdout, "[app]", log.Lshortfile))
	m.Map(app)
	m.Use(sessions.Sessions("app", sessions.NewCookieStore([]byte("BfyfgIyngIOUgmOIUgt87thrg5RHn78b"))))
	m.Use(auth.SessionUser(auth.GenerateAnonymousUser))
	m.Use(render.Renderer(render.Options{
		Directory:  "templates",
		Extensions: []string{".tmpl"},
	}))

	auth.RedirectUrl = "/login"
	auth.RedirectParam = "next"

	m.Get("/404", func(r render.Render) {
		r.HTML(200, "404", nil)
	})
	m.Get("/login", func(r render.Render) {
		r.HTML(200, "login", nil)
	})
	m.Post("/login", binding.Bind(auth.UserModel{}), handlers.PostLogin)

	m.Get("/logout", auth.LoginRequired, func(session sessions.Session, user auth.User, r render.Render) {
		auth.Logout(session, user)
		r.Redirect("/")
	})

	// Регистрация пользователя, после которой нас перебрасывает на страницу логина
	m.Get("/signup", handlers.GetSignup)
	m.Post("/signup", binding.Bind(auth.UserModel{}), handlers.PostSignup)

	m.Get("/subscribe", handlers.GetSubscribe)
	m.Get("/unsubscribe", handlers.GetUnSubscribe)

	//Анкета текущего пользователя
	m.Get("/", auth.LoginRequired, handlers.GetHome)

	m.Get("/list", auth.LoginRequired, handlers.GetUserList)
	m.Post("/list", auth.LoginRequired, handlers.PostUserList)

	m.Get("/search", handlers.GetUserList)
	m.Post("/search", handlers.PostUserSearch)

	m.NotFound(func(r render.Render) {
		r.HTML(404, "404", nil)
	})
	if err := http.ListenAndServe(net.JoinHostPort(app.Config.Server.Address, app.Config.Server.Port), m); err != nil {
		log.Fatalln(err.Error())
	}
}
