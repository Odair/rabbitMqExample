package evento

import (
	"database/sql"
)

type UseCase interface {
	GetAll() ([]*Evento, error)
	Get(ID int64) (*Evento, error)
	Store(b *Evento) error
}

type Service struct {
	DB *sql.DB
}

func NewService(db *sql.DB) *Service {
	return &Service{
		DB: db,
	}
}

func (s *Service) GetAll() ([]*Evento, error) {

	var result []*Evento

	rows, err := s.DB.Query("Select * from Evento")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var e Evento

		err := rows.Scan(&e.IdEvento, &e.Ip, &e.Estado, &e.Valor)
		if err != nil {
			return nil, err
		}
		result = append(result, &e)
	}

	return result, nil
}

func (s *Service) Get(ID int64) (*Evento, error) {

	var e Evento

	stmt, err := s.DB.Prepare("select IdEvento, Ip, Estado, Valor from beer where IdEvento = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	err = stmt.QueryRow(ID).Scan(&e.IdEvento, &e.Ip, &e.Estado, &e.Valor)
	if err != nil {
		return nil, err
	}

	return &e, nil

}

func (s *Service) Store(e *Evento) error {

	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("insert into beer(Ip, Estado, Valor) values (?,?,?)")
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(e.Ip, e.Estado, e.Valor)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil
}
