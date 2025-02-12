package dto

type CreateClientDTO struct {
	ClientId string `json:"client_id" validate:"required"`
	Login    string `json:"login" validate:"required"`
	Age      int    `json:"age" validate:"required,min=0,max=200"`
	Location string `json:"location" validate:"required"`
	Gender   string `json:"gender" validate:"required"`
}

type GetClientByIdDTO struct {
	ClientId string `params:"clientId"`
}
