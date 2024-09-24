package repositories

import (
	"database/sql"

	"api-server/models"
)

type ItemRepository struct {
	db *sql.DB
}

func NewItemRepository(db *sql.DB) *ItemRepository {
	return &ItemRepository{db: db}
}

func (r *ItemRepository) GetAll() ([]models.Item, error) {
	rows, err := r.db.Query("SELECT id, name FROM items")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.Item
	for rows.Next() {
		var item models.Item
		if err := rows.Scan(&item.ID, &item.Name); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

func (r *ItemRepository) GetByID(id string) (*models.Item, error) {
	var item models.Item
	err := r.db.QueryRow("SELECT id, name FROM items WHERE id = $1", id).Scan(&item.ID, &item.Name)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *ItemRepository) Create(name string) (*models.Item, error) {
	var item models.Item
	err := r.db.QueryRow("INSERT INTO items (name) VALUES ($1) RETURNING id, name", name).Scan(&item.ID, &item.Name)
	if err != nil {
		return nil, err
	}
	return &item, nil
}
