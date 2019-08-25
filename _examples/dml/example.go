package main

import (
	"errors"
	"net/http"

	"github.com/g8rswimmer/go-sfdc"
	"github.com/g8rswimmer/go-sfdc/_examples/common"
	"github.com/g8rswimmer/go-sfdc/credentials"
	"github.com/g8rswimmer/go-sfdc/session"
)

func main() {
	payload := common.RetrievePayload()

	common.Title("Salesforce REST API DML Example")

	if payload.DML == nil {
		common.Error(errors.New("The payload does not contain the DML object"))
	}

	creds, err := credentials.NewPasswordCredentials(payload.Credentials)
	common.Error(err)

	config := sfdc.Configuration{
		Credentials: creds,
		Client:      http.DefaultClient,
		Version:     payload.Version,
	}

	session, err := session.Open(config)
	common.Error(err)

	payload.DML.Run(session)
}
