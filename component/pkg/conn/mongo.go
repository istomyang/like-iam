package conn

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	mongoOptions "go.mongodb.org/mongo-driver/mongo/options"
	"istomyang.github.com/like-iam/component/pkg/interfaces"
	"istomyang.github.com/like-iam/component/pkg/options"
	"strings"
	"time"
)

type MongoClient struct {
	ctx        context.Context
	cancel     context.CancelFunc
	config     *options.MongoOpts
	client     *mongo.Client
	collection *mongo.Collection

	retry        int
	defaultRetry int
}

func NewMongoClient(ctx context.Context, opts *options.MongoOpts) *MongoClient {
	var c MongoClient
	c.ctx, c.cancel = context.WithCancel(ctx)
	c.config = opts
	c.defaultRetry = 5
	return &c
}

func (m *MongoClient) Connect() error {
	return m.connect(nil)
}

func (m *MongoClient) connect(e error) error {
	if m.retry > m.defaultRetry {
		return fmt.Errorf("can't connect to mongodb, got error: %s", e.Error())
	}
	var err error
	var opts = mongoOptions.Client().ApplyURI(m.buildUrl())
	if m.config.CAFile != "" && m.config.KeyFile != "" {
		credential := mongoOptions.Credential{
			AuthMechanism: "MONGODB-X509",
		}
		opts = opts.SetAuth(credential)
	}
	m.client, err = mongo.Connect(m.ctx, opts)
	if err != nil {
		m.retry++
		time.Sleep(time.Second * 5)
		return m.connect(err)
	}
	return nil
}

func (m *MongoClient) buildUrl() string {
	var c = m.config
	var url strings.Builder
	url.WriteString("mongodb://")
	url.WriteString(c.Addr)
	if m.config.CAFile != "" && m.config.KeyFile != "" {
		url.WriteString("?")
		url.WriteString("tlsCAFile")
		url.WriteString("=")
		url.WriteString(m.config.CAFile)
		url.WriteString("tlsCertificateKeyFile")
		url.WriteString("=")
		url.WriteString(m.config.KeyFile)
	}
	return url.String()
}

func (m *MongoClient) Send(dataset []map[string]any) error {
	var ddd = make([]any, len(dataset))
	for i, data := range dataset {
		var d bson.D
		for k, v := range data {
			d = append(d, bson.E{Key: k, Value: v})
		}
		ddd[i] = d
	}
	_, err := m.collection.InsertMany(m.ctx, ddd)
	return err
}

func (m *MongoClient) Run() (err error) {
	err = m.Connect()
	m.collection = m.client.Database(m.config.DbName).Collection(m.config.CollectionName)
	return
}

func (m *MongoClient) Close() error {
	m.cancel()
	return m.client.Disconnect(context.TODO())
}

var _ interfaces.ComponentCommon = &MongoClient{}
