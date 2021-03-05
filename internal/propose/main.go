package propose

import (
	"gorm.io/gorm"
	"io"
)

type ProposeUsecases interface {
	Create(proposeRequest Request) error
	FindById(proposeId string) Propose
	ParsePropose(request io.ReadCloser) (Request, error)
	Validate(proposeRequest Request) error
}

type Main struct {
	DB *gorm.DB
}


func NewMain(db *gorm.DB) ProposeUsecases {
	return Main{db}
}