package form3_test

import (
	"os"
	"testing"
	"sync"
	"github.com/stretchr/testify/assert"
	guuid "github.com/google/uuid"
	"form3"
)

const ORGANIZATION_ID = "eb0bd6f5-c3f5-44b2-b677-acd23cdde73c"

func buildAccountAttributes() *form3.AccountAttributes {
	var names []string
	return &form3.AccountAttributes{
		Country: "GB",
		BaseCurrency: "GBP",
		BankId: "400302",
		BankIdCode: "GBDSC",
		AccountNumber: "10000004",
		CustomerId: "234",
		Iban: "GB28NWBK40030212764204",
		AlternativeBankAccountNames: names,
    Bic: "NWBKGB42",
    Title: "Sir",
  	FirstName: "Mary-Jane Doe",
    AccountClassification: "Personal",
    JointAccount: false,
    AccountMatchingOptOut: false,
	}
}

func TestAccounts_Create(t *testing.T) {
	assert := assert.New(t)
	client := form3.New(&form3.Form3Config{
		Url: os.Getenv("API_URL"),
	})
	attributes := buildAccountAttributes()
	
	account, err := client.Accounts.Create(ORGANIZATION_ID, attributes)
	assert.Nil(err)
	assert.NotNil(account.Id)
	assert.Equal(0, account.Version)
	assert.Equal(ORGANIZATION_ID, account.OrganisationId)
	assert.Equal(attributes.Country, account.Attributes.Country)
	assert.Equal(attributes.BaseCurrency, account.Attributes.BaseCurrency)
	assert.Equal(attributes.BankId, account.Attributes.BankId)
	assert.Equal(attributes.BankIdCode, account.Attributes.BankIdCode)
	assert.Equal(attributes.AccountNumber, account.Attributes.AccountNumber)
	assert.Equal(attributes.CustomerId, account.Attributes.CustomerId)
	assert.Equal(attributes.Iban, account.Attributes.Iban)
	assert.Equal(attributes.Bic, account.Attributes.Bic)
	assert.Equal(attributes.Title, account.Attributes.Title)
	assert.Equal(attributes.FirstName, account.Attributes.FirstName)
	assert.Equal(attributes.AccountClassification, account.Attributes.AccountClassification)
	assert.Equal(attributes.JointAccount, account.Attributes.JointAccount)
	assert.Equal(attributes.AccountMatchingOptOut, account.Attributes.AccountMatchingOptOut)
}

func TestAccounts_CreateWrongPayload(t *testing.T) {
	assert := assert.New(t)
	client := form3.New(&form3.Form3Config{
		Url: os.Getenv("API_URL"),
	})
	attributes := buildAccountAttributes()
	attributes.Country = "G"
	
	_, err := client.Accounts.Create(ORGANIZATION_ID, attributes)
	assert.NotNil(err)
	assert.Equal("validation failure list:\nvalidation failure list:\nvalidation failure list:\ncountry in body should match '^[A-Z]{2}$'", err.Error())
}

func TestAccounts_List(t *testing.T) {
	assert := assert.New(t)
	client := form3.New(&form3.Form3Config{
		Url: os.Getenv("API_URL"),
	})
	attributes := buildAccountAttributes()

	lastAccount, _ := client.Accounts.Create(ORGANIZATION_ID, attributes)
	accounts, err := client.Accounts.List(&form3.AccountListParams{})
	assert.Nil(err)
	
	if (len(accounts) == 0) {
		t.Fatalf("expected to have at least one account but got 0")
	}

	assert.Equal(lastAccount.Id, accounts[len(accounts) - 1].Id)
	assert.Equal(lastAccount.Version, accounts[len(accounts) - 1].Version)
}

func TestAccounts_ListPaginate(t *testing.T) {
	assert := assert.New(t)
	client := form3.New(&form3.Form3Config{
		Url: os.Getenv("API_URL"),
	})
	attributes := buildAccountAttributes()

	wg := &sync.WaitGroup{}
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func(){
			client.Accounts.Create(ORGANIZATION_ID, attributes)
			wg.Done()
		}()
	}
	wg.Wait()

	accounts, err := client.Accounts.List(&form3.AccountListParams{
		form3.AccountListPagination{
			Number: 2,
			Size: 2,
		},
	})
	assert.Nil(err)
	assert.Equal(2, len(accounts))
	
	otherAccounts, err := client.Accounts.List(&form3.AccountListParams{
		form3.AccountListPagination{
			Number: 3,
			Size: 2,
		},
	})
	assert.Nil(err)
	assert.Equal(2, len(otherAccounts))

	assert.NotEqual(accounts[0].Id, otherAccounts[0].Id)
	assert.NotEqual(accounts[1].Id, otherAccounts[1].Id)

	tooFar, err := client.Accounts.List(&form3.AccountListParams{
		form3.AccountListPagination{
			Number: 3,
			Size: 20000,
		},
	})
	assert.Nil(err)
	assert.Equal(0, len(tooFar))
}

func TestAccounts_Fetch(t *testing.T) {
	assert := assert.New(t)
	client := form3.New(&form3.Form3Config{
		Url: os.Getenv("API_URL"),
	})

	attributes := buildAccountAttributes()
	newAccount, _ := client.Accounts.Create(ORGANIZATION_ID, attributes)

	account, err := client.Accounts.Fetch(newAccount.Id)
	assert.Nil(err)
	assert.Equal(newAccount.Id, account.Id)
	assert.Equal(newAccount.Version, account.Version)
	assert.Equal(newAccount.OrganisationId, account.OrganisationId)
}

func TestAccounts_FetchFakeAccount(t *testing.T) {
	assert := assert.New(t)
	client := form3.New(&form3.Form3Config{
		Url: os.Getenv("API_URL"),
	})
	fakeAccountId := guuid.New().String()

	account, err := client.Accounts.Fetch(fakeAccountId)
	assert.NotNil(err)
	assert.NotNil(account)
	assert.Equal("", account.Id)
	assert.Equal("record " + fakeAccountId + " does not exist", err.Error())
}

func TestAccounts_Delete(t *testing.T) {
	assert := assert.New(t)
	client := form3.New(&form3.Form3Config{
		Url: os.Getenv("API_URL"),
	})

	attributes := buildAccountAttributes()
	newAccount, _ := client.Accounts.Create(ORGANIZATION_ID, attributes)

	err := client.Accounts.Delete(newAccount.Id, newAccount.Version)
	assert.Nil(err)

	nonExistingAccount, err := client.Accounts.Fetch(newAccount.Id)
	assert.NotNil(err)
	assert.NotNil(nonExistingAccount)
	assert.NotEqual(nonExistingAccount.Id, newAccount.Id)
	assert.Equal("", nonExistingAccount.Id)
}
