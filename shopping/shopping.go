package shopping

import (
	"database/sql"
	"fmt"
	"time"
)

type ItemId int

type item struct {
	Id        ItemId
	Name      string
	Quantity  string
	DateAdded time.Time
}

type items []item

type Repository interface {
	SaveItem(i item) (*ItemId, error)
	UpdateItem(id *ItemId, i item) error
	ListItems() (items, error)
}

type sqliteRepository struct {
	db *sql.DB
}

func (r *sqliteRepository) SaveItem(i item) (*ItemId, error) {
	res, err := r.db.Exec("INSERT INTO shopping_list (item, quantity, date_added) VALUES (?, ?, ?)", i.Name, i.Quantity, i.DateAdded)

	if err != nil {
		return nil, fmt.Errorf("failed to insert record: %w", err)
	}

	last, err := res.LastInsertId()

	if err != nil {
		return nil, fmt.Errorf("failed to get last inserted id: %w", err)
	}
	id := ItemId(last)

	return &id, nil
}

func (r *sqliteRepository) UpdateItem(id *ItemId, i item) error {
	_, err := r.db.Exec("UPDATE shopping_list SET item = ?, quantity = ? WHERE id = ?", i.Name, i.Quantity, *id)

	if err != nil {
		return fmt.Errorf("failed to update item %d: %w", i.Id, err)
	}

	return nil
}

func (r *sqliteRepository) ListItems() (items, error) {
	rows, err := r.db.Query("SELECT id, item, quantity, date_added FROM shopping_list")

	if err != nil {
		return nil, fmt.Errorf("error querying database: %w", err)

	}

	res := items{}

	for rows.Next() {
		i := item{}
		err = rows.Scan(&i.Id, &i.Name, &i.Quantity, &i.DateAdded)

		if err != nil {
			return nil, fmt.Errorf("failed to deserialize results: %w", err)
		}

		res = append(res, i)
	}

	return res, nil
}

func NewRepository(db *sql.DB) Repository {
	return &sqliteRepository{
		db: db,
	}
}

func New(name string, quantity string, dateAdded time.Time) item {
	return item{
		Name:      name,
		Quantity:  quantity,
		DateAdded: dateAdded,
	}
}
