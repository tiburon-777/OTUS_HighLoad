package handlers

import (
	"time"

	"github.com/codegangsta/martini-contrib/render"

	"github.com/tiburon-777/OTUS_HighLoad/internal/application"
	"github.com/tiburon-777/OTUS_HighLoad/internal/auth"
)

func GetFeed(app application.App, r render.Render, user auth.User) {
	h := user.(*auth.UserModel).BirthDate
	user.(*auth.UserModel).YearsOld = int(time.Since(h).Hours() / 8760)
	doc := make(map[string]interface{})
	doc["user"] = user.(*auth.UserModel)
	var tmpTime string
	var post Post
	var posts []Post
	var results, err = app.DBMaster.Query(`SELECT
			posts.ID AS Id,
			users.Username AS Author,
			posts.Created AS Created,
			posts.Subject AS Subject,
			posts.Body AS Body
		FROM
			users JOIN relations JOIN posts
		WHERE
			relations.friendId=users.Id
			AND posts.Author=relations.friendId
			AND relations.userId=?
		ORDER by Created DESC`,
		user.(*auth.UserModel).ID)
	if err != nil || results == nil {
		err500("can't get feed from DB: ", err, r)
	}
	defer results.Close()
	for results.Next() {
		err = results.Scan(&post.ID, &post.Author, &tmpTime, &post.Subject, &post.Body)
		if err != nil {
			err500("can't scan result from DB: ", err, r)
		}
		post.Created = str2Time(tmpTime, r)
		posts = append(posts, post)
	}
	doc["posts"] = posts

	r.HTML(200, "feed", doc)
}

