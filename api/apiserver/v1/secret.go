package v1

import (
	"gorm.io/gorm"
	metav1 "istomyang.github.com/like-iam/component-base/meta/v1"
	"istomyang.github.com/like-iam/component-base/util/idutil"
)

type Secret struct {
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Username string `json:"username"           gorm:"column:username"  validate:"omitempty"`

	SecretID  string `json:"secretID"           gorm:"column:secret-id"  validate:"omitempty"`
	SecretKey string `json:"secretKey"          gorm:"column:secret-key" validate:"omitempty"`

	// Required: true
	Expires     int64  `json:"expires"     gorm:"column:expires"     validate:"omitempty"`
	Description string `json:"description" gorm:"column:description" validate:"description"`
}

func (u *Secret) TableName() string {
	return "secret"
}

func (u *Secret) AfterCreate(tx *gorm.DB) error {
	var err error
	if u.InstanceID, err = idutil.GetInstanceId(u.ID, "secret", 6); err != nil {
		return err
	}

	return tx.Save(u).Error
}

type SecretList struct {
	metav1.ListMeta `json:",inline"`

	Items []*Secret `json:"items"`
}
