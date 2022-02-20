package core

import "fmt"

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

type ErrConflict struct {
	conflicting interface{}
}

func NewErrConflict(conflicting interface{}) ErrConflict {
	return ErrConflict{conflicting: conflicting}
}

func (e ErrConflict) Error() string {
	return fmt.Sprintf("Conflict detected. Conflicting objects are: %#v", e.conflicting)
}
