package db

import (
	"context"
	"fmt"
	"log"
	"net"

	"cloud.google.com/go/cloudsqlconn"
	gsd "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Settings for CloudSQL
const (
	dbUser                 = "packmanrei"
	dbPwd                  = "Reiand0123"
	dbName                 = "voting-db"
	instanceConnectionName = "voting0195:asia-northeast1:voting-instance"
	usePrivate             = ""
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
	db := OpenDB()
	db.AutoMigrate(&Account{}, &Voting{}, &Condition{}, &Contact{})
}

func OpenDB() *gorm.DB {

	d, err := cloudsqlconn.NewDialer(context.Background())
	CheckError(err)
	var opts []cloudsqlconn.DialOption
	if usePrivate != "" {
		opts = append(opts, cloudsqlconn.WithPrivateIP())
	}
	gsd.RegisterDialContext("cloudsqlconn",
		func(ctx context.Context, addr string) (net.Conn, error) {
			return d.Dial(ctx, instanceConnectionName, opts...)
		})

	dbURI := fmt.Sprintf("%s:%s@cloudsqlconn(localhost:3306)/%s?parseTime=true",
		dbUser, dbPwd, dbName)

	db, err := gorm.Open(mysql.Open(dbURI), &gorm.Config{})
	CheckError(err)
	return db
}

func CheckError(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
