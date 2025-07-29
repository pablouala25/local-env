package dto

import (
	shareddomain "github.com/Bancar/reauth-bff-aws-lambda/internal/shared/domain"
)

type Item struct {
	ID      string `json:"id"`
	Message string `json:"message"`
}

func FromDomain(item *shareddomain.Item) *Item {
	return &Item{
		ID:      item.ID,
		Message: item.Message,
	}
}

func (i *Item) ToDomain() *shareddomain.Item {
	return &shareddomain.Item{
		ID:      i.ID,
		Message: i.Message,
	}
}
