package propose

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"training/internal/card"
	"training/internal/utils/http"
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
	var clientResponseBytes,_ = json.Marshal(response)

	headersMap := make(map[string]string)
	headersMap["Content-Type"] = "application/json"
	responseRequest, err := http.MakeRequest(
		"POST",
		"http://localhost:8888/api/cartoes",
		clientResponseBytes, headersMap, nil )
	if err != nil {
		 fmt.Errorf("Error parsing response", err.Error())
	}
	responseBytes, _ := ioutil.ReadAll(responseRequest.Body)
	defer responseRequest.Body.Close()
	if responseRequest.StatusCode != 201 {
		logrus.Error(fmt.Errorf("Error making request", string(responseBytes)));
	}
}
func (main Main) CheckIfCardWasGenerated(propose Propose) (card.ClientCardResponse, error) {
	headersMap := make(map[string]string)
	paramsMap := make(map[string]string)
	cardResponse := card.ClientCardResponse{}
	headersMap["Content-Type"] = "application/json"
	paramsMap["idProposta"] = propose.IdPropose

	response, err := http.MakeRequest("GET", "http://localhost:8888/api/cartoes", nil, headersMap, paramsMap)
	defer response.Body.Close()
	if err  != nil {
		return cardResponse, err
	}
	err = json.NewDecoder(response.Body).Decode(&cardResponse)


	if err != nil {

		return cardResponse, fmt.Errorf("Error parsing response", err.Error())
	}
	logrus.Println("card generated", cardResponse)
	return cardResponse, nil
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
	clientResponse := ClientResponse{}

	headersMap := make(map[string]string)
	headersMap["Content-Type"] = "application/json"
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
	if err != nil {
		return  ClientResponse{}, fmt.Errorf("error parsing client request", err.Error())
	}

	response, err := http.MakeRequest(
		"POST",
		"http://localhost:9999/api/solicitacao",
		byteRequest,
		headersMap,
		nil)

	if err != nil {
		return  ClientResponse{}, fmt.Errorf("Error creating client request", err.Error())

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