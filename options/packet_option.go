package options

import (
	"errors"
	"github.com/ZYKJShadow/tuic-protocol-go/address"
	"github.com/ZYKJShadow/tuic-protocol-go/utils"
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
	dataLen := uint32(len(payload))
	if dataLen < maxPktSize {
		p.FragTotal = 1
	} else {
		p.FragTotal = uint8(dataLen / maxPktSize)
		if dataLen%maxPktSize != 0 {
			p.FragTotal++
		}
	}
}
