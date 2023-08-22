package v1

import (
	"encoding/json"
	"github.com/ory/ladon"
	"gorm.io/gorm"
	metav1 "istomyang.github.com/like-iam/component-base/meta/v1"
	"istomyang.github.com/like-iam/component-base/util/idutil"
)

// AuthzPolicy defines iam policy type.
type AuthzPolicy struct {
	ladon.DefaultPolicy
}

// String returns the string format of Policy.
func (ap AuthzPolicy) String() string {
	data, _ := json.Marshal(ap)
	return string(data)
}

func (ap AuthzPolicy) Load(policyShadow string) error {
	if err := json.Unmarshal([]byte(policyShadow), &ap); err != nil {
		return err
	}
	return nil
}

type Policy struct {
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// The user of the policy.
	Username string `json:"username" gorm:"column:username" validate:"omitempty"`

	// AuthzPolicy policy, will not be stored in db.
	Policy AuthzPolicy `json:"policy,omitempty" gorm:"-" validate:"omitempty"`
	// Policy authorize.DefaultPolicy `json:"policy,omitempty" gorm:"-" validate:"omitempty"`

	// The authorize policy content, just a string format of ladon.DefaultPolicy. DO NOT modify directly.
	PolicyShadow string `json:"-" gorm:"column:policyShadow" validate:"omitempty"`
}

func (p *Policy) TableName() string {
	return "policy"
}

func (p *Policy) BeforeCreate(tx *gorm.DB) error {
	if err := p.ObjectMeta.BeforeCreate(tx); err != nil {
		return err
	}

	p.Policy.ID = p.Name
	p.PolicyShadow = p.Policy.String()

	return nil
}

func (p *Policy) AfterCreate(tx *gorm.DB) error {
	var err error
	if p.InstanceID, err = idutil.GetInstanceId(p.ID, "policy", 6); err != nil {
		return err
	}

	return tx.Save(p).Error
}

func (p *Policy) BeforeUpdate(tx *gorm.DB) error {
	if err := p.ObjectMeta.BeforeUpdate(tx); err != nil {
		return err
	}

	p.Policy.ID = p.Name
	p.PolicyShadow = p.Policy.String()

	return nil
}

func (p *Policy) AfterFind(tx *gorm.DB) error {
	if err := p.ObjectMeta.AfterFind(tx); err != nil {
		return err
	}

	p.Policy = AuthzPolicy{}
	return p.Policy.Load(p.PolicyShadow)
}

type PolicyList struct {
	metav1.ListMeta `json:",inline"`

	Items []*Policy `json:"items"`
}
