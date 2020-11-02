package http

import (
	"github.com/jmoiron/sqlx"
	"goacs/models/cpe"
	"goacs/repository"
	"log"
	"net/http"
	"strings"
	"time"
)

type ACSRequest struct {
	DBConnection *sqlx.DB
	CPE          *cpe.CPE
	Body         string
	Response     *http.Response
}

func NewACSRequest(cpe *cpe.CPE) *ACSRequest {
	return &ACSRequest{
		DBConnection: repository.GetConnection(),
		CPE:          cpe,
		Body:         "",
	}
}

func (acsRequest *ACSRequest) Kick() {
	err := acsRequest.Send()

	if err != nil {
		log.Println("Kick error", err.Error())
	}
}

func (acsRequest *ACSRequest) Send() error {
	request, err := http.NewRequest("GET", acsRequest.CPE.ConnectionRequestUrl, strings.NewReader(acsRequest.Body))

	if err != nil {
		log.Println("Request send error", err.Error())

		return err
	}

	request.SetBasicAuth(acsRequest.CPE.ConnectionRequestUser, acsRequest.CPE.ConnectionRequestPassword)

	client := http.Client{
		Timeout: time.Second * 5,
	}

	response, err := client.Do(request)

	if err != nil {
		return err
	}

	defer response.Body.Close()

	acsRequest.Response = response

	return nil

}
