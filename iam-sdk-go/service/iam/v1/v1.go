package v1

import "istomyang.github.com/like-iam/iam-sdk-go/pkg/client"

type Api interface {
	User() User
	Secret() Secret
	Policy() Policy
}

type apiV1 struct {
	client client.Client
}

func NewApiV1(client client.Client) Api {
	return &apiV1{client: client}
}

func (a *apiV1) User() User {
	return newUser(a.client)
}

func (a *apiV1) Secret() Secret {
	return newSecret(a.client)
}

func (a *apiV1) Policy() Policy {
	return newPolicy(a.client)
}

var _ Api = &apiV1{}
