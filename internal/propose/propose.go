package propose

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"net/http"
	"training/internal/card"
	"training/internal/utils/validator"
)

type Propose struct {
	Email       string    	`json:"email"`
	Name        string	    `json:"name"`
	Cpf     	string		`json:"cpf"`
    Address 	string  	`json:"Address"`
	Salary 		int32   	`json:"salary"`
	IdPropose 	string 		`json:"idPropose"`
	State 	string 			`json:"state"`
}

type Request struct {
	Email                string		`json:"email"`
	Name                 string		`json:"name"`
	Cpf                  string 	 `json:"cpf"`
	Address              string  	`json:"Address"`
	Salary               int32    	`json:"salary"`
}

type ClientRequest struct {
	Documento       string			`json:"documento"`
	Nome     		string 			`json:"nome"`
	IdProposta      string 			`json:"idProposta"`
}

type ClientResponse struct {
	Documento            string			    `json:"documento"`
	Nome     		     string 			`json:"nome"`
	IdProposta           string 			`json:"idProposta"`
	ResultadoSolicitacao string				`json:"resultadoSolicitacao"`
}

func (main Main) Create(propose Propose ) error {
	tx := main.DB.Exec(
		fmt.Sprintf("INSERT INTO PROPOSE(email,name,cpf,address,salary, id, state) VALUES('%s','%s','%s','%s',%d, '%s', '%s')",
			propose.
			Email, propose.Name, propose.Cpf, propose.Address, propose.Salary, propose.IdPropose,propose.State,
		))
	if tx.Error != nil {
		return fmt.Errorf("Error: "+tx.Error.Error())
	}
	return nil
}
func (main Main) GenerateCard(response ClientResponse ) {
	client := http.Client{}
	logrus.Info("Creating response")
	responseBytes,_ := json.Marshal(response)
	req, err := http.NewRequest("POST", "http://localhost:8888/api/cartoes", bytes.NewReader(responseBytes))
	if err != nil {
		logrus.Error(fmt.Errorf("Error creating request", err.Error()))
	}
	req.Header.Add("Content-Type","application/json")
	res, err := client.Do(req)
	if err != nil {
		logrus.Error(fmt.Errorf("Error making request", err.Error()))
	}
	responseBytes, _ = ioutil.ReadAll(res.Body)
	if res.StatusCode != 201 {
		logrus.Error(fmt.Errorf("Error making request", req.URL));
		logrus.Error(fmt.Errorf("Error making request", res.StatusCode));
		logrus.Error(fmt.Errorf("Error making request", string(responseBytes)));
	}
	defer res.Body.Close()
}
func (main Main) CheckIfCardWasGenerated(propose Propose) error {
	client := http.Client{}
	cardResponse := card.ClientCardResponse{}
	logrus.Info("Creating response")
	req, err := http.NewRequest("GET", "http://localhost:8888/api/cartoes", nil)
	if err != nil {
		return fmt.Errorf("Error creating request", err.Error())
	}

	req.Header.Add("Content-Type","application/json")
	q := req.URL.Query()
	q.Add("idProposta", propose.IdPropose)
	req.URL.RawQuery = q.Encode()
	logrus.Info("Making request")
	response,err := client.Do(req)
	if err  != nil {
		return fmt.Errorf("Error making request", err.Error())
	}

	err = json.NewDecoder(response.Body).Decode(&cardResponse)
	if err != nil {

		return fmt.Errorf("Error parsing response", err.Error())
	}
	defer response.Body.Close()
	logrus.Info("Response status code ", response.StatusCode)
	logrus.Info("Req url ", req.URL.String())
	logrus.Info("Response body ",cardResponse)

	return nil
}

func (main Main) CheckIfUserAlreadyHasPropose(cpf string) error {
	var count int32
	row, err  := main.DB.Raw(
		fmt.Sprintf("SELECT count(*) FROM PROPOSE p where p.cpf = '%s'",
			cpf,
		)).Rows()
	if err != nil {
		return fmt.Errorf("error finding propose %s", row.Err())
	}
	if row.Next() {
		_ = row.Scan(&count)
		if count > 0 {
			return fmt.Errorf("Propose already registered" )
		}
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

func (main Main) CheckIfProposeIsAvailable(propose Request) (ClientResponse, error) {

	client := &http.Client{
	}
	clientResponse := ClientResponse{
	}
	logrus.Println("Generating uuid", propose)
	uuid, err := uuid.NewUUID()
	stringId, err := uuid.MarshalText()
	if err != nil {
		return ClientResponse{}, fmt.Errorf("Error generating uuid")
	}
	var clientRequest = ClientRequest{
		Nome: propose.Name,
		Documento: propose.Cpf,
		IdProposta: string(stringId),
	}

	byteRequest, err :=json.Marshal(clientRequest)
	logrus.Println("Generating client request", string(byteRequest))
	if err != nil {
		return  ClientResponse{}, fmt.Errorf("error parsing client request", err.Error())
	}
	logrus.Println("Making request", clientRequest)

	req, err := http.NewRequest("POST", "http://localhost:9999/api/solicitacao", bytes.NewReader(byteRequest))

	if err != nil {
		return  ClientResponse{}, fmt.Errorf("Error creating client request", err.Error())

	}

	req.Header.Set("Content-type","application/json")
	response, err := client.Do(req)
	if err != nil {
		return  ClientResponse{}, fmt.Errorf("Error on client request", err.Error())
	}

	err = json.NewDecoder(response.Body).Decode(&clientResponse)
	logrus.Println("Response parsed", clientResponse)
	if err != nil {
		return ClientResponse{}, fmt.Errorf("Error reading response", err.Error())
	}
	defer response.Body.Close()
	return  clientResponse, nil
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

func (p ClientResponse) ToPropose(r Request) Propose {
	return Propose {
		IdPropose: p.IdProposta,
		Name: r.Name,
		Cpf: r.Cpf,
		Salary: r.Salary,
		Address: r.Address,
		Email: r.Email,
		State: p.ResultadoSolicitacao,
	}
}