package propose

import (
	"gorm.io/gorm"
	"io"
	"training/internal/card"
)

type ProposeUsecases interface {
	Create(propose Propose) error
	FindById(proposeId string) Propose
	ParsePropose(request io.ReadCloser) (Request, error)
	Validate(proposeRequest Request) error
	CheckIfUserAlreadyHasPropose(cpf string) error
	CheckIfProposeIsAvailable(propose Request) (ClientResponse, error)
	CheckIfCardWasGenerated(propose Propose)  (card.ClientCardResponse, error)
	GenerateCard(response ClientResponse)
}

type Main struct {
	DB *gorm.DB
}


func NewMain(db *gorm.DB) ProposeUsecases {
	return Main{db}
}