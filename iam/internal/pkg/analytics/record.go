package analytics

import (
	"istomyang.github.com/like-iam/component-base/json"
	"time"
)

type RecordInfoBytes []byte

type RecordInfo struct {
	Timestamp int64     `json:"timestamp" msgpack:"timestamp"`
	ExpireAt  time.Time `json:"expireAt" msgpack:"expireAt"`

	UserName   string `json:"userName" msgpack:"userName"`
	Effect     string `json:"effect" msgpack:"effect"`         // Allow or Deny
	Conclusion string `json:"conclusion" msgpack:"conclusion"` // Who allow and who deny, use their id
	Request    string `json:"request" msgpack:"request"`       // current ladon.Request
	Policies   string `json:"policies" msgpack:"policies"`     // corresponding policies
	Deciders   string `json:"deciders" msgpack:"deciders"`     // who visits this request
}

func (info *RecordInfo) Marshal() (RecordInfoBytes, error) {
	return json.MPMarshal(info)
}
