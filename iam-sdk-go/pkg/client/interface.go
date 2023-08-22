package client

import (
	"context"
	"net/http"
	"net/url"
)

type Verb string

// Verbs depends on routers of apiserver for constraint.
const (
	VerbGET    Verb = http.MethodGet
	VerbPUT    Verb = http.MethodPut
	VerbPost   Verb = http.MethodPost
	VerbDelete Verb = http.MethodDelete
)

type Res string

// ResXXX depends on routers of apiserver for constraint.
const (
	ResUser   Res = "users"
	ResPolicy Res = "policies"
	ResSecret Res = "secrets"
)

type V string

// Vx is Version of api for constraint.
const (
	V1 V = "v1"
	V2 V = "v2"
)

// Client defines interface of rest client.
//
// It is recommended to use Verb, Get, Post, Put, Delete to create Request.
// Generally Client use Request to get Response, but in this package, I use chain calls, you can do
// like that: `DefaultClient.Post().Resource(ResUser).Body(&new_user).Send(ctx).Error()
type Client interface {
	Verb(method Verb) Request
	Get() Request
	Post() Request
	Put() Request
	Delete() Request
}

// Request defines interface of rest request which its api nears service layer.
type Request interface {
	// Verb represents a method for http, like VerbGET, VerbPUT and VerbPost.
	Verb(v Verb) Request
	// Resource represents resource group for api, like ResUser, ResPolicy and ResSecret.
	Resource(s Res) Request
	// Name is a string of username.
	Name(s string) Request
	// Action is action for resource, in established by popular usage, `change-password` in /user/change-password, /login, /logout and so on.
	Action(s string) Request
	// Params is an array like "?key1=value1&key2=value2".
	Params(p url.Values) Request
	// Version sets api 's version, use V1 or V2.
	Version(v V) Request
	// Header sets headers to request.
	Header(k, v string) Request
	// Meta sets meta like GetOperateMeta, DeleteOperateMeta and so on.
	Meta(m any) Request

	// Body is http request's body.
	Body(v any) Request

	// Send does the real request work.
	Send(ctx context.Context) Response
}

// Response defines interface of result for response.
type Response interface {
	Into(v any) error
	Raw() ([]byte, error)
	// Error returns server error like 500 and other errors.
	Error() error
}
