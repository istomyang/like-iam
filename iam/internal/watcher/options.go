package watcher

import generaloptions "istomyang.github.com/like-iam/component/pkg/options"

type Options struct {
	RedisOptions *generaloptions.RedisOpts `json:"redis-options,omitempty" mapstructure:"redis-options"`
	MysqlOptions *generaloptions.MySQLOpts `json:"mysql-options,omitempty" mapstructure:"mysql-options"`

	Watcher *WatchOpts `json:"watcher,omitempty" mapstructure:"watcher"`
}

type WatchOpts struct {
	Clean *CleanOpts `json:"clean" mapstructure:"clean"`
}

type CleanOpts struct {
	MaxUserActiveDays int `json:"max-user-active-days" mapstructure:"max-user-active-days"`
}
