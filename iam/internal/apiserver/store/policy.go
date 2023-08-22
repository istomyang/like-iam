package store

import (
	"context"
	v1 "istomyang.github.com/like-iam/api/apiserver/v1"
	metav1 "istomyang.github.com/like-iam/component-base/meta/v1"
)

type PolicyStore interface {
	Create(c context.Context, policy *v1.Policy, opts metav1.CreateOperateMeta) error
	Update(c context.Context, policy *v1.Policy, opts metav1.UpdateOperateMeta) error
	Delete(c context.Context, username string, name string, opts metav1.DeleteOperateMeta) error
	DeleteCollection(c context.Context, username string, names []string, opts metav1.DeleteOperateMeta) error
	Get(c context.Context, username string, name string, opts metav1.GetOperateMeta) (*v1.Policy, error)
	List(c context.Context, username string, opts metav1.ListOperateMeta) (*v1.PolicyList, error)

	// ClearOutdated cleans outdated policies.
	// Use DeletedAt field, this means Delete operation just mark item should delete now.
	ClearOutdated(c context.Context, maxReserveDays int) (int64, error)
}
