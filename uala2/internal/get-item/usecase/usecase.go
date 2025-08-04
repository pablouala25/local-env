package usecase

import (
	"context"

	shareddomain "github.com/Bancar/reauth-bff-aws-lambda/internal/shared/domain"
	shareaderrors "github.com/Bancar/reauth-bff-aws-lambda/internal/shared/errors"
)

type ItemRepository interface {
	GetItem(ctx context.Context, id string) (*shareddomain.Item, error)
}

type UseCase struct {
	repo ItemRepository
}

func NewUseCase(r ItemRepository) *UseCase {
	return &UseCase{repo: r}
}

func (u *UseCase) GetItem(ctx context.Context, id string) (*shareddomain.Item, error) {
	if id == "" {
		return nil, shareaderrors.ErrInvalidPayload
	}
	return u.repo.GetItem(ctx, id)
}
