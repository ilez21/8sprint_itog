package main

import (
	"database/sql"
)

type ParcelStore struct {
	db *sql.DB
}

func NewParcelStore(db *sql.DB) ParcelStore {
	return ParcelStore{db: db}
}

func (s ParcelStore) Add(p Parcel) (int, error) {
	// реализуем добавление строки в таблицу parcel
	result, err := s.db.Exec(
		"INSERT INTO parcel (client, status, address, created_at) VALUES (?, ?, ?, ?)",
		p.Client, p.Status, p.Address, p.CreatedAt,
	)
	if err != nil {
		return 0, err
	}

	// получаем идентификатор последней добавленной записи
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (s ParcelStore) Get(number int) (Parcel, error) {
	// реализуем чтение строки по заданному number
	row := s.db.QueryRow(
		"SELECT number, client, status, address, created_at FROM parcel WHERE number = ?",
		number,
	)

	// заполняем объект Parcel данными из таблицы
	p := Parcel{}
	err := row.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
	if err != nil {
		return p, err
	}

	return p, nil
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	// реализуем чтение строк из таблицы parcel по заданному client
	rows, err := s.db.Query(
		"SELECT number, client, status, address, created_at FROM parcel WHERE client = ?",
		client,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// заполняем срез Parcel данными из таблицы
	var res []Parcel
	for rows.Next() {
		p := Parcel{}
		err := rows.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
		if err != nil {
			return nil, err
		}
		res = append(res, p)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return res, nil
}

func (s ParcelStore) SetStatus(number int, status string) error {
	// реализуем обновление статуса в таблице parcel
	_, err := s.db.Exec(
		"UPDATE parcel SET status = ? WHERE number = ?",
		status, number,
	)
	return err
}

func (s ParcelStore) SetAddress(number int, address string) error {
	// реализуем обновление адреса в таблице parcel
	// менять адрес можно только если значение статуса registered
	_, err := s.db.Exec(
		"UPDATE parcel SET address = ? WHERE number = ? AND status = ?",
		address, number, ParcelStatusRegistered,
	)
	return err
}

func (s ParcelStore) Delete(number int) error {
	// реализуем удаление строки из таблицы parcel
	// удалять строку можно только если значение статуса registered
	_, err := s.db.Exec(
		"DELETE FROM parcel WHERE number = ? AND status = ?",
		number, ParcelStatusRegistered,
	)
	return err
}
