package conn

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/olivere/elastic/v7"
	"istomyang.github.com/like-iam/component/pkg/options"
	"istomyang.github.com/like-iam/log"
	"net/http"
	"sync"
	"time"
)

type ElasticSearchClient struct {
	apiKey   string
	apiKeyId string
	ctx      context.Context
	cancel   context.CancelFunc
	opts     *options.ElasticSearchOpts

	esClient      *elastic.Client
	bulkProcessor *elastic.BulkProcessor
}

var (
	singletonES *ElasticSearchClient
	esMutex     sync.Mutex
)

func GetElasticSearchClient() *ElasticSearchClient {
	return singletonES
}

// NewElasticSearchClient returns if error, sleep seconds to retry.
func NewElasticSearchClient(ctx context.Context, opts *options.ElasticSearchOpts) (*ElasticSearchClient, error) {
	var client = &ElasticSearchClient{}
	client.ctx, client.cancel = context.WithCancel(ctx)
	singletonES = client
	if err := client.connect(); err != nil {
		return nil, err
	}
	return singletonES, nil
}

func (c *ElasticSearchClient) Client() *elastic.Client {
	return c.esClient
}

func (c *ElasticSearchClient) BulkProcessor() *elastic.BulkProcessor {
	return c.bulkProcessor
}

func (c *ElasticSearchClient) Run() error {
	var err error
	c.Client().Start()
	if err = c.bulkProcessor.Start(c.ctx); err != nil {
		return err
	}
	return nil
}

func (c *ElasticSearchClient) Close() error {
	var err error
	c.cancel()
	c.esClient.Stop()
	if err = c.bulkProcessor.Flush(); err != nil {
		return err
	}
	if err = c.bulkProcessor.Close(); err != nil {
		return err
	}
	return nil
}

// Reconnect use to for external package calls.
// If error, retry after seconds, default is 5s.
func (c *ElasticSearchClient) Reconnect() {
	if err := c.connect(); err != nil {
		time.Sleep(time.Second * 5)
		c.Reconnect()
	}
}

func (c *ElasticSearchClient) connect() error {
	esMutex.Lock()
	defer esMutex.Unlock()

	var opts = c.opts
	var client = c
	var httpClient *http.Client
	if opts.APIKey != "" && opts.APIKeyID != "" {
		httpClient = http.DefaultClient
		httpClient.Transport = client
	}

	// elastic.NewClientFromConfig use config but a little simple.
	esClient, err := elastic.NewClient(
		elastic.SetURL(opts.UrlsToSlice()...),
		elastic.SetBasicAuth(opts.Username, opts.Password),
		elastic.SetHttpClient(httpClient),
		elastic.SetSniff(opts.SniffEnabled),
		elastic.SetGzip(opts.GZip),
		elastic.SetErrorLog(client),
		elastic.SetInfoLog(client))
	if err != nil {
		return err
	}

	esBulker := esClient.
		BulkProcessor().
		Name(opts.Bulk.Name).
		Workers(opts.Bulk.Workers).
		BulkSize(opts.Bulk.BulkSize).
		BulkActions(opts.Bulk.BulkActions).
		FlushInterval(opts.Bulk.FlushInterval)

	if opts.Bulk.Name != "" {
		esBulker.Name("ES Bulk Processor")
	}

	client.esClient = esClient
	client.bulkProcessor, err = esBulker.Do(client.ctx)
	if err != nil {
		return err
	}
	return nil
}

var _ elastic.Logger = &ElasticSearchClient{}

func (c *ElasticSearchClient) Printf(format string, v ...interface{}) {
	log.Infof("ElasticSearchClient: "+format, v)
}

var _ http.RoundTripper = &ElasticSearchClient{}

func (c *ElasticSearchClient) RoundTrip(req *http.Request) (*http.Response, error) {
	tokenString := fmt.Sprintf("%s:%s", c.apiKeyId, c.apiKey)
	var token = make([]byte, base64.StdEncoding.EncodedLen(len(tokenString)))
	base64.StdEncoding.Encode(token, []byte(tokenString))
	h := fmt.Sprintf("ApiKey %s", token)

	req.Header.Set("Authorization", h)

	return http.DefaultClient.Transport.RoundTrip(req)
}
