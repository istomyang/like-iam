package v1

// OperateMeta defines metadata for all operation.
// For example, fetch data from database in page form must use `limit` and `offset` metadata.
type OperateMeta struct {
	ApiVersion string `json:"api-version,omitempty"`
}

type ListOperateMeta struct {
	OperateMeta `json:",inline"`

	// FieldSelector restricts the list of returned objects by their fields. Defaults to everything.
	// TODO: more research.
	// In IAM project, selector like css's selector, but is not useful.
	FieldSelector string `json:"field-selector,omitempty" form:"fieldSelector"`

	// TimeoutSeconds specifies the seconds of ClientIP type session sticky time.
	TimeoutSeconds *int64 `json:"timeout-seconds,omitempty"`

	// Offset specify the number of records to skip before starting to return the records.
	Offset *int64 `json:"offset,omitempty" form:"offset"`

	// Limit specify the number of records to be retrieved.
	Limit *int64 `json:"limit,omitempty" form:"limit"`
}

type GetOperateMeta struct {
	OperateMeta `json:",inline"`
}

type DeleteOperateMeta struct {
	OperateMeta `json:",inline"`

	// Unscoped replace soft delete operation with hard delete operation.
	// Default gorm db use DeleteAt field to mark this entry need be deleted.
	// +optional
	Unscoped bool `json:"unscoped"`
}

type CreateOperateMeta struct {
	OperateMeta `json:",inline"`

	// DryRun doesn't send operation to database driver. It generates SQL without executing.
	// It can be used to prepare or test generated SQL.
	// https://gorm.io/docs/session.html#DryRun
	// +optional
	DryRun bool `json:"dry-run,omitempty"`
}

type UpdateOperateMeta struct {
	OperateMeta `json:",inline"`

	// Same with CreateOperateMeta.DryRun
	DryRun bool `json:"dry-run,omitempty"`
}
