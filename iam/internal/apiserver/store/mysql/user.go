package mysql

import (
	"context"
	"gorm.io/gorm"
	v1 "istomyang.github.com/like-iam/api/apiserver/v1"
	"istomyang.github.com/like-iam/component-base/errors"
	metav1 "istomyang.github.com/like-iam/component-base/meta/v1"
	"istomyang.github.com/like-iam/iam/internal/apiserver/store"
)

type user struct {
	db *gorm.DB
}

func newUser(ds *datastore) store.UserStore {
	return &user{db: ds.db}
}

func (u *user) Create(c context.Context, user *v1.User, opts metav1.CreateOperateMeta) error {
	return u.db.WithContext(c).Create(&user).Error
}

func (u *user) Update(c context.Context, user *v1.User, opts metav1.UpdateOperateMeta) error {
	return u.db.WithContext(c).Save(&user).Error
}

func (u *user) Delete(c context.Context, username string, opts metav1.DeleteOperateMeta) error {
	var err error
	if err = (&secret{db: u.db}).deleteByUser(c, username, opts); err != nil {
		return err
	}
	if err = (&policy{db: u.db}).deleteByUser(c, username, opts); err != nil {
		return err
	}
	if opts.Unscoped {
		u.db = u.db.Unscoped()
	}
	err = u.db.WithContext(c).Where("username = ?", username).Delete(&v1.User{}).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.WithCode(errors.ErrDatabase, err.Error())
	}
	return err
}

func (u *user) DeleteCollection(c context.Context, usernames []string, opts metav1.DeleteOperateMeta) error {
	var err error
	if err = (&secret{db: u.db}).deleteCollectionByUser(c, usernames, opts); err != nil {
		return err
	}
	if err = (&policy{db: u.db}).deleteCollectionByUser(c, usernames, opts); err != nil {
		return err
	}
	if opts.Unscoped {
		u.db = u.db.Unscoped()
	}
	err = u.db.WithContext(c).Where("username in (?)", usernames).Delete(&v1.User{}).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.WithCode(errors.ErrDatabase, err.Error())
	}
	return err
}

func (u *user) Get(c context.Context, username string, opts metav1.GetOperateMeta) (*v1.User, error) {
	var user *v1.User
	if err := u.db.Where("username = ?", username).First(&user).Error; err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.WithCode(errors.ErrDatabase, err.Error())
	}
	return user, nil
}

func (u *user) List(c context.Context, opts metav1.ListOperateMeta) (*v1.UserList, error) {
	var users v1.UserList
	d := u.db.Where("username LIKE ?", "%"+opts.FieldSelector+"%").
		Limit(int(*opts.Limit)).
		Offset(int(*opts.Offset)).
		Order("id desc").
		Find(&users).
		Offset(-1).
		Limit(-1).
		Count(&users.TotalCount)
	return &users, d.Error
}
