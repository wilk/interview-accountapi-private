# Form3 Take Home Exercise

## Overview
This coding challenge has been made by Vincenzo Ferrari (aka Wilk) for the interview with Form3.

## Technical Decisions
### The client
The library consists in just one file (`form3.go`) plus the tests (`form3_test.go`).
It's a client library for the `accountapi` and it exposes one method called `New` used to initialise the lib with certain config: these config can then be extended in the future. For now, they're just used to bring the `API_URL`.

An example (taken from the test file):

```
client := form3.New(&form3.Form3Config{
  Url: os.Getenv("API_URL"),
})
```

The `Form3` client has no authentication method (it wasn't required) but it can easily extended with one and called this way: `client.Auth(credentials)`
Because the specs only required one resource (Account), `Form` has just one available resource with the 4 required actions (Create, Fetch, List, Delete).

An example (taken from the test file):

```
account, err := client.Accounts.Create(ORGANIZATION_ID, attributes)
```

`client.Accounts` is the way someone can access to the resource and `Create` is the action required.
All the actions follow the common error handling style guide: `account, err` if an error occurred, err will be filled with it, otherwise it will be `nil`.

The client uses a couple of external lib:

- github.com/google/uuid
- github.com/parnurzeal/gorequest

They are included to speed up the development.

The client is simple and it does not validate the input: it just bubbles up the error coming from the APIs (or from `gorequest`) towards the user.
So, the user needs to know what kind of errors the API can throw.

Golang structs have been preferred on maps because of the type safetyness.

The `golang mod` has been used to define the module and to save the deps with the right versioning.

### The tests
The tests are E2E.
They use the `Form3` client to perform full API requests and then the response is checked against a set of assertions.
The tests are isolated and indipendent from each other.
To achieve this I had to perform different calls to the API in the same test, like the following example (taken from the tests):

```
newAccount, _ := client.Accounts.Create(ORGANIZATION_ID, attributes)
account, err := client.Accounts.Fetch(newAccount.Id)
assert.Equal(newAccount.Id, account.Id)
assert.Equal(newAccount.Version, account.Version)
assert.Equal(newAccount.OrganisationId, account.OrganisationId)
```

Usually, this is done using `fixtures`.
However, I didn't want to deep dive down the database schema so I used the API to populate it step by step.

The tests contain both expected and unexpected cases, like `TestAccounts_Fetch` and `TestAccounts_FetchFakeAccount`.

Finally, two files have been used, with different packages, like `form3` and `form3_test`, and that's because I wanted to test the client as if I'd use it as a final user.

### The Docker Image
Everything starts within a custom docker image, created using the `Dockerfile` in this repo.
I used a full image (instead of an alpine, for instance) because of `cgo` used within some external deps.
Also, I've had to install `netcat` to sync the tests with the `accountapi` bootstrapping phase.

### Run it (and enjoy!)
Of course you can use `docker-compose up` to run it but I'd suggest to use the `--abort-on-container-exit` flag, thus when the tests finish, all the infrastructure comes down.

---- ORIGINAL SPECS ----

## Instructions

This exercise has been designed to be completed in 4-8 hours. The goal of this exercise is to write a client library 
in Go to access our fake [account API](http://api-docs.form3.tech/api.html#organisation-accounts) service. 

### Should
- Client library should be written in Go
- Document your technical decisions
- Implement the `Create`, `Fetch`, `List` and `Delete` operations on the `accounts` resource. Note that filtering of the List operation is not required, but you should support paging
- Focus on writing full-stack tests that cover the full range of expected and unexpected use-cases
 - Tests can be written in Go idiomatic style or in BDD style. Form3 engineers tend to favour BDD. Make sure tests are easy to read
 - If you encounter any problems running the fake accountapi we would encourage you to do some debugging first, 
before reaching out for help

#### Docker-compose

 - Add your solution to the provided docker-compose file
 - We should be able to run `docker-compose up` and see your tests run against the provided account API service 

### Please don't
- Use a code generator to write the client library
- Implement an authentication scheme

## How to submit your exercise
- Create a private [GitHub](https://help.github.com/en/articles/create-a-repo) repository, copy the `docker-compose` from this repository
- [Invite](https://help.github.com/en/articles/inviting-collaborators-to-a-personal-repository) @form3tech-interviewer-1 and @form3tech-interviewer-2 to your private repo
- Let us know you've completed the exercise using the link provided at the bottom of the email from our recruitment team
- Usernames of the developers reviewing your code will then be provided for you to grant them access to your private repository
- Put your name in the README
