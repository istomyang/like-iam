package pumps

import (
	"context"
	"istomyang.github.com/like-iam/iam/internal/pkg/analytics"
	"time"
)

// Pump abstracts a storage object.
type Pump interface {
	// Init do initialized work with config Pump 's back service use.
	// Config is PumpConfig.Extend.
	Init(ctx context.Context, config map[string]any) error
	Run() error
	Close() error
	// GetName used to identify this impl in logger.
	GetName() string
	// Write can include data with structure.
	Write([]any) error
	SetFilter(*Filter)
	GetFilter() *Filter
	SetTimeout(time.Duration)
	GetTimeout() time.Duration
	SetOmitDetail(bool)
	GetOmitDetail() bool
}

// PumpConfig is a common config struct which impl decoding it by mapstructure with more config fields.
type PumpConfig struct {
	Name       string         `json:"name" mapstructure:"name"`
	Filter     *Filter        `json:"filter" mapstructure:"filter"`
	OmitDetail bool           `json:"omit-detail" mapstructure:"recordDetail"`
	Timeout    time.Duration  `json:"timeout" mapstructure:"timeout"`
	Extend     map[string]any `json:"extend" mapstructure:"extend"`
}

// Filter determines whose data will skipped.
type Filter struct {
	Usernames     []string `json:"usernames,omitempty"`
	SkipUsernames []string `json:"skipUsernames,omitempty"`
}

// ShouldFilter tells this info whether omitted or not.
// If duplication in SkipUsernames and Usernames, SkipUsernames first.
func (f *Filter) ShouldFilter(info *analytics.RecordInfo) bool {
	if len(f.SkipUsernames) > 0 && f.in(info.UserName, f.SkipUsernames) {
		return true
	}
	if len(f.Usernames) > 0 && !f.in(info.UserName, f.Usernames) {
		return true
	}

	return false
}

func (f *Filter) in(name string, names []string) bool {
	for _, n := range names {
		if name == n {
			return true
		}
	}
	return false
}
