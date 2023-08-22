package store

import "github.com/ory/ladon"

// PolicyStore lists data from apiserver server.
type PolicyStore interface {
	List() (map[string][]*ladon.DefaultPolicy, error)
}
