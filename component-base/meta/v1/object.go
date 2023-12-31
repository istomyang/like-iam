package v1

import (
	"gorm.io/gorm"
	"time"
)

// ObjectMeta defines metadata for every resource, which used in persistent media.
// Generally, different resources has different metadata, so in the future, we use common metadata + own-metadata.
type ObjectMeta struct {
	// ID is the unique in time and space value for this object. It is typically generated by
	// the storage on successful creation of a resource and is not allowed to change on PUT
	// operations.
	//
	// Populated by the system.
	// Read-only.
	ID uint64 `json:"id,omitempty" gorm:"primary_key;AUTO_INCREMENT;column:id"`

	// InstanceID defines a string type resource identifier,
	// use prefixed to distinguish resource types, easy to remember, Url-friendly.
	// like: user-x89snb
	InstanceID string `json:"instanceID,omitempty" gorm:"unique;column:instanceID;type:varchar(32);not null"`

	// Required: true
	// Name must be unique. Is required when creating resources.
	// Name is primarily intended for creation idempotence and configuration
	// definition.
	// It will be generated automated only if Name is not specified.
	// Cannot be updated.
	Name string `json:"name,omitempty" gorm:"column:name;type:varchar(64);not null" validate:"name"`

	// Extend store the fields that need to be added, but do not want to add a new table column, will not be stored in db.
	Extend Extend `json:"extend,omitempty" gorm:"-" validate:"omitempty"`

	// ExtendShadow is the shadow of Extend. DO NOT modify directly.
	ExtendShadow string `json:"-" gorm:"column:extendShadow" validate:"omitempty"`

	// CreatedAt is a timestamp representing the server time when this object was
	// created. It is not guaranteed to be set in happens-before order across separate operations.
	// Clients may not set this value. It is represented in RFC3339 form and is in UTC.
	//
	// Populated by the system.
	// Read-only.
	// Null for lists.
	CreatedAt time.Time `json:"createdAt,omitempty" gorm:"column:createdAt"`

	// UpdatedAt is a timestamp representing the server time when this object was updated.
	// Clients may not set this value. It is represented in RFC3339 form and is in UTC.
	//
	// Populated by the system.
	// Read-only.
	// Null for lists.
	UpdatedAt time.Time `json:"updatedAt,omitempty" gorm:"column:updatedAt"`

	// DeletedAt is RFC 3339 date and time at which this resource will be deleted. This
	// field is set by the server when a graceful deletion is requested by the user, and is not
	// directly settable by a client.
	//
	// Populated by the system when a graceful deletion is requested.
	// Read-only.
	DeletedAt gorm.DeletedAt `json:"-" gorm:"column:deletedAt;index:idx_deletedAt"`
}

func (m *ObjectMeta) BeforeCreate(tx *gorm.DB) (err error) {
	m.ExtendShadow = m.Extend.String()
	return nil
}

func (m *ObjectMeta) BeforeUpdate(tx *gorm.DB) (err error) {
	m.ExtendShadow = m.Extend.String()
	return nil
}

func (m *ObjectMeta) AfterFind(tx *gorm.DB) (err error) {
	m.Extend = Extend{}
	return m.Extend.Load(m.ExtendShadow)
}

type ListMeta struct {
	TotalCount int64 `json:"totalCount,omitempty"`
}
