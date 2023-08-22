package fake

import (
	"fmt"
	"github.com/ory/ladon"
	v1 "istomyang.github.com/like-iam/api/apiserver/v1"
	metav1 "istomyang.github.com/like-iam/component-base/meta/v1"
	"istomyang.github.com/like-iam/iam/internal/apiserver/store"
	"sync"
)

const ResourceCount = 1000

type datastore struct {
	sync.RWMutex
	users    []*v1.User
	secrets  []*v1.Secret
	policies []*v1.Policy
}

func (s *datastore) User() store.UserStore {
	return nil
}

func (s *datastore) Secret() store.SecretStore {
	return nil
}

func (s *datastore) Policy() store.PolicyStore {
	return nil
}

func (s *datastore) Run() error {
	return nil
}

func (s *datastore) Close() error {
	return nil
}

var (
	factory store.Factory
	once    sync.Once
)

func GetFakeFactory() (store.Factory, error) {
	once.Do(func() {
		var r *datastore
		r.users = createUsers(ResourceCount)
		r.secrets = createSecrets(ResourceCount)
		r.policies = createPolicies(ResourceCount)
		factory = r
	})
	return factory, nil
}

func createUsers(count int) []*v1.User {
	var rs []*v1.User

	for i := 0; i < count; i++ {
		rs = append(rs, &v1.User{
			ObjectMeta: metav1.ObjectMeta{
				ID:   uint64(i),
				Name: fmt.Sprintf("user-%d", i),
			},
			Username: fmt.Sprintf("username-%d", i),
			Password: fmt.Sprintf("password-%d", i),
		})
	}

	return rs
}

func createSecrets(count int) []*v1.Secret {
	var rs []*v1.Secret

	for i := 0; i < count; i++ {
		rs = append(rs, &v1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				ID:   uint64(i),
				Name: fmt.Sprintf("secret-%d", i),
			},
			Username:  fmt.Sprintf("username-%d", i),
			SecretID:  fmt.Sprintf("username-%d", i),
			SecretKey: fmt.Sprintf("username-%d", i),
		})
	}

	return rs
}

func createPolicies(count int) []*v1.Policy {
	var rs []*v1.Policy

	for i := 0; i < count; i++ {
		rs = append(rs, &v1.Policy{
			ObjectMeta: metav1.ObjectMeta{
				ID:   uint64(i),
				Name: fmt.Sprintf("policy-%d", i),
			},
			Username: fmt.Sprintf("username-%d", i),
			Policy: v1.AuthzPolicy{
				DefaultPolicy: ladon.DefaultPolicy{},
			},
		})
	}

	return rs
}
