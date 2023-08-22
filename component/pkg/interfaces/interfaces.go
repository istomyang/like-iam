package interfaces

type ComponentCommon interface {
	Run() error
	Close() error
}
