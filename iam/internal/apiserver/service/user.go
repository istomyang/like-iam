package service

import (
	"context"
	v1 "istomyang.github.com/like-iam/api/apiserver/v1"
	metav1 "istomyang.github.com/like-iam/component-base/meta/v1"
	"sync"
)

type UserSvc interface {
	Create(ctx context.Context, user *v1.User, opts metav1.CreateOperateMeta) error
	Update(ctx context.Context, user *v1.User, opts metav1.UpdateOperateMeta) error
	Delete(ctx context.Context, username string, opts metav1.DeleteOperateMeta) error
	DeleteCollection(ctx context.Context, usernames []string, opts metav1.DeleteOperateMeta) error
	Get(ctx context.Context, username string, opts metav1.GetOperateMeta) (*v1.User, error)
	List(ctx context.Context, opts metav1.ListOperateMeta) (*v1.UserList, error)
	ListWithBadPerformance(ctx context.Context, opts metav1.ListOperateMeta) (*v1.UserList, error)
	ChangePassword(ctx context.Context, user *v1.User) error
}

type userSvc struct {
	svc *service
}

func newUserSvc(svc *service) UserSvc {
	return &userSvc{svc: svc}
}

func (u *userSvc) Create(ctx context.Context, user *v1.User, opts metav1.CreateOperateMeta) error {
	return u.svc.store.User().Create(ctx, user, opts)
}

func (u *userSvc) Update(ctx context.Context, user *v1.User, opts metav1.UpdateOperateMeta) error {
	return u.svc.store.User().Update(ctx, user, opts)
}

func (u *userSvc) Delete(ctx context.Context, username string, opts metav1.DeleteOperateMeta) error {
	return u.svc.store.User().Delete(ctx, username, opts)
}

func (u *userSvc) DeleteCollection(ctx context.Context, usernames []string, opts metav1.DeleteOperateMeta) error {
	return u.svc.store.User().DeleteCollection(ctx, usernames, opts)
}

func (u *userSvc) Get(ctx context.Context, username string, opts metav1.GetOperateMeta) (*v1.User, error) {
	return u.svc.store.User().Get(ctx, username, opts)
}

func (u *userSvc) List(ctx context.Context, opts metav1.ListOperateMeta) (*v1.UserList, error) {
	userList, err := u.svc.store.User().List(ctx, opts)
	if err != nil {
		return nil, err
	}

	var sMap sync.Map

	var wg sync.WaitGroup
	errChan := make(chan error, 1)
	finishChan := make(chan bool, 1)

	for _, user := range userList.Items {
		wg.Add(1)
		go func(ur *v1.User) {
			defer wg.Done()

			policy, err := u.svc.store.Policy().List(ctx, ur.Username, metav1.ListOperateMeta{})
			if err != nil {
				errChan <- err
				return
			}
			ur.TotalPolicy = policy.TotalCount
			sMap.Store(ur.InstanceID, ur)
		}(user)
	}

	go func() {
		wg.Wait()
		finishChan <- true
	}()

	select {
	case e := <-errChan:
		return nil, e
	case <-finishChan:
		break
	}

	rs := make([]*v1.User, userList.TotalCount)
	sMap.Range(func(key, value any) bool {
		rs = append(rs, value.(*v1.User))
		return true
	})

	return &v1.UserList{
		ListMeta: metav1.ListMeta{TotalCount: userList.TotalCount},
		Items:    rs,
	}, nil
}

func (u *userSvc) ListWithBadPerformance(ctx context.Context, opts metav1.ListOperateMeta) (*v1.UserList, error) {
	userList, err := u.svc.store.User().List(ctx, opts)
	if err != nil {
		return nil, err
	}

	for _, user := range userList.Items {
		policy, err := u.svc.store.Policy().List(ctx, user.Username, metav1.ListOperateMeta{})
		if err != nil {
			return nil, err
		}
		user.TotalPolicy = policy.TotalCount
	}

	return userList, nil
}

func (u *userSvc) ChangePassword(ctx context.Context, user *v1.User) error {
	return u.svc.store.User().Update(ctx, user, metav1.UpdateOperateMeta{})
}
