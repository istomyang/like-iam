package fake

import (
	"context"
	v1 "istomyang.github.com/like-iam/api/apiserver/v1"
	"istomyang.github.com/like-iam/component-base/errors"
	metav1 "istomyang.github.com/like-iam/component-base/meta/v1"
	"istomyang.github.com/like-iam/iam/internal/apiserver/store"
	"istomyang.github.com/like-iam/iam/internal/pkg/codes"
	"strings"
)

type user struct {
	db *datastore
}

func newUser(ds *datastore) store.UserStore {
	return &user{db: ds}
}

func (u *user) Create(c context.Context, user *v1.User, opts metav1.CreateOperateMeta) error {
	u.db.Lock()
	defer u.db.Unlock()

	for _, u := range u.db.users {
		if u.Username == user.Username {
			return errors.WithCode(codes.ErrUserAlreadyExist, "username \"%s\" has already existed.", user.Username)
		}
	}

	user.ID = uint64(len(u.db.users) + 1)
	u.db.users = append(u.db.users, user)

	return nil
}

func (u *user) Update(c context.Context, user *v1.User, opts metav1.UpdateOperateMeta) error {
	u.db.Lock()
	defer u.db.Unlock()

	for i, v := range u.db.users {
		if v.Username == user.Username {
			u.db.users[i] = user
			return nil
		}
	}

	return errors.WithCode(codes.ErrUserNotFound, "username \"%s\" update fail, not found.", user.Username)
}

func (u *user) Delete(c context.Context, username string, opts metav1.DeleteOperateMeta) error {
	u.db.Lock()
	defer u.db.Unlock()

	ps := u.db.policies
	for i, v := range ps {
		if v.Username == username {
			ps[i] = ps[len(ps)-1]
			u.db.policies = ps[:len(ps)-1]
		}
	}

	ss := u.db.secrets
	for i, v := range ss {
		if v.Username == username {
			ss[i] = ss[len(ss)-1]
			u.db.secrets = ss[:len(ss)-1]
		}
	}

	us := u.db.users
	for i, v := range us {
		if v.Username == username {
			us[i] = us[len(us)-1]
			u.db.users = us[:len(us)-1]
			return nil
		}
	}

	return errors.WithCode(codes.ErrUserNotFound, "username \"%s\" update fail, not found.", username)
}

func (u *user) DeleteCollection(c context.Context, usernames []string, opts metav1.DeleteOperateMeta) error {
	u.db.Lock()
	defer u.db.Unlock()

	us := u.db.users
	uss := usernames
	for i, v := range u.db.users {
		for j, username := range usernames {
			if v.Username == username {
				uss[j] = uss[len(uss)-1]
				uss = uss[:len(uss)-1]

				us[i] = us[len(us)-1]
				us = us[:len(us)-1]
			}
		}
	}

	u.db.users = us

	return errors.WithCode(codes.ErrUserNotFound, "username \"%v\" update fail, not found.", uss)
}

func (u *user) Get(c context.Context, username string, opts metav1.GetOperateMeta) (*v1.User, error) {
	u.db.Lock()
	defer u.db.Unlock()

	for _, v := range u.db.users {
		if v.Username == username {
			return v, nil
		}
	}

	return nil, errors.WithCode(codes.ErrUserNotFound, "username \"%s\" update fail, not found.", username)
}

func (u *user) List(c context.Context, opts metav1.ListOperateMeta) (*v1.UserList, error) {
	u.db.Lock()
	defer u.db.Unlock()

	var us []*v1.User

	for _, v := range u.db.users {
		if strings.Contains(v.Username, opts.FieldSelector) {
			us = append(us, v)
		}
	}

	return &v1.UserList{
		ListMeta: metav1.ListMeta{TotalCount: int64(len(us))},
		Items:    us,
	}, nil
}
