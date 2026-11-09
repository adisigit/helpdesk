package handler

import (
	"context"
	"errors"
	"helpdesk/backend/internal/database/models"
	"helpdesk/backend/internal/service"
	"time"

	"github.com/danielgtaylor/huma/v2"
)

type UserHandler struct{ svc *service.UserService }

func NewUserHandler(svc *service.UserService) *UserHandler { return &UserHandler{svc: svc} }

type AdminUserBody struct {
	UserBody
	CreatedAt time.Time `json:"created_at"`
}

func toAdminUserBody(u *models.User) AdminUserBody {
	return AdminUserBody{UserBody: toUserBody(u), CreatedAt: u.CreatedAt}
}

func mapUserErr(err error) error {
	switch {
	case errors.Is(err, service.ErrUserNotFound):
		return huma.Error404NotFound("user not found")
	case errors.Is(err, service.ErrLastAdmin):
		return huma.Error409Conflict(err.Error())
	case errors.Is(err, service.ErrInvalidRole), errors.Is(err, service.ErrInvalidFilter):
		return huma.Error422UnprocessableEntity(err.Error())
	}
	return err
}

type ListUsersInput struct {
	Role     string `query:"role" enum:"requester,agent,admin"`
	Q        string `query:"q" maxLength:"100" doc:"Cari nama atau email"`
	Page     int    `query:"page" minimum:"1" default:"1"`
	PageSize int    `query:"page_size" minimum:"1" maximum:"100" default:"20"`
}

type ListUsersOutput struct {
	Body struct {
		Items      []AdminUserBody `json:"items"`
		Total      int64           `json:"total"`
		Page       int             `json:"page"`
		PageSize   int             `json:"page_size"`
		TotalPages int             `json:"total_pages"`
	}
}

func (h *UserHandler) List(ctx context.Context, in *ListUsersInput) (*ListUsersOutput, error) {
	res, err := h.svc.List(ctx, service.ListUsersInput{
		Role: models.Role(in.Role), Query: in.Q, Page: in.Page, PageSize: in.PageSize,
	})
	if err != nil {
		return nil, mapUserErr(err)
	}
	out := &ListUsersOutput{}
	out.Body.Items = make([]AdminUserBody, 0, len(res.Items))
	for i := range res.Items {
		out.Body.Items = append(out.Body.Items, toAdminUserBody(&res.Items[i]))
	}
	out.Body.Total = res.Total
	out.Body.Page = res.Page
	out.Body.PageSize = res.PageSize
	out.Body.TotalPages = int((res.Total + int64(res.PageSize) - 1) / int64(res.PageSize))
	return out, nil
}

type ChangeRoleInput struct {
	ID   string `path:"id" format:"uuid"`
	Body struct {
		Role string `json:"role" enum:"requester,agent,admin"`
	}
}

type UserOutput struct{ Body AdminUserBody }

func (h *UserHandler) ChangeRole(ctx context.Context, in *ChangeRoleInput) (*UserOutput, error) {
	id, err := parseID(in.ID)
	if err != nil {
		return nil, err
	}
	u, err := h.svc.ChangeRole(ctx, id, models.Role(in.Body.Role))
	if err != nil {
		return nil, mapUserErr(err)
	}
	return &UserOutput{Body: toAdminUserBody(u)}, nil
}

type AssigneesOutput struct {
	Body struct {
		Items []UserBody `json:"items"`
	}
}

func (h *UserHandler) ListAssignees(ctx context.Context, _ *struct{}) (*AssigneesOutput, error) {
	us, err := h.svc.ListStaff(ctx)
	if err != nil {
		return nil, err
	}
	out := &AssigneesOutput{}
	out.Body.Items = make([]UserBody, 0, len(us))
	for i := range us {
		out.Body.Items = append(out.Body.Items, toUserBody(&us[i]))
	}
	return out, nil
}
