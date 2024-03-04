package main

import (
	"database/sql"
	"fmt"
)

type ParcelStore struct {
	db *sql.DB
}

func NewParcelStore(db *sql.DB) ParcelStore {
	return ParcelStore{db: db}
}

func (s ParcelStore) Add(p Parcel) (int, error) {

	res, err := s.db.Exec("INSERT INTO parcel(client, status, address, created_at) VALUES(:client, :status, :address, :created_at)",
		sql.Named("Client", p.Client),
		sql.Named("Status", p.Status),
		sql.Named("Address", p.Address),
		sql.Named("Created_at", p.CreatedAt),
	)
	if err != nil {

		return 0, fmt.Errorf("Ошибка добавления: %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("Ошибка добавления: %w", err)
	}
	// верните идентификатор последней добавленной записисс
	return int(id), nil
}

func (s ParcelStore) Get(number int) (Parcel, error) {

	row := s.db.QueryRow("SELECT * FROM parcel WHERE number = :number",
		sql.Named("number", number))

	p := Parcel{}
	err := row.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
	if err != nil {
		return p, fmt.Errorf("Получена ошибка :%w", err)
	}

	return p, nil
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {

	// здесь из таблицы может вернуться несколько строк
	rows, err := s.db.Query("SELECT * FROM parcel WHERE client = :client",
		sql.Named("client", client))
	if err != nil {
		return nil, fmt.Errorf("Ошибка клиента: %w", err)
	}
	defer rows.Close()

	var res []Parcel
	for rows.Next() {
		p := Parcel{}
		err := rows.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("Ошибка клиента: %w", err)
		}
		res = append(res, p)
	}

	return res, nil
}

func (s ParcelStore) SetStatus(number int, status string) error {

	_, err := s.db.Exec("UPDATE parsel SET status = :status where number = :number",
		sql.Named("status", status),
		sql.Named("number", number),
	)
	if err != nil {
		return fmt.Errorf("Ошибка обновления статуса: %w", err)
	}
	return nil
}

func (s ParcelStore) SetAddress(number int, address string) error {

	status := "registered"

	_, err := s.db.Exec("UPDATE parcel SET address = :address WHERE number = :number AND status = :status",
		sql.Named("address", address),
		sql.Named("number", number),
		sql.Named("status", status))

	if err != nil {
		return fmt.Errorf("Ошибка обновления адреса: %w", err)
	}

	return nil
}

func (s ParcelStore) Delete(number int) error {
	// удалять строку можно только если значение статуса registered
	status := "registered"
	_, err := s.db.Exec("DELETE FROM parcel WHERE number = :number AND status = :status",
		sql.Named("number", number),
		sql.Named("status", status))
	if err != nil {
		return fmt.Errorf("Ошибка удаления: %w", err)
	}
	return nil
}
