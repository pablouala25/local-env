package usecase

import (
	"context"

	shareddomain "github.com/Bancar/reauth-bff-aws-lambda/internal/shared/domain"
	shareaderrors "github.com/Bancar/reauth-bff-aws-lambda/internal/shared/errors"
)

type ItemRepository interface {
	CreateItem(context.Context, *shareddomain.Item) error
}

type UseCase struct {
	repo ItemRepository
}

func NewUseCase(r ItemRepository) *UseCase {
	return &UseCase{repo: r}
}

func (u *UseCase) CreateItem(ctx context.Context, it *shareddomain.Item) error {
	if it.ID == "" || it.Message == "" {
		return shareaderrors.ErrInvalidPayload
	}
	return u.repo.CreateItem(ctx, it)
}

