package elasticsearch

import (
	"context"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/olivere/elastic/v7"
	"istomyang.github.com/like-iam/component/pkg/conn"
	"istomyang.github.com/like-iam/component/pkg/options"
	"istomyang.github.com/like-iam/iam/internal/pkg/analytics"
	"istomyang.github.com/like-iam/iam/internal/pumper/pumps"
	"istomyang.github.com/like-iam/log"
	"time"
)

type esPump struct {
	ctx    context.Context
	cancel context.CancelFunc

	timeout    time.Duration
	omitDetail bool
	filter     *pumps.Filter

	config *options.ElasticSearchOpts
	client *conn.ElasticSearchClient
}

func New() pumps.Pump {
	var cp = esPump{}
	return &cp
}

func (e *esPump) Init(ctx context.Context, config map[string]any) error {
	var err error
	e.config = &options.ElasticSearchOpts{}
	if err = mapstructure.Decode(config, e.config); err != nil {
		return err
	}

	if e.config.Urls == "" {
		e.config.Urls = "http://localhost:9200"
	}
	if e.config.IndexName == "" {
		e.config.IndexName = "iam_analytics"
	}
	if e.config.DocumentType == "" {
		e.config.DocumentType = "iam_analytics"
	}
	log.Infof("Elasticsearch URLs: %s", e.config.Urls)
	log.Infof("Elasticsearch Index name: %s", e.config.IndexName)
	log.Infof("Elasticsearch DocumentType: %s", e.config.DocumentType)

	if e.config.RollingIndex {
		log.Infof("Index will have date appended to it in the format %s -YYYY.MM.DD", e.config.IndexName)
	}

	e.ctx, e.cancel = context.WithCancel(ctx)
	if e.client, err = conn.NewElasticSearchClient(e.ctx, e.config); err != nil {
		return err
	}
	return nil
}

func (e *esPump) Run() error {
	var err error
	if err = e.client.Run(); err != nil {
		return err
	}
	return nil
}

func (e *esPump) Close() error {
	var err error
	e.cancel()
	if err = e.client.Close(); err != nil {
		return err
	}
	return nil
}

func (e *esPump) GetName() string {
	return "ElasticSearch Pump"
}

func (e *esPump) SetFilter(filter *pumps.Filter) {
	e.filter = filter
}

func (e *esPump) GetFilter() *pumps.Filter {
	return e.filter
}

func (e *esPump) SetTimeout(duration time.Duration) {
	e.timeout = duration
}

func (e *esPump) GetTimeout() time.Duration {
	return e.timeout
}

func (e *esPump) SetOmitDetail(b bool) {
	e.omitDetail = b
}

func (e *esPump) GetOmitDetail() bool {
	return e.omitDetail
}

func (e *esPump) Write(i []interface{}) error {
	for _, data := range i {
		if e.ctx.Err() != nil {
			continue
		}

		info, ok := data.(*analytics.RecordInfo)
		if !ok {
			return fmt.Errorf("data %v is fail to convert to RecordInfo", data)

		}

		if e.config.UseBulk {
			req := elastic.NewBulkCreateRequest().Index(e.indexName()).Type(e.config.DocumentType).Doc(e.mapping(info))
			e.client.BulkProcessor().Add(req)
		} else {
			if _, err := e.client.Client().Index().Index(e.indexName()).BodyJson(e.mapping(info)).Do(e.ctx); err != nil {
				return err
			}
		}
	}
	return nil
}

func (e *esPump) mapping(info *analytics.RecordInfo) map[string]any {
	return map[string]any{
		"@timestamp": info.Timestamp,
		"username":   info.UserName,
		"effect":     info.Effect,
		"request":    info.Request,
		"policies":   info.Policies,
		"deciders":   info.Deciders,
		"conclusion": info.Conclusion,
		"expire-at":  info.ExpireAt}
}

func (e *esPump) indexName() string {
	var index string
	if e.config.RollingIndex {
		index += "-" + time.Now().Format("2006-01-02")
	}
	return index
}

var _ pumps.Pump = &esPump{}
