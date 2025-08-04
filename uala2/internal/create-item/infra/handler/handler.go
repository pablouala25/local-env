package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Bancar/uala-auth-team-go-dependencies/libs/http/response"
	"github.com/Bancar/uala-auth-team-go-dependencies/libs/tracing"

	shareaderrors "github.com/Bancar/reauth-bff-aws-lambda/internal/shared/errors"

	"github.com/Bancar/reauth-bff-aws-lambda/internal/create-item/infra/handler/dto"
	shareddomain "github.com/Bancar/reauth-bff-aws-lambda/internal/shared/domain"
	apierrors "github.com/Bancar/uala-auth-team-go-dependencies/libs/errors"
)

type usecase interface {
	CreateItem(context.Context, *shareddomain.Item) error
}

type Handler struct {
	uc usecase
}

func NewHandler(uc usecase) *Handler {
	return &Handler{uc: uc}
}

func (h *Handler) CreateItem(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracing.Span(r.Context(), "Items/Create")
	defer span.End()

	var in dto.Item
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		response.BadRequest(ctx, w, "invalid JSON body")
		return
	}

	err := h.uc.CreateItem(ctx, in.ToDomain())
	if err != nil {
		switch err {
		case shareaderrors.ErrInvalidPayload:
			response.BadRequest(ctx, w, err.Error())
		default:
			span.SetErrorAttributes("usecase_create_error", tracing.Error(err.Error()))
			response.Error(ctx, w, apierrors.UnexpectedError(err))
		}
		return
	}

	response.Ok(ctx, w, map[string]string{"status": "created", "id": in.ID})
}

// Opcional: health/status
func (h *Handler) Status(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_, _ = fmt.Fprint(w, `{"status":"ok"}`)
}
