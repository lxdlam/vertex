package container

type StringMap interface {
	Set([]*StringContainer, []*StringContainer) error
	Get([]*StringContainer) ([]*StringContainer, error)

	Len(*StringContainer) error
	Append(*StringContainer, *StringContainer) error

	Increase(*StringContainer, int64) (int64, error)
	Decrease(*StringContainer, int64) (int64, error)

	GetRange(*StringContainer, int, int) (*StringContainer, error)
	SetRange(*StringContainer, int, *StringContainer) error
}
