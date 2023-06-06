package models

type LoginRequest struct {
	Email    string
	Password string
}

type CreateMembershipRequest struct {
	UserID int `json:"user_id"`
	PlanID int `json:"plan_id"`
}
