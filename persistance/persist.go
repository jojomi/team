package persistance

import (
	"context"
	"fmt"
	"log"
	"path/filepath"

	"github.com/jojomi/go-script/v2"
	"github.com/jojomi/team/ent"
	_ "github.com/mattn/go-sqlite3"
	"github.com/mitchellh/go-homedir"
)

var (
	client       *ent.Client
	databasePath = "~/.team/history.sqlite"
)

func createDatabaseClient() {
	var err error

	sc := script.NewContext()
	databasePath, err := homedir.Expand(databasePath)
	if err != nil {
		log.Fatalf("failed getting home dir: %v", err)
	}
	sc.EnsureDirExists(filepath.Dir(databasePath), 0700)

	client, err = ent.Open("sqlite3", fmt.Sprintf("file:%s?cache=shared&_fk=1", databasePath))
	if err != nil {
		log.Fatalf("failed opening connection to sqlite: %v", err)
	}
	// defer client.Close()

	// run the auto migration tool
	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}
}

func GetDatabaseClient() *ent.Client {
	if client == nil {
		createDatabaseClient()
	}
	return client
}
