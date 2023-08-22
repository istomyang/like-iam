package pumps

import (
	"context"
	"fmt"
	"istomyang.github.com/like-iam/iam/internal/pkg/analytics"
	"istomyang.github.com/like-iam/iam/internal/pumper/pumps/csv"
	"istomyang.github.com/like-iam/iam/internal/pumper/pumps/elasticsearch"
	"istomyang.github.com/like-iam/iam/internal/pumper/pumps/empty"
	"istomyang.github.com/like-iam/iam/internal/pumper/pumps/influx"
	"istomyang.github.com/like-iam/iam/internal/pumper/pumps/kafka"
	"istomyang.github.com/like-iam/iam/internal/pumper/pumps/mongo"
	"istomyang.github.com/like-iam/iam/internal/pumper/pumps/prometheus"
	"istomyang.github.com/like-iam/iam/internal/pumper/pumps/sysloger"
)

var pumps map[string]Pump
var presets map[string]Pump

func init() {
	presets = map[string]Pump{
		"empty":         empty.New(),
		"csv":           csv.New(),
		"elasticsearch": elasticsearch.New(),
		"influx":        influx.New(),
		"kafka":         kafka.New(),
		"mongo":         mongo.New(),
		"prometheus":    prometheus.New(),
		"sysloger":      sysloger.New(),
	}
}

// RegisterRun registers pumps from presetting impl.
func RegisterRun(ctx context.Context, p map[string]any) []error {
	var errs []error

	for name, rawConfig := range p {
		var config = rawConfig.(*PumpConfig)
		pmp, exist := presets[name]
		if !exist {
			errs = append(errs, fmt.Errorf("must assgin impl, or implment it, got: %s", name))
			continue
		}
		pumps[name] = pmp
		if err := pmp.Init(ctx, config.Extend); err != nil {
			errs = append(errs, fmt.Errorf("you assign wrong config, please check it: %v", config.Extend))
			continue
		}
		pmp.SetOmitDetail(config.OmitDetail)
		pmp.SetTimeout(config.Timeout)
		pmp.SetFilter(config.Filter)
		errs = append(errs, pmp.Run())
	}

	return errs
}

// Do is a convenient function to run codes for all registered pumps.
func Do(f func(pump Pump) error) []error {
	var errs []error
	for _, pump := range pumps {
		errs = append(errs, f(pump))
	}
	return errs
}

func Close() []error {
	return Do(func(pump Pump) error {
		return pump.Close()
	})
}

// HandleData handles data for every Pump.
// If return nil, skip this data.
func HandleData(config *PumpConfig, info *analytics.RecordInfo) *analytics.RecordInfo {
	if config.OmitDetail {
		info.Deciders = ""
		info.Policies = ""
	}
	if config.Filter.ShouldFilter(info) {
		return nil
	}
	return info
}
