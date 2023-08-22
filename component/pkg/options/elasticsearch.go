package options

import (
	"github.com/spf13/pflag"
	"strings"
	"time"
)

type ElasticSearchOpts struct {
	IndexName    string `json:"index-name" mapstructure:"index-name"`
	RollingIndex bool   `json:"rolling-index" mapstructure:"rolling-index"` // Index name by Date.
	DocumentType string `json:"document-type" mapstructure:"document-type"`
	Urls         string `json:"urls" mapstructure:"urls"`

	Username     string `json:"username" mapstructure:"username"`
	Password     string `json:"password" mapstructure:"password"`
	SniffEnabled bool   `json:"sniff-enabled" mapstructure:"sniff-enabled"`
	GZip         bool   `json:"gzip" mapstructure:"gzip"`

	// APIKey and APIKeyID are merged with "%s:%s" format and put into HTTP's Authorization Header.
	// Like: Authorization APIKey APIKeyID:APIKey, which encoded by base64.
	APIKey   string `json:"auth_api_key" mapstructure:"auth_api_key"`
	APIKeyID string `json:"auth_api_id" mapstructure:"auth_api_id"`

	UseBulk bool     `json:"use-bulk" mapstructure:"use-bulk"`
	Bulk    BulkOpts `json:"bulk" mapstructure:"bulk"`
}

type BulkOpts struct {
	Name          string        `json:"name" mapstructure:"name"`
	Workers       int           `json:"workers" mapstructure:"workers"`
	BulkSize      int           `json:"bulk_size" mapstructure:"bulk-size"`
	BulkActions   int           `json:"bulk_actions" mapstructure:"bulk_actions"`
	FlushInterval time.Duration `json:"flush_interval" mapstructure:"flush_interval"`
}

func (o *ElasticSearchOpts) Validate() []error {
	var err []error

	// TODO: check number is correct.

	return err
}

func (o *ElasticSearchOpts) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.IndexName, "index-name", o.IndexName, "Index name.")
	fs.BoolVar(&o.RollingIndex, "enable-rolling-index", o.RollingIndex, "enable rolling index with name like user-2006.3.12")
	fs.StringVar(&o.DocumentType, "urls", o.DocumentType, "Document type.")

	fs.StringVar(&o.Urls, "urls", o.Urls, "Urls split with `,`, and no space.")
	fs.StringVar(&o.Username, "username", o.Username, "Username.")
	fs.StringVar(&o.Password, "password", o.Urls, "Password.")
	fs.BoolVar(&o.SniffEnabled, "sniff-enabled", o.SniffEnabled, "Enable Sniff.")
	fs.BoolVar(&o.GZip, "gzip", o.SniffEnabled, "Enable Gzip.")
	fs.StringVar(&o.APIKey, "auth-api-key", o.APIKey, "If set api-key, ignore username and password.")
	fs.StringVar(&o.APIKeyID, "auth-api-key-id", o.APIKey, "If set api-key-id, ignore username and password.")
	fs.BoolVar(&o.UseBulk, "use-bulk", o.UseBulk, "Use es bulk, set false will ignore bulk config.")

	fs.StringVar(&o.Bulk.Name, "bulk.name", o.APIKey, "Bulk name.")
	fs.IntVar(&o.Bulk.Workers, "bulk.worker-number", o.Bulk.Workers, "Bulk's worker number.")
	fs.IntVar(&o.Bulk.BulkSize, "bulk.bulk-size", o.Bulk.BulkSize, "Bulk size.")
	fs.IntVar(&o.Bulk.BulkActions, "bulk.bulk-actions", o.Bulk.BulkSize, "Bulk actions number.")
	fs.DurationVar(&o.Bulk.FlushInterval, "bulk.flush-interval", o.Bulk.FlushInterval, "Bulk's flush interval.")
}

func (o *ElasticSearchOpts) UrlsToSlice() []string {
	return strings.Split(o.Urls, ",")
}
