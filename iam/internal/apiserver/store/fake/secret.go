package fake

import (
	"context"
	v1 "istomyang.github.com/like-iam/api/apiserver/v1"
	"istomyang.github.com/like-iam/component-base/errors"
	metav1 "istomyang.github.com/like-iam/component-base/meta/v1"
	"istomyang.github.com/like-iam/component-base/util/idutil"
	"istomyang.github.com/like-iam/iam/internal/pkg/codes"
	"strings"
)

type secret struct {
	db *datastore
}

func (s *secret) Create(c context.Context, secret *v1.Secret, opts metav1.CreateOperateMeta) error {
	s.db.Lock()
	defer s.db.Unlock()

	for _, v := range s.db.secrets {
		if v.Username == secret.Username && v.SecretID == secret.SecretID && v.Name == secret.Name {
			return errors.WithCode(codes.ErrSecretAlreadyExit, "secret has already exits.")
		}
	}

	secret.ID = uint64(len(s.db.secrets))
	secret.InstanceID, _ = idutil.GetInstanceId(secret.ID, "secret", 6)

	s.db.secrets = append(s.db.secrets, secret)

	return nil
}

func (s *secret) Update(c context.Context, secret *v1.Secret, opts metav1.UpdateOperateMeta) error {
	s.db.Lock()
	defer s.db.Unlock()

	for i, v := range s.db.secrets {
		if v.Username == secret.Username && v.SecretID == secret.SecretID {
			s.db.secrets[i] = secret
			return nil
		}
	}

	return errors.WithCode(codes.ErrSecretNotFound, "secret not found.")
}

func (s *secret) Delete(c context.Context, username, secretID string, opts metav1.DeleteOperateMeta) error {
	s.db.Lock()
	defer s.db.Unlock()

	for i, v := range s.db.secrets {
		if v.Username == username && v.SecretID == secretID {
			s.db.secrets[i] = s.db.secrets[len(s.db.secrets)-1]
			s.db.secrets = s.db.secrets[:len(s.db.secrets)-1]
			return nil
		}
	}

	return errors.WithCode(codes.ErrSecretNotFound, "secret not found.")
}

func (s *secret) DeleteCollection(c context.Context, username string, secretIDs []string, opts metav1.DeleteOperateMeta) error {
	s.db.Lock()
	defer s.db.Unlock()

	ss := s.db.secrets
	for i, v := range ss {
		if v.Username == username {
			for _, d := range secretIDs {
				if v.SecretID == d {
					s.db.secrets[i] = s.db.secrets[len(s.db.secrets)-1]
					s.db.secrets = s.db.secrets[:len(s.db.secrets)-1]
				}
			}
		}
	}

	return nil
}

func (s *secret) Get(c context.Context, username, secretID string, opts metav1.GetOperateMeta) (*v1.Secret, error) {
	s.db.Lock()
	defer s.db.Unlock()

	for _, v := range s.db.secrets {
		if v.Username == username && v.SecretID == secretID {
			return v, nil
		}
	}

	return nil, errors.WithCode(codes.ErrSecretNotFound, "secret-id `%s` in user `%s` not found", secretID, username)
}

func (s *secret) List(c context.Context, username string, opts metav1.ListOperateMeta) (*v1.SecretList, error) {
	s.db.Lock()
	defer s.db.Unlock()

	var r []*v1.Secret
	for _, v := range s.db.secrets {
		if strings.Contains(v.Username, username) {
			r = append(r, v)
		}
	}

	return &v1.SecretList{
		ListMeta: metav1.ListMeta{TotalCount: int64(len(r))},
		Items:    r,
	}, nil
}
