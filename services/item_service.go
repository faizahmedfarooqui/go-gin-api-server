package services

import (
	"database/sql"

	"api-server/models"
	"api-server/repositories"
)

type ItemService struct {
	repo *repositories.ItemRepository
}

func NewItemService(db *sql.DB) *ItemService {
	return &ItemService{
		repo: repositories.NewItemRepository(db),
	}
}

func (s *ItemService) GetAllItems() ([]models.Item, error) {
	return s.repo.GetAll()
}

func (s *ItemService) GetItemByID(id string) (*models.Item, error) {
	return s.repo.GetByID(id)
}

func (s *ItemService) CreateItem(name string) (*models.Item, error) {
	return s.repo.Create(name)
}
