package auth

import (
	"database/sql"
	"fmt"
	"github.com/tiburon-777/OTUS_HighLoad/internal/application"
	"time"
)

type UserModel struct {
	Id            int64  `form:"id" db:"id"`
	Username      string `form:"name" db:"username"`
	Password      string `form:"password" db:"password"`
	Name        string		`form:"name" db:"name"`
	Surname     string		`form:"surname" db:"surname"`
	BirthDate   time.Time	`form:"birthdate" db:"birthdate"`
	Male		bool		`form:"male" db:"male"`
	City		string		`form:"city" db:"city"`
	Interests	string		`form:"interests" db:"interests"`
	authenticated bool   `form:"-" db:"-"`
	Db *sql.DB
}

func GenerateAnonymousUser() User {
	return &UserModel{}
}

func (u *UserModel) Login() {
	// Update last login time
	// Add to logged-in user's list
	// etc ...
	u.authenticated = true
}

// Logout will preform any actions that are required to completely
// logout a user.
func (u *UserModel) Logout() {
	// Remove from logged-in user's list
	// etc ...
	u.authenticated = false
}

func (u *UserModel) IsAuthenticated() bool {
	return u.authenticated
}

func (u *UserModel) UniqueId() interface{} {
	return u.Id
}

func (u *UserModel) GetById(app application.App, id interface{}) error {
	query := fmt.Sprintf("SELECT username FROM users WHERE id=%d", id)
	var v []uint8
	err := app.DB.QueryRow(query).Scan(&v)
	if err != nil {
		return err
	}
	return nil
}