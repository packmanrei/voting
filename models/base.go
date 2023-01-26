package db

import (
	"fmt"
	"log"

	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

const (
	DB       = "test.db"
	server   = "votingserver.database.windows.net"
	port     = 1433
	user     = "packmanrei"
	password = "******"
	database = "votingDB"
)

type Account struct {
	gorm.Model
	Username string
	Password string
	Sex      string
	Age      int
	Email    string
}

type Voting struct {
	gorm.Model
	PosterID int
	VotersID string
	RoomName string
	Title    string
	Content  string
	Choices  string
	Votes    string
	Watchers string
}

type Condition struct {
	gorm.Model
	VotingID int
	Sex      string
	Age      string
	AgeNum   int
}

type VotingResult struct {
	Choice string
	Vote   string
	Result int
}

type Contact struct {
	Email   string
	Content string
}

func init() {
	db := openDB()
	db.AutoMigrate(&Account{}, &Voting{}, &Condition{}, &Contact{})
}

func openDB() *gorm.DB {

	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s;",
		server, user, password, port, database)
	db, err := gorm.Open(sqlserver.Open(connString), &gorm.Config{})

	CheckError(err)
	return db
}

func CheckError(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
