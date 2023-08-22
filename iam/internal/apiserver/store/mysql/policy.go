package mysql

import (
	"context"
	"gorm.io/gorm"
	v1 "istomyang.github.com/like-iam/api/apiserver/v1"
	"istomyang.github.com/like-iam/component-base/errors"
	metav1 "istomyang.github.com/like-iam/component-base/meta/v1"
	"istomyang.github.com/like-iam/iam/internal/apiserver/store"
	"istomyang.github.com/like-iam/iam/internal/pkg/codes"
)

type policy struct {
	db *gorm.DB
}

func newPolicy(ds *datastore) store.PolicyStore {
	return &policy{db: ds.db}
}

func (p *policy) Create(c context.Context, policy *v1.Policy, opts metav1.CreateOperateMeta) error {
	return p.db.WithContext(c).Create(&policy).Error
}

func (p *policy) Update(c context.Context, policy *v1.Policy, opts metav1.UpdateOperateMeta) error {
	return p.db.WithContext(c).Save(&policy).Error
}

func (p *policy) Delete(c context.Context, username string, name string, opts metav1.DeleteOperateMeta) error {
	if opts.Unscoped {
		p.db = p.db.Unscoped()
	}
	err := p.db.WithContext(c).Where("username = ?", username).Delete(&v1.Secret{}).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.WithCode(errors.ErrDatabase, err.Error())
	}
	return nil
}

func (p *policy) deleteByUser(c context.Context, username string, opts metav1.DeleteOperateMeta) error {
	if opts.Unscoped {
		p.db = p.db.Unscoped()
	}
	return p.db.WithContext(c).Where("username = ?", username).Delete(&v1.Secret{}).Error
}

func (p *policy) DeleteCollection(c context.Context, username string, names []string, opts metav1.DeleteOperateMeta) error {
	var err error
	if opts.Unscoped {
		p.db = p.db.Unscoped()
	}
	err = p.db.WithContext(c).Where("username = ? and instanceID in (?)", username, names).Delete(&v1.Policy{}).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.WithCode(errors.ErrDatabase, err.Error())
	}
	return err
}

func (p *policy) deleteCollectionByUser(c context.Context, usernames []string, opts metav1.DeleteOperateMeta) error {
	if opts.Unscoped {
		p.db = p.db.Unscoped()
	}
	return p.db.WithContext(c).Where("username in (?)", usernames).Delete(&v1.Policy{}).Error
}

func (p *policy) Get(c context.Context, username string, name string, opts metav1.GetOperateMeta) (*v1.Policy, error) {
	r := &v1.Policy{}
	err := p.db.Where("username = ? and name = ?", username, name).First(&r).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.WithCode(codes.ErrPolicyNotFound, err.Error())
		}

		return nil, errors.WithCode(errors.ErrDatabase, err.Error())
	}

	return r, nil
}

func (p *policy) List(c context.Context, username string, opts metav1.ListOperateMeta) (*v1.PolicyList, error) {
	var r v1.PolicyList
	d := p.db.Where("username LIKE ?", "%"+opts.FieldSelector+"%").
		Limit(int(*opts.Limit)).
		Offset(int(*opts.Offset)).
		Order("id desc").
		Find(&r).
		Offset(-1).
		Limit(-1).
		Count(&r.TotalCount)
	return &r, d.Error
}

func (p *policy) ClearOutdated(c context.Context, maxReserveDays int) (int64, error) {
	//TODO implement me
	panic("implement me")
}
