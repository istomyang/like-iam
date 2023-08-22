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

type secret struct {
	db *gorm.DB
}

func newSecret(ds *datastore) store.SecretStore {
	return &secret{db: ds.db}
}

func (s *secret) Create(c context.Context, secret *v1.Secret, opts metav1.CreateOperateMeta) error {
	return s.db.WithContext(c).Create(&secret).Error
}

func (s *secret) Update(c context.Context, secret *v1.Secret, opts metav1.UpdateOperateMeta) error {
	return s.db.WithContext(c).Save(&secret).Error
}

func (s *secret) Delete(c context.Context, username, secretID string, opts metav1.DeleteOperateMeta) error {
	if opts.Unscoped {
		s.db = s.db.Unscoped()
	}
	err := s.db.WithContext(c).Where("username = ? and secret-id = ?", username, secretID).Delete(&v1.Secret{}).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.WithCode(errors.ErrDatabase, err.Error())
	}
	return nil
}

func (s *secret) deleteByUser(c context.Context, username string, opts metav1.DeleteOperateMeta) error {
	if opts.Unscoped {
		s.db = s.db.Unscoped()
	}
	return s.db.WithContext(c).Where("username = ?", username).Delete(&v1.Secret{}).Error
}

func (s *secret) DeleteCollection(c context.Context, username string, secretIDs []string, opts metav1.DeleteOperateMeta) error {
	var err error
	if opts.Unscoped {
		s.db = s.db.Unscoped()
	}
	err = s.db.WithContext(c).Where("username = ? and secret-id in (?)", username, secretIDs).Delete(&v1.Secret{}).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.WithCode(errors.ErrDatabase, err.Error())
	}
	return err
}

func (s *secret) deleteCollectionByUser(c context.Context, usernames []string, opts metav1.DeleteOperateMeta) error {
	if opts.Unscoped {
		s.db = s.db.Unscoped()
	}
	return s.db.WithContext(c).Where("username in (?)", usernames).Delete(&v1.Secret{}).Error
}

func (s *secret) Get(c context.Context, username, secretID string, opts metav1.GetOperateMeta) (*v1.Secret, error) {
	se := &v1.Secret{}
	err := s.db.Where("username = ? and secret-id = ?", username, secretID).First(&se).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.WithCode(codes.ErrSecretNotFound, err.Error())
		}
		return nil, errors.WithCode(errors.ErrDatabase, err.Error())
	}

	return se, nil
}

func (s *secret) List(c context.Context, username string, opts metav1.ListOperateMeta) (*v1.SecretList, error) {
	var r v1.SecretList
	d := s.db.Where("username LIKE ?", "%"+opts.FieldSelector+"%").
		Limit(int(*opts.Limit)).
		Offset(int(*opts.Offset)).
		Order("id desc").
		Find(&r).
		Offset(-1).
		Limit(-1).
		Count(&r.TotalCount)
	return &r, d.Error
}
