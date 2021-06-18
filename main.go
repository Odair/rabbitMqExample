package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/joho/godotenv"
)

const (
	Migration = `CREATE TABLE IF NOT EXISTS Entrega (
		IdEvento serial PRIMARY KEY,
		Ip text,
		Valor int,
		CreatedAt timestamp with time zone DEFAULT current_timestamp,
		UpdatedAt timestamp)`
)

func main() {
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.NewSyncLogger(logger)
		logger = log.With(logger,
			"service", "evento",
			"time:", log.DefaultTimestampUTC,
			"caller", log.DefaultCaller,
		)
	}

	e := godotenv.Load()

	if e != nil {
		fmt.Println(e)
	}

	level.Info(logger).Log("msg", "service started")
	defer level.Info(logger).Log("msg", "service ended")

	var db *sql.DB
	{
		var err error

		conexao := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", os.Getenv("DbHost"), os.Getenv("DbUser"), os.Getenv("DbPassword"), os.Getenv("DbName"))
		db, err = sql.Open("postgres", conexao)
		if err != nil {
			level.Error(logger).Log("exit", err)
			os.Exit(-1)
		}

		defer db.Close()

		_, err = db.Query(Migration)
		if err != nil {
			level.Info(logger).Log("msg", "migration failed "+err.Error())
			return
		}

	}
}
