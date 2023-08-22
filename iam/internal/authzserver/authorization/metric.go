package authorization

import "github.com/ory/ladon"

type metric struct {
}

var _ ladon.Metric = &metric{}

func newMetric() (ladon.Metric, error) {
	return &metric{}, nil
}

func (m *metric) RequestDeniedBy(request ladon.Request, policy ladon.Policy) {
	//TODO implement me
}

func (m *metric) RequestAllowedBy(request ladon.Request, policies ladon.Policies) {
	//TODO implement me
}

func (m *metric) RequestNoMatch(request ladon.Request) {
	//TODO implement me
}

func (m *metric) RequestProcessingError(request ladon.Request, policy ladon.Policy, err error) {
	//TODO implement me
}
