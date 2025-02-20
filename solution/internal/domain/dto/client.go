package dto

type CreateClientDTO struct {
	ClientId string `json:"client_id" validate:"required,uuid"`
	Login    string `json:"login" validate:"required,max=100"`
	Age      int    `json:"age" validate:"required,min=0,max=200"`
	Location string `json:"location" validate:"required,max=255"`
	Gender   string `json:"gender" validate:"required"`
}

type GetClientByIdDTO struct {
	ClientId string `params:"clientId"`
}
