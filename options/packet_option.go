package options

import (
	"errors"
	"tuic-protocol-go/address"
	"tuic-protocol-go/utils"
)

type PacketOptions struct {
	AssocID   uint16
	PacketID  uint16
	FragTotal uint8
	FragID    uint8
	Size      uint16
	Addr      address.Address
}

var _ IOption = (*PacketOptions)(nil)

func (p *PacketOptions) Marshal() ([]byte, error) {
	b := make([]byte, 8+p.Addr.Len())
	utils.WriteUint16(b[0:2], p.AssocID)
	utils.WriteUint16(b[2:4], p.PacketID)
	utils.WriteUint8(b[4:5], p.FragTotal)
	utils.WriteUint8(b[5:6], p.FragID)
	utils.WriteUint16(b[6:8], p.Size)

	addr, err := p.Addr.Marshal()
	if err != nil {
		return nil, err
	}

	copy(b[8:], addr)

	return b, nil
}

func (p *PacketOptions) Unmarshal(b []byte) error {
	if len(b) < 8 {
		return errors.New("invalid packet options length")
	}

	p.AssocID = utils.ReadUint16(b[0:2])
	p.PacketID = utils.ReadUint16(b[2:4])
	p.FragTotal = utils.ReadUint8(b[4:5])
	p.FragID = utils.ReadUint8(b[5:6])
	p.Size = utils.ReadUint16(b[6:8])

	return nil
}

func (p *PacketOptions) Len() uint32 {
	return 2 + 2 + 1 + 1 + 2 + uint32(p.Addr.Len())
}

func (p *PacketOptions) CalFragTotal(payload []byte, maxPktSize uint32) {
	firstFragSize := maxPktSize - p.Len()
	fragSizeAddrNone := maxPktSize - 8
	i := uint32(len(payload))
	if firstFragSize < i {
		p.FragTotal = uint8(1 + (i-firstFragSize)/fragSizeAddrNone + 1)
	}
}
