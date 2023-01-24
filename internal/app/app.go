package app

import (
	"database/sql"
	"encoding/json"
	"errors"
	"firebird/internal/handlers"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	_ "github.com/nakagami/firebirdsql"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

type App struct {
	Config  Config
	Queries Queries
	DB      sql.DB
}

type Config struct {
	DatabaseUrl  string
	DatabaseUser string
	DatabasePass string
	Host         string
	Uuid         string
	ServerName   string
	ApiHost      string
	LIMIT        string
}

type Queries struct {
	Queries map[string]Query
}

type Query struct {
	Query   string
	Enabled bool
}

const (
	YYYYMMDDHHMMSS = "2006-01-02 15:04:05"
)

func (config *Config) ConfigLoad(filename string) error {
	err := godotenv.Load(filename)
	if err != nil {
		return errors.New("Error loading " + filename)
	}

	config.DatabaseUrl = os.Getenv("DATABASE_URL")
	if config.DatabaseUrl == "" {
		return errors.New("missing DATABASE_URL")
	}

	config.Uuid = os.Getenv("UUID")
	if config.Uuid == "" {
		return errors.New("missing UUID")
	}

	config.ServerName = os.Getenv("SERVER_NAME")
	if config.ServerName == "" {
		return errors.New("missing SERVER_NAME")
	}

	config.ApiHost = os.Getenv("API_HOST")
	if config.ApiHost == "" {
		return errors.New("missing API_HOST")
	}
	config.LIMIT = os.Getenv("LIMIT")
	if config.LIMIT == "" {
		return errors.New("missing LIMIT")
	}

	return nil
}

func (queries *Queries) QueriesLoad(filename string) error {
	ymlfile, err := ioutil.ReadFile(filename)
	if err != nil {
		return errors.New("Error loading " + filename)
	}

	err = yaml.Unmarshal(ymlfile, &queries.Queries)
	if err != nil {
		return errors.New("Error loading yml")
	}
	if queries.Queries == nil {
		return errors.New("missing Queries")
	}

	return nil
}

func NewApp() App {
	return App{}
}

func (app *App) Setup() error {
	db, err := sql.Open("firebirdsql", app.Config.DatabaseUrl)
	if err != nil {
		return err
	}

	app.DB = *db

	return nil
}

func (app *App) Run() error {
	log.Printf(
		"App started for server '%s', datetime '%s'",
		app.Config.Uuid,
		time.Now().Format(YYYYMMDDHHMMSS),
	)

	var tables []string
	for tableName, query := range app.Queries.Queries {
		if query.Enabled != true {
			continue
		}
		tables = append(tables, tableName)
	}
	tablesData := handlers.LastIds(app.Config.ApiHost, handlers.TablesRequest{
		TablesName: tables,
		Uuid:       app.Config.Uuid,
	})

	for tableName, query := range app.Queries.Queries {
		if query.Enabled != true {
			continue
		}
		lastId, ok := tablesData.TablesIds[tableName]
		if !ok {
			continue
		}
		log.Printf("Start query data for tableName '%s' with lastId - '%d'", tableName, lastId)

		queryFormated := fmt.Sprintf(query.Query, lastId, app.Config.LIMIT)
		queryData := handlers.DBExecute(app.DB, queryFormated)
		queryJson, err := json.Marshal(queryData)
		if err != nil {
			log.Printf("Error with query data for tableName '%s' with lastId - '%d'", tableName, lastId)
		}

		tableData := handlers.SendTableData(app.Config.ApiHost, handlers.TableRequest{
			Uuid: app.Config.Uuid,
			Data: queryJson,
		})
		if tableData.Status != 200 {
			log.Println("SendTableData failed")
		} else {
			log.Println("SendTableData success")
		}

	}

	log.Printf(
		"App finished for server '%s', datetime '%s'",
		app.Config.Uuid,
		time.Now().Format(YYYYMMDDHHMMSS),
	)

	return nil
}
