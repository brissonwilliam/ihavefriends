package core

type ErrNotFound struct {
	id string
}

func NewErrNotFound(id string) ErrNotFound {
	return ErrNotFound{
		id: id,
	}
}

func (e ErrNotFound) Error() string {
	return "Could not find " + e.id
}
