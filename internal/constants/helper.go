package constants

func (h Header) String() string {
	return string(h)
}

func (h ContentTypes) String() string {
	return string(h)
}

func (h Compression) String() string {
	return string(h)
}

func (c Context) String() string {
	return string(c)
}

func (s Status) IsValid() bool {
	switch s {
	case StatusActive, StatusInactive, StatusSuspended, StatusBlocked:
		return true
	default:
		return false
	}
}

func (s Status) String() string {
	return string(s)
}

func (r Roles) String() string {
	return string(r)
}
