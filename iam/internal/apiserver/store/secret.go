package store

import (
	"context"
	v1 "istomyang.github.com/like-iam/api/apiserver/v1"
	metav1 "istomyang.github.com/like-iam/component-base/meta/v1"
)

type SecretStore interface {
	Create(c context.Context, secret *v1.Secret, opts metav1.CreateOperateMeta) error
	Update(c context.Context, secret *v1.Secret, opts metav1.UpdateOperateMeta) error
	Delete(c context.Context, username, secretID string, opts metav1.DeleteOperateMeta) error
	DeleteCollection(c context.Context, username string, secretIDs []string, opts metav1.DeleteOperateMeta) error
	Get(c context.Context, username, secretID string, opts metav1.GetOperateMeta) (*v1.Secret, error)
	List(c context.Context, username string, opts metav1.ListOperateMeta) (*v1.SecretList, error)
}
