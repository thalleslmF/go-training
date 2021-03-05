package propose

import (
	"encoding/json"
	"fmt"
	"io"
	"training/internal/utils/validator"
)

type Propose struct {
	Email   string    `json:"email"`
	Name    string	    `json:"name"`
	Cpf     string		`json:"cpf"`
	Address string  `json:"Address"`
	Salary  int32   `json:"salary"`
}

type Request struct {
	Email   string	`json:"email"`
	Name    string		`json:"name"`
	Cpf     string 		`json:"cpf"`
	Address string  `json:"Address"`
	Salary  int32    `json:"salary"`
}

func (main Main) Create(proposeRequest Request ) error {
	tx := main.DB.Exec(
		fmt.Sprintf("INSERT INTO PROPOSE(email,name,cpf,address,salary) VALUES('%s','%s','%s','%s',%d)",
			proposeRequest.
			Email, proposeRequest.Name, proposeRequest.Cpf, proposeRequest.Address, proposeRequest.Salary,
		))
	if tx.Error != nil {
		return fmt.Errorf("Error: "+tx.Error.Error())
	}
	return nil
}

func (main Main) Validate(proposeRequest Request ) error {
	err := ValidateNonNull(proposeRequest)
	if err != nil {
		return err
	}
	err =  validator.ValidateCpf(proposeRequest.Cpf)
	if err != nil {
		return err
	}
	err= ValidateSalary(proposeRequest.Salary)
	if err != nil {
		return err
	}
	return nil
}


func (main Main) FindById(proposeId string) Propose {
	return Propose{}
}

func (main Main) ParsePropose(request io.ReadCloser)  (Request , error) {
	var proposeRequest Request
   	err := json.NewDecoder(request).Decode(&proposeRequest)
   	if err != nil {
		return Request{}, err
   	}
	return proposeRequest, nil
}

func ValidateNonNull(request Request ) error {
	if  request.Name == "" {
		return fmt.Errorf("empty name")
	}
	if  request.Address == "" {
		return fmt.Errorf("empty address")
	}
	if  request.Email == "" {
		return fmt.Errorf("empty email")
	}
	if  request.Cpf == "" {
		return fmt.Errorf("empty cpf")
	}
	return nil
}

func ValidateSalary(salary int32) error {
	if salary < 0 {
		return fmt.Errorf("Negative salaryf")
	}
	return nil
}