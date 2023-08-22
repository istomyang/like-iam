package store

import (
	"context"
	v1 "istomyang.github.com/like-iam/api/apiserver/v1"
	metav1 "istomyang.github.com/like-iam/component-base/meta/v1"
)

type UserStore interface {
	Create(c context.Context, user *v1.User, opts metav1.CreateOperateMeta) error
	Update(c context.Context, user *v1.User, opts metav1.UpdateOperateMeta) error
	Delete(c context.Context, username string, opts metav1.DeleteOperateMeta) error
	DeleteCollection(c context.Context, usernames []string, opts metav1.DeleteOperateMeta) error
	Get(c context.Context, username string, opts metav1.GetOperateMeta) (*v1.User, error)
	List(c context.Context, opts metav1.ListOperateMeta) (*v1.UserList, error)
}
