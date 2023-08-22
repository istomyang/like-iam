package authorization

import (
	"encoding/json"
	"fmt"
	"github.com/ory/ladon"
	"istomyang.github.com/like-iam/iam/internal/authzserver/analytics"
	common "istomyang.github.com/like-iam/iam/internal/pkg/analytics"
	"strings"
	"time"
)

type auditor struct {
}

func newAuditor() (ladon.AuditLogger, error) {
	return &auditor{}, nil
}

var _ ladon.AuditLogger = &auditor{}

func (a *auditor) LogRejectedAccessRequest(request *ladon.Request, pool ladon.Policies, deciders ladon.Policies) {
	var conclusion string

	// you can refer to ladon 's caller of LogRejectedAccessRequest to see how to work.
	if len(deciders) == 1 {
		conclusion = fmt.Sprintf("request rejectd by policy %s.", deciders[0].GetID())
	} else if len(deciders) > 1 {
		conclusion = fmt.Sprintf("request allowed by %s and rejectd by policy %s.",
			mergeId2String(deciders[:len(deciders)-1]),
			deciders[len(deciders)-1].GetID())
	} else {
		conclusion = fmt.Sprintf("request rejectd by all policies.")
	}

	r, _ := json.Marshal(request)
	p, _ := json.Marshal(pool)
	d, _ := json.Marshal(deciders)

	info := common.RecordInfo{
		Timestamp:  time.Now().Unix(),
		UserName:   request.Context["username"].(string),
		Effect:     ladon.DenyAccess,
		Conclusion: conclusion,
		Request:    string(r),
		Policies:   string(p),
		Deciders:   string(d),
	}

	_ = analytics.GetAnalytics().Record(&info)
}

func (a *auditor) LogGrantedAccessRequest(request *ladon.Request, pool ladon.Policies, deciders ladon.Policies) {
	var conclusion = fmt.Sprintf("request allowed by policies: %s", mergeId2String(deciders))
	r, _ := json.Marshal(request)
	p, _ := json.Marshal(pool)
	d, _ := json.Marshal(deciders)

	info := common.RecordInfo{
		Timestamp:  time.Now().Unix(),
		UserName:   request.Context["username"].(string),
		Effect:     ladon.AllowAccess,
		Conclusion: conclusion,
		Request:    string(r),
		Policies:   string(p),
		Deciders:   string(d),
	}

	_ = analytics.GetAnalytics().Record(&info)
}

func mergeId2String(policies ladon.Policies) string {
	var builder strings.Builder
	for i, policy := range policies {
		builder.WriteString(policy.GetID())
		if i != len(policies)-1 {
			builder.WriteString(",")
		}
	}
	return builder.String()
}
