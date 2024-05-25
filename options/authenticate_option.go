package options

type AuthenticateOptions struct {
	UUID  []byte
	Token []byte
}

var _ IOption = (*AuthenticateOptions)(nil)

func (a *AuthenticateOptions) Marshal() ([]byte, error) {
	b := make([]byte, 48)
	copy(b[:16], a.UUID[:])
	copy(b[16:], a.Token[:])
	return b, nil
}

func (a *AuthenticateOptions) Unmarshal(b []byte) error {
	a.UUID = make([]byte, 16)
	a.Token = make([]byte, 32)

	copy(a.UUID[:], b[:16])
	copy(a.Token[:], b[16:])
	return nil
}

func (a *AuthenticateOptions) Len() uint32 {
	return 16 + 32
}
