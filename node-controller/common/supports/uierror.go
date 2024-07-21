package supports

import (
	"fmt"
)

type UIError struct {
	text string
}

// NewUIError 状态码 400
func NewUIError(text string) error {
	return &UIError{text: text}
}

func NewUIErrorf(format string, args ...interface{}) error {
	return &UIError{
		text: fmt.Sprintf(format, args...),
	}
}

func (u *UIError) Error() string {
	return u.text
}
