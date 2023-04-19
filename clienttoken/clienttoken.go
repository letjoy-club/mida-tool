package clienttoken

import "strings"

type ClientToken string

func (c ClientToken) IsInvalid() bool {
	return strings.HasPrefix(string(c), "$")
}

func (c ClientToken) IsUser() bool {
	return strings.HasPrefix(string(c), "u_")
}

func (c ClientToken) IsAdmin() bool {
	return strings.HasPrefix(string(c), "1")
}

func (c ClientToken) IsInternal() bool {
	return strings.HasPrefix(string(c), "2")
}

func (c ClientToken) IsAnonymous() bool {
	return c == ""
}

func (c ClientToken) UserID() string {
	if c.IsUser() {
		return c.String()
	}
	return ""
}

func (c ClientToken) Ptr() *string {
	s := string(c)
	return &s
}

func (c ClientToken) String() string {
	return string(c)
}
