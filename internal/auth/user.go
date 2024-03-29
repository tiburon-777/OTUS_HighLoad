package auth

import (
	"fmt"
	"time"

	"github.com/tiburon-777/OTUS_HighLoad/internal/application"
)

type UserModel struct {
	ID            int64     `db:"id" form:"id"`
	Username      string    `db:"username" form:"username"`
	Password      string    `db:"password" form:"password"`
	Name          string    `db:"name" form:"name"`
	Surname       string    `db:"surname" form:"surname"`
	BirthDate     time.Time `db:"birthdate"`
	YearsOld      int       `db:"-" form:"-"`
	FormBirthDate string    `form:"birthdate"`
	Gender        string    `db:"gender" form:"gender"`
	City          string    `db:"city" form:"city"`
	Interests     string    `db:"interests" form:"interests"`
	IsFriend      bool      `db:"-" form:"-"`
	authenticated bool      `db:"-" form:"-"`
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

func (u *UserModel) UniqueID() interface{} {
	return u.ID
}

func (u *UserModel) GetByID(app application.App, id interface{}) error {
	var v string
	query := fmt.Sprintf("SELECT username, name, surname, birthdate, gender, city, interests FROM users WHERE id=%d", id)
	err := app.DBMaster.QueryRow(query).Scan(&u.Username, &u.Name, &u.Surname, &v, &u.Gender, &u.City, &u.Interests)
	if err != nil {
		return err
	}
	u.BirthDate, err = time.Parse("2006-01-02 15:04:05", v)
	if err != nil {
		return err
	}
	u.ID = id.(int64)
	return nil
}
