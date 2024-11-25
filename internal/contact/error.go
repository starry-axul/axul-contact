package contact

import "fmt"

type ErrNotFound struct {
	ContactID string
}

func (e ErrNotFound) Error() string {
	return fmt.Sprintf("contact '%s' doesn't exist", e.ContactID)
}
