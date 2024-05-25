package address

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"tuic-go/tuic/utils"
)

const (
	AddrTypeNone   = 0xff
	AddrTypeDomain = 0x00
	AddrTypeIPv4   = 0x01
	AddrTypeIPv6   = 0x02
)

type Address interface {
	TypeCode() uint8
	Len() int
	String() string
	Marshal() ([]byte, error)
	ResolveDNS() ([]net.TCPAddr, error)
}

var _ Address = (*NoneAddress)(nil)

type NoneAddress struct{}

func (a *NoneAddress) Marshal() ([]byte, error) {
	return nil, nil
}
func (a *NoneAddress) TypeCode() uint8 { return 0xff }
func (a *NoneAddress) Len() int        { return 1 }
func (a *NoneAddress) String() string  { return "none" }
func (a *NoneAddress) ResolveDNS() ([]net.TCPAddr, error) {
	return nil, errors.New("none address")
}

var _ Address = (*DomainAddress)(nil)

type DomainAddress struct {
	Domain string
	Port   uint16
}

func (a *DomainAddress) Marshal() ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteByte(a.TypeCode())
	buf.WriteByte(byte(len(a.Domain)))
	buf.WriteString(a.Domain)
	err := binary.Write(&buf, binary.BigEndian, a.Port)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (a *DomainAddress) TypeCode() uint8 { return 0x00 }
func (a *DomainAddress) Len() int        { return 1 + 1 + len(a.Domain) + 2 }
func (a *DomainAddress) String() string  { return fmt.Sprintf("%s:%d", a.Domain, a.Port) }
func (a *DomainAddress) ResolveDNS() ([]net.TCPAddr, error) {
	ips, err := net.DefaultResolver.LookupIPAddr(context.Background(), a.Domain)
	if err != nil {
		return nil, err
	}
	var result []net.TCPAddr
	for _, ip := range ips {
		result = append(result, net.TCPAddr{
			IP:   ip.IP,
			Port: int(a.Port),
		})
	}
	return result, nil
}

var _ Address = (*SocketAddress)(nil)

type SocketAddress struct {
	Addr net.TCPAddr
}

func (a *SocketAddress) Marshal() ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteByte(a.TypeCode())

	if a.Addr.IP.To4() != nil {
		buf.Write(a.Addr.IP.To4())
	} else {
		buf.Write(a.Addr.IP)
	}

	err := binary.Write(&buf, binary.BigEndian, uint16(a.Addr.Port))
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (a *SocketAddress) TypeCode() uint8 {
	if a.Addr.IP.To4() != nil {
		return 0x01
	}
	return 0x02
}

func (a *SocketAddress) Len() int {
	if a.Addr.IP.To4() != nil {
		return 1 + 4 + 2
	}
	return 1 + 16 + 2
}

func (a *SocketAddress) String() string { return a.Addr.String() }

func (a *SocketAddress) ResolveDNS() ([]net.TCPAddr, error) {
	return []net.TCPAddr{a.Addr}, nil
}

func UnMarshalAddr(r io.Reader) (Address, error) {
	var buf [1]byte
	_, err := io.ReadFull(r, buf[:])
	if err != nil {
		return nil, err
	}

	typeCode := buf[0]

	switch typeCode {
	case AddrTypeNone:
		return &NoneAddress{}, nil
	case AddrTypeDomain:
		_, err = io.ReadFull(r, buf[:])
		if err != nil {
			return nil, err
		}

		lenBytes := int(buf[0])

		data := make([]byte, lenBytes+2)
		_, err = io.ReadFull(r, data)
		if err != nil {
			return nil, err
		}

		port := utils.ReadUint16(data[lenBytes:])
		domain := string(data[:lenBytes])

		return &DomainAddress{Domain: domain, Port: port}, nil
	case AddrTypeIPv4:
		var b [6]byte
		_, err = io.ReadFull(r, b[:])
		if err != nil {
			return nil, err
		}
		ip := net.IPv4(b[0], b[1], b[2], b[3])
		port := binary.BigEndian.Uint16(b[4:])
		return &SocketAddress{Addr: net.TCPAddr{IP: ip, Port: int(port)}}, nil
	case AddrTypeIPv6:
		var b [18]byte
		_, err := io.ReadFull(r, b[:])
		if err != nil {
			return nil, err
		}
		ip := make(net.IP, net.IPv6len)
		for i := 0; i < net.IPv6len; i += 2 {
			ip[i] = b[i]
			ip[i+1] = b[i+1]
		}
		port := binary.BigEndian.Uint16(b[16:])
		return &SocketAddress{Addr: net.TCPAddr{IP: ip, Port: int(port)}}, nil
	default:
		return nil, errors.New("unsupported address type")
	}
}
