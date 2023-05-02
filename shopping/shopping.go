package shopping

import (
	"database/sql"
	"fmt"
	"time"
)

type itemId string

type item struct {
	id        itemId
	name      string
	quantity  string
	dateAdded time.Time
}

type items []item

type Repository interface {
	SaveItem(i item) error
	ListItems() (items, error)
}

type sqliteRepository struct {
	db *sql.DB
}

func (r *sqliteRepository) SaveItem(i item) error {
	_, err := r.db.Exec("INSERT INTO shopping_list (id, item, quantity, date_added) VALUES (?, ?, ?, ?)", i.id, i.name, i.quantity, i.dateAdded)

	if err != nil {
		return fmt.Errorf("failed to insert record: %w", err)
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
		err = rows.Scan(&i.id, &i.name, &i.quantity, &i.dateAdded)

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

func New(id string, name string, quantity string, dateAdded time.Time) item {
	return item{
		id:        itemId(id),
		name:      name,
		quantity:  quantity,
		dateAdded: dateAdded,
	}
}
