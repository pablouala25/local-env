package sharedmodels

import shareddomain "github.com/Bancar/reauth-bff-aws-lambda/internal/shared/domain"

type Item struct {
	ID      string `dynamodbav:"id"`
	Message string `dynamodbav:"message"`
}

// FromDomain convierte un domain.Item a models.Item
func FromDomain(item *shareddomain.Item) *Item {
	return &Item{
		ID:      item.ID,
		Message: item.Message,
	}
}

// ToDomain convierte un models.Item a domain.Item
func (m *Item) ToDomain() *shareddomain.Item {
	return &shareddomain.Item{
		ID:      m.ID,
		Message: m.Message,
	}
}
