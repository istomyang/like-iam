package v1

import (
	"gorm.io/gorm"
	"istomyang.github.com/like-iam/component-base/auth"
	metav1 "istomyang.github.com/like-iam/component-base/meta/v1"
	"istomyang.github.com/like-iam/component-base/util/idutil"
	"time"
)

type User struct {
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Username string `json:"username,omitempty" gorm:"username" validate:"required,min=1,max=30"`

	Password string `json:"password,omitempty" gorm:"password"`

	IsAdmin string `json:"isAdmin,omitempty" gorm:"isAdmin"`

	LoginAt time.Time `json:"loginAt" gorm:"loginAt"`

	TotalPolicy int64 `json:"totalPolicy" gorm:"-" validate:"omitempty"`
}

func (u *User) TableName() string {
	return "user"
}

func (u *User) Compare(password string) bool {
	return auth.Compare(u.Password, password)
}

func (u *User) AfterCreate(tx *gorm.DB) error {
	var err error
	if u.InstanceID, err = idutil.GetInstanceId(u.ID, "user", 6); err != nil {
		return err
	}

	return tx.Save(u).Error
}

type UserList struct {
	metav1.ListMeta `json:",inline"`
	Items           []*User `json:"items"`
}
