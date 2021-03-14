package propose

import (
	"gorm.io/gorm"
	"io"
)

type ProposeUsecases interface {
	Create(propose Propose) error
	FindById(proposeId string) Propose
	ParsePropose(request io.ReadCloser) (Request, error)
	Validate(proposeRequest Request) error
	CheckIfUserAlreadyHasPropose(cpf string) error
	CheckIfProposeIsAvailable(propose Request) (ClientResponse, error)
	CheckIfCardWasGenerated(propose Propose)  error
	GenerateCard(response ClientResponse)
}

type Main struct {
	DB *gorm.DB
}


func NewMain(db *gorm.DB) ProposeUsecases {
	return Main{db}
}