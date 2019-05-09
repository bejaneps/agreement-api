// Package auth authenticates user through google oauth2
package auth

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"github.com/jinzhu/gorm"
	// for sqlite dialects
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var (
	scopes = []string{
		"https://www.googleapis.com/auth/documents",
		"https://www.googleapis.com/auth/drive",
	}
	credentials = os.Getenv("SECRET_CREDENTIALS")
	database    = os.Getenv("DATABASE_FILE")
	config      *oauth2.Config

	//GoogleAccount is user google account for saving agreements
	GoogleAccount string
)

// GetClient returns a google client instance
func GetClient() *http.Client {
	b, err := ioutil.ReadFile(credentials)
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	cred, err := google.CredentialsFromJSON(context.Background(), b, scopes...)
	if err != nil {
		log.Fatalf("Unable to parse credentials file: %v", err)
	}

	tok, err := cred.TokenSource.Token()
	if err != nil {
		log.Fatalf("Unable to get token: %v", err)
	}

	return config.Client(context.Background(), tok)
}

// GetDB returns a connection to sqlite database
func GetDB() *gorm.DB {
	DB, err := gorm.Open("sqlite3", database)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}

	return DB
}
