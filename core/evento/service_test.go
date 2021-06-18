package evento_test

import (
	"database/sql"
	"testing"

	"github.com/Odair/rabbitMqExample/core/evento"
	_ "github.com/mattn/go-sqlite3"
)

const (
	Migration = `CREATE TABLE IF NOT EXISTS Evento (
		IdEvento serial PRIMARY KEY,
		Ip text,
		Estado text,
		Valor int,
		CreatedAt timestamp with time zone DEFAULT current_timestamp,
		UpdatedAt timestamp)`
)

func TestStore(t *testing.T) {
	e := &evento.Evento{
		IdEvento: 1,
		Ip:       "000.000.0.0",
		Estado:   "Minas Gerais",
		Valor:    9,
	}
	db, err := sql.Open("sqlite3", "../../data/evento_test.db")
	if err != nil {
		t.Fatalf("Erro conectando ao banco de dados %s", err.Error())
	}

	err = startDB(db)
	if err != nil {
		t.Fatalf("Erro iniciar Db%s", err.Error())
	}

	err = clearDB(db)
	if err != nil {
		t.Fatalf("Erro limpando o banco de dados: %s", err.Error())
	}
	defer db.Close()
	service := evento.NewService(db)
	err = service.Store(e)
	if err != nil {
		t.Fatalf("Erro salvando no banco de dados: %s", err.Error())
	}
	saved, err := service.GetAll()
	if err != nil {
		t.Fatalf("Erro buscando do banco de dados: %s", err.Error())
	}
	if len(saved) == 0 {
		t.Fatalf("Dados inválidos. Esperado %d, recebido %d", 1, len(saved))
	}
}

func clearDB(db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Exec("delete from Evento")
	tx.Commit()
	return err
}

func startDB(db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Exec(Migration)
	tx.Commit()
	return err
}
