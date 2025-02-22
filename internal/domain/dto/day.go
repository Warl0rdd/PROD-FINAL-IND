package dto

type SetDayDTO struct {
	CurrentDate int `json:"current_date"  validate:"required"`
}
