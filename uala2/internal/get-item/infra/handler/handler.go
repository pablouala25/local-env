package handler

import (
	"context"
	"net/http"
	"strings"

	shareddomain "github.com/Bancar/reauth-bff-aws-lambda/internal/shared/domain"
	shareaderrors "github.com/Bancar/reauth-bff-aws-lambda/internal/shared/errors"

	apierrors "github.com/Bancar/uala-auth-team-go-dependencies/libs/errors"
	"github.com/Bancar/uala-auth-team-go-dependencies/libs/http/response"
	"github.com/Bancar/uala-auth-team-go-dependencies/libs/tracing"
)

type usecase interface {
	GetItem(context.Context, string) (*shareddomain.Item, error)
}

type Handler struct {
	uc usecase
}

func NewHandler(uc usecase) *Handler {
	return &Handler{uc: uc}
}

func (h *Handler) GetItem(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracing.Span(r.Context(), "Items/Get")
	defer span.End()

	id := lastSegment(r.URL.Path)
	if id == "" {
		id = r.URL.Query().Get("id")
	}
	it, err := h.uc.GetItem(ctx, strings.TrimSpace(id))
	if err != nil {
		switch err {
		case shareaderrors.ErrInvalidPayload:
			response.BadRequest(ctx, w, err.Error())
		case shareaderrors.ErrNotFound:
			response.Error(ctx, w, &apierrors.APIError{
				Message:    "item not found",
				Code:       "not_found",
				StatusCode: http.StatusNotFound,
			})
		default:
			span.SetErrorAttributes("usecase_get_error", tracing.Error(err.Error()))
			response.Error(ctx, w, apierrors.UnexpectedError(err))
		}
		return
	}

	response.Ok(ctx, w, it)
}

func lastSegment(path string) string {
	if path == "" {
		return ""
	}
	p := strings.TrimSuffix(path, "/")
	parts := strings.Split(p, "/")
	if len(parts) == 0 {
		return ""
	}
	return parts[len(parts)-1]
}
