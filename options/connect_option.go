package options

import (
	"bytes"
	"tuic-protocol-go/address"
)

type ConnectOptions struct {
	Addr address.Address
}

var _ IOption = (*ConnectOptions)(nil)

func (c *ConnectOptions) Marshal() ([]byte, error) {
	return c.Addr.Marshal()
}

func (c *ConnectOptions) Unmarshal(b []byte) error {
	addr, err := address.UnMarshalAddr(bytes.NewReader(b))
	if err != nil {
		return err
	}

	c.Addr = addr
	return nil
}

func (c *ConnectOptions) Len() uint32 {
	return uint32(c.Addr.Len())
}
