package options

import (
	"tuic-go/tuic/utils"
)

type DissociateOptions struct {
	AssocID uint16
}

var _ IOption = (*DissociateOptions)(nil)

func (d *DissociateOptions) Marshal() ([]byte, error) {
	b := make([]byte, 2)
	utils.WriteUint16(b, d.AssocID)
	return b, nil
}

func (d *DissociateOptions) Unmarshal(b []byte) error {
	d.AssocID = utils.ReadUint16(b)
	return nil
}

func (d *DissociateOptions) Len() uint32 {
	return 2
}
