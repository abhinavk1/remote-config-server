package versioncontrol

type Abstract interface {
	Checkout(branch string) error

	Pull() error
}
