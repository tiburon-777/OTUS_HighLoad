package dataset

import (
	"database/sql"
	"github.com/mdigger/translit"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

type Person struct {
	FirstName  string
	SecondName string
	Password   string
	BirthDate  time.Time
	Gender     string
	City       string
	Interests  []string
}

func NewPerson() (p Person) {
	rand.Seed(time.Now().UnixNano())
	p.Gender = genders[rand.Intn(len(genders))]
	if p.Gender == "male" {
		p.FirstName = manNames[rand.Intn(len(manNames))]
		p.SecondName = secondNames[rand.Intn(len(secondNames))]
	} else {
		p.FirstName = womanNames[rand.Intn(len(womanNames))]
		p.SecondName = secondNames[rand.Intn(len(secondNames))] + "а"
	}
	t := make([]byte, 16)
	rand.Read(t)
	p.Password = string(t)
	p.City = cities[rand.Intn(len(cities))]
	for i := 0; i < (rand.Intn(4) + 3); i++ {
		p.Interests = append(p.Interests, interests[rand.Intn(len(interests))])
	}
	s, _ := time.ParseDuration(strconv.Itoa(rand.Intn(700000)) + "h")
	p.BirthDate = time.Now().Add(-s)
	return
}

func FillDB(db *sql.DB, lim int) {
	var uCount int
	if err := db.QueryRow(`SELECT COUNT(*) FROM users`).Scan(&uCount); err != nil {
		log.Fatalf("can't get total of user profiles from DB: ")
	}
	uCount = lim - uCount
	if uCount <= 0 {
		log.Printf("Ok. We have more users then %s.", lim)
	}
	log.Printf("Try to generate %d rows and fill the DB...", uCount)
	for i := 1; i < uCount; i++ {
		if i%100 == 0 {
			log.Printf("Successfully inserted %d rows", i)
		}
		p := NewPerson()
		if _, err := db.Exec(`INSERT INTO users ( Username, Password, Name, Surname, BirthDate, Gender, City, Interests ) values (?, ?, ?, ?, ?, ?, ?, ?)`,
			translit.Ru(p.FirstName)+strconv.Itoa(i),
			p.Password,
			p.FirstName,
			p.SecondName,
			p.BirthDate.Format("2006-01-02 15:04:05"),
			p.Gender,
			p.City,
			strings.Join(p.Interests, ","),
		); err != nil {
			log.Fatalf("can't insert row in DB. Inserted %d rows of %d: %s", i, uCount, err.Error())
		}
	}
	log.Println("Table USERS filled successfully")
}
