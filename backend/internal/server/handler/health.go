package handler

import "context"

type HealthOutput struct {
	Body struct {
		Status string `json:"status" example:"ok"`
	}
}

func Health(ctx context.Context, _ *struct{}) (*HealthOutput, error) {
	out := &HealthOutput{}
	out.Body.Status = "ok"
	return out, nil
}
