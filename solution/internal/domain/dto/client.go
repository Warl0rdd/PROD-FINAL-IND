package dto

type CreateClientDTO struct {
	ClientId string `json:"client_id" validate:"required,uuid"`
	Login    string `json:"login" validate:"required,email"`
	Age      int    `json:"age" validate:"required"`
	Location string `json:"location" validate:"required"`
	Gender   string `json:"gender" validate:"required"`
}

type GetClientByIdDTO struct {
	ClientId string `params:"clientId"`
}
