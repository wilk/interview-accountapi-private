package form3

import (
	"strconv"
	"encoding/json"
	"errors"
	guuid "github.com/google/uuid"
	"github.com/parnurzeal/gorequest"
)

var BASE_API_URL string

type accounts struct {}
// main wrapper: in the future, it could support other resources
type Form3 struct {
	Accounts accounts
}

type AccountAttributes struct {
	Country string `json:"country"`
	BaseCurrency string `json:"base_currency"`
	BankId string `json:"bank_id"`
	BankIdCode string `json:"bank_id_code"`
	AccountNumber string `json:"account_number"`
	Bic string `json:"bic"`
	Iban string `json:"iban"`
	CustomerId string `json:"customer_id"`
	Title string `json:"title"`
	FirstName string `json:"first_name"`
	BankAccountName string `json:"bank_account_name"`
	AlternativeBankAccountNames []string `json:"alternative_bank_account_names"`
	AccountClassification string `json:"account_classification"`
	JointAccount bool `json:"joint_account"`
	AccountMatchingOptOut bool `json:"account_matching_opt_out"`
	SecondaryIdentification string `json:"secondary_identification"`
}

type Account struct {
	Id string `json:"id"`
	OrganisationId string `json:"organisation_id"`
	Version int `json:"version"`
	Attributes AccountAttributes `json:"attributes"`
}

type AccountListPagination struct {
	Number int
	Size int
}

// list params: in the future it could support filtering too
type AccountListParams struct {
	Pagination AccountListPagination
}

// lib config: in the future it could support other config
type Form3Config struct {
	Url string
}

type accountResponseBody struct {
	Data Account `json:"data"`
}

type accountsResponseBody struct {
	Data []Account `json:"data"`
}

// generic error handler, both for HTTP responses and gorequest errors
func handleResponseError(res gorequest.Response, responseBody []byte, errs []error) error {
	if len(errs) > 0 {
		return errs[0]
	}

	if res.StatusCode > 399 {
		var responseMap map[string]string
		if err := json.Unmarshal(responseBody, &responseMap); err != nil {
			return err	
		}
		return errors.New(responseMap["error_message"])
	}

	return nil
}

func New(config *Form3Config) *Form3 {
	if config.Url != "" {
		BASE_API_URL = config.Url
	} else {
		// this should be converted into the prod URL
		BASE_API_URL = "http://localhost:8080/v1"
	}

	f := &Form3{}

	return f
}

func (a *accounts) Create(organisationId string, attributes *AccountAttributes) (Account, error) {
	var body accountResponseBody

	accountId := guuid.New().String()
	payload := map[string]interface{}{
		"data": map[string]interface{}{
			"id": accountId,
			"organisation_id": organisationId,
			"type": "accounts",
			"attributes": attributes,
		},
	}
	
	request := gorequest.New()
	res, responseBody, errs := request.Post(BASE_API_URL + "/organisation/accounts").
		Send(payload).
		EndStruct(&body)
	
	err := handleResponseError(res, responseBody, errs)
	
	return body.Data, err
}

func (a *accounts) List(params *AccountListParams) ([]Account, error) {
	var body accountsResponseBody

	query := map[string]int{
		"page[number]": params.Pagination.Number,
		"page[size]": params.Pagination.Size,
	}

	request := gorequest.New()
	res, responseBody, errs := request.Get(BASE_API_URL + "/organisation/accounts").
		Query(query).
		EndStruct(&body)
	
	err := handleResponseError(res, responseBody, errs)
	
	return body.Data, err
}

func (a *accounts) Fetch(accountId string) (Account, error) {
	var body accountResponseBody

	request := gorequest.New()
	res, responseBody, errs := request.Get(BASE_API_URL + "/organisation/accounts/" + accountId).
		EndStruct(&body)
	
	err := handleResponseError(res, responseBody, errs)

	return body.Data, err
}

func (a *accounts) Delete(accountId string, version int) error {
	request := gorequest.New()
	res, responseBody, errs := request.Delete(BASE_API_URL + "/organisation/accounts/" + accountId).
		Query("version=" + strconv.Itoa(version)).
		End()
	
	err := handleResponseError(res, []byte(responseBody), errs)

	return err
}
