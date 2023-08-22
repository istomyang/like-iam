package codes

// iam-apiserver: user codes.
const (
	// ErrUserNotFound - 400: Secret reach the max count.
	ErrUserNotFound int = iota + 110001

	// ErrUserAlreadyExist - 404: Secret not found.
	ErrUserAlreadyExist
)

// iam-apiserver: secret codes.
const (
	// ErrReachMaxCount - 400: Secret reach the max count.
	ErrReachMaxCount int = iota + 110101

	// ErrSecretNotFound - 404: Secret not found.
	ErrSecretNotFound

	// ErrSecretAlreadyExit - 304: Secret not found.
	ErrSecretAlreadyExit
)

// iam-apiserver: policy errors.
const (
	// ErrPolicyNotFound - 404: Policy not found.
	ErrPolicyNotFound int = iota + 110201

	// ErrPolicyAlreadyExit - 304: Secret not found.
	ErrPolicyAlreadyExit
)
