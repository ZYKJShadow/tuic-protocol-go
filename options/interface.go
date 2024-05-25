package options

type IOption interface {
	Marshal() ([]byte, error)
	Unmarshal(b []byte) error
	Len() uint32
}
