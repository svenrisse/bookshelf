package mocks

import (
	"github.com/svenrisse/bookshelf/internal/models"
)

func NewMockModels() models.Models {
	return models.Models{
		Books: &BookModel{},
	}
}
