package client

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"io"
	"istomyang.github.com/like-iam/component-base/base"
	"istomyang.github.com/like-iam/component-base/web"
	"istomyang.github.com/like-iam/iam-sdk-go/pkg/util/coder"
	"net/http"
	"net/url"
	"runtime"
	"strings"
)

type request struct {
	coder base.Coder

	certPEM  []byte
	keyPEM   []byte
	caPEM    []byte
	insecure bool

	// verb is Method of http.Request, required.
	verb string

	// resource is rest spec's resource, required.
	resource string

	// version is api version, like v1 or v2, required.
	version string

	// name is username of iam.
	name string

	// params is http spec's param.
	params url.Values

	// meta is operation meta.
	meta *url.Values

	// action is change-password of /user/change-password or /login or /logout
	action string

	// header is http spec's header.
	header http.Header

	// body is http spec's body.
	body []byte

	// u contains baseUrl.
	u *url.URL
}

// newRequest advises on using Client to create Request.
func newRequest(u *url.URL, coderName string, cert, key, ca []byte, insecure bool) Request {
	var r request
	r.u = u
	r.coder = coder.Get(coderName)
	r.certPEM = cert
	r.keyPEM = key
	r.caPEM = ca
	r.insecure = insecure
	return &r
}

var _ Request = &request{}

func (r *request) Verb(v Verb) Request {
	r.verb = string(v)
	return r
}

func (r *request) Resource(s Res) Request {
	r.resource = string(s)
	return r
}

func (r *request) Name(s string) Request {
	r.name = s
	return r
}

func (r *request) Action(s string) Request {
	r.action = s
	return r
}

func (r *request) Params(p url.Values) Request {
	r.params = p
	return r
}

func (r *request) Version(v V) Request {
	r.version = string(v)
	return r
}

func (r *request) Meta(m any) Request {
	if r.meta == nil {
		r.meta = &url.Values{}
	}
	u, _ := web.UrlValues(m)
	for k, _ := range map[string][]string(*u) {
		r.meta.Set(k, u.Get(k))
	}
	return r
}

func (r *request) Header(k, v string) Request {
	if r.header == nil {
		r.header = http.Header{}
	}
	r.header.Set(k, v)
	return r
}

func (r *request) Body(v any) Request {
	var err error
	if r.body, err = r.coder.Encode(v); err != nil {
		// Panic directly.
		panic(fmt.Errorf("request's body can't encode, print it: %v", v))
	}
	return r
}

func (r *request) Send(ctx context.Context) Response {
	build := r.build1
	send := r.send1

	req, err := build(ctx)
	if err != nil {
		return newResponse(nil, r.coder, err)
	}
	res, err := send(ctx, req)
	if err != nil {
		return newResponse(nil, r.coder, err)
	}
	return newResponse(res, r.coder, nil)
}

func (r *request) build1(ctx context.Context) (req *http.Request, err error) {
	req, err = http.NewRequestWithContext(ctx, r.verb, r.buildUrl(), r.buildBody())
	{
		req.Header = r.header
		req.Header.Set("Accept", "application/json")
		req.Header.Set("User-Agent", r.buildUA())

		// Client should tell server body's mime type.
		if len(r.body) > 0 && (r.verb == string(VerbPost) || r.verb == string(VerbPUT)) {
			req.Header.Set("Content-Type", "application/json")
		}
	}
	return
}

func (r *request) send1(ctx context.Context, req *http.Request) (*http.Response, error) {
	var err error
	var httpClient = &(*http.DefaultClient)
	httpClient.Transport, err = r.buildRoundTripper()
	{
		req.WithContext(ctx)
	}
	res, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// buildUrl assumes all field is right.
func (r *request) buildUrl() string {
	var result strings.Builder
	{
		// Established by popular usage
		// Here result: https://localhost:443/login
		var a = r.action == "login"
		var b = r.action == "logout"
		var c = r.action == "refresh"
		if a || b || c {
			result.WriteString(r.u.Scheme)
			result.WriteString("//")
			result.WriteString(r.u.Host)
			result.WriteString("/")
			result.WriteString(r.action)
			return result.String()
		}
	}
	{
		// Here result: https://localhost:443/v1/users
		result.WriteString(r.u.Scheme)
		result.WriteString("//")
		result.WriteString(r.u.Host)
		result.WriteString("/")
		result.WriteString(r.version)
		result.WriteString("/")
		result.WriteString(r.resource)
	}
	{
		if r.name != "" {
			result.WriteString("/")
			result.WriteString(r.name)
		}
		if r.action != "" {
			result.WriteString("/")
			result.WriteString(r.action)
		}
		if r.params != nil {
			result.WriteString("?")
			var uv = url.Values{}
			for k, _ := range map[string][]string(r.params) {
				uv.Set(k, r.params.Get(k))
			}
			for k, _ := range map[string][]string(*r.meta) {
				uv.Set(k, r.meta.Get(k))
			}
			result.WriteString(uv.Encode())
		}
	}

	return result.String()
}

func (r *request) buildBody() io.Reader {
	if r.body == nil {
		return nil
	}
	var buf bytes.Buffer
	buf.Write(r.body)
	return &buf
}

func (r *request) buildUA() string {

	var result bytes.Buffer
	{
		s := runtime.Version()                       // go1.19.2
		g := cases.Title(language.Und).String(s[:2]) // Go
		v := s[2:]                                   // 1.19.2
		result.WriteString(g)
		result.WriteString("/")
		result.WriteString(v)

		result.WriteString(" ")

		result.WriteString("(")
		switch runtime.GOOS {
		case "darwin":
			result.WriteString("Darwin; MacOS;")
		case "linux":
			result.WriteString("Linux;")
		case "freebsd":
			result.WriteString("FreeBSD;")
		default:
			result.WriteString("Other Systems;")
		}
		result.WriteString(runtime.GOARCH) // amd64
		result.WriteString(")")

		// Result: Go/1.19.2 (Darwin; MacOS; amd64)
	}
	result.WriteString(" ")
	{
		result.WriteString("IAM-SDK-GO/1.0 ")
		result.WriteString("(")
		result.WriteString("api:" + r.version)
		result.WriteString(")")

		// Result: IAM-SDK-GO/1.0 (api:v1)
	}

	// Result: Go/1.19.2 (Darwin; MacOS; amd64) IAM-SDK-GO/1.0 (api:v1)

	return result.String()
}

// buildRoundTripper handle tls config.
func (r *request) buildRoundTripper() (http.RoundTripper, error) {
	var ts = http.DefaultTransport.(*http.Transport).Clone()
	var tlsConfig = ts.TLSClientConfig
	var err error
	{
		var pool = tlsConfig.RootCAs
		if pool == nil {
			if pool, err = x509.SystemCertPool(); err != nil {
				return nil, err
			}
		}
		if r.caPEM != nil {
			pool.AppendCertsFromPEM(r.caPEM)
		}
		tlsConfig.RootCAs = pool
	}
	{
		var cert tls.Certificate
		if cert, err = tls.X509KeyPair(r.certPEM, r.keyPEM); err != nil {
			return nil, err
		}
		tlsConfig.Certificates = []tls.Certificate{cert}
	}
	tlsConfig.InsecureSkipVerify = r.insecure
	ts.TLSClientConfig = tlsConfig

	return ts, err
}
