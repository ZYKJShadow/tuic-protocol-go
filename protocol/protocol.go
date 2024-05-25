package protocol

import (
	"errors"
	"github.com/sirupsen/logrus"
	"io"
	"tuic-go/tuic/address"
	"tuic-go/tuic/options"
)

const VersionMajor = 0x05

const (
	CmdAuthenticate = 0x00
	CmdConnect      = 0x01
	CmdPacket       = 0x02
	CmdDissociate   = 0x03
	CmdHeartbeat    = 0x04
)

const (
	HeaderLen       = 2
	PacketLen       = 8
	AuthenticateLen = 48
	DissociateLen   = 2
)

const (
	NetworkTcp = "tcp"
	NetworkUdp = "udp"
)

//goland:noinspection ALL
const (
	UdpRelayModeQuic   = "quic"
	UdpRelayModeNative = "native"
)

const DefaultConcurrentStreams int64 = 32

type Command struct {
	Version uint8
	Type    uint8
	Options options.IOption
}

type PacketResponse struct {
	PacketID uint32
	Data     []byte
}

func (cmd *Command) Marshal() ([]byte, error) {
	// 创建一个字节数组,长度为2(版本号和类型各占1字节)+2(Options长度占2字节)+Options的长度
	totalLen := HeaderLen
	if cmd.Options != nil {
		totalLen += int(cmd.Options.Len())
	}

	cmdBytes := make([]byte, totalLen)

	// 将版本号写入第1个字节
	cmdBytes[0] = cmd.Version

	// 将类型写入第2个字节
	cmdBytes[1] = cmd.Type

	if cmd.Options != nil {
		// 将Options写入剩余的字节
		b, err := cmd.Options.Marshal()
		if err != nil {
			return nil, err
		}

		copy(cmdBytes[2:], b)
	}

	return cmdBytes, nil
}

func (cmd *Command) Unmarshal(stream io.Reader) error {
	header := make([]byte, HeaderLen)
	_, err := io.ReadFull(stream, header)
	if err != nil {
		logrus.Errorf("header io.ReadFull err:%v", err)
		return err
	}

	cmd.Version = header[0]
	cmd.Type = header[1]

	switch cmd.Type {
	case CmdAuthenticate:
		authBytes := make([]byte, AuthenticateLen)
		_, err = io.ReadFull(stream, authBytes)
		if err != nil {
			logrus.Errorf("authBytes io.ReadFull err:%v", err)
			return err
		}

		var opt options.AuthenticateOptions
		err = opt.Unmarshal(authBytes)
		if err != nil {
			logrus.Errorf("opt.Unmarshal err:%v", err)
			return err
		}

		cmd.Options = &opt
	case CmdConnect:
		protocolAddr, err := address.UnMarshalAddr(stream)
		if err != nil {
			logrus.Errorf("address.UnMarshalAddr err:%v", err)
			return err
		}

		var opt options.ConnectOptions
		opt.Addr = protocolAddr
		cmd.Options = &opt
	case CmdPacket:
		packetBytes := make([]byte, PacketLen)
		_, err = io.ReadFull(stream, packetBytes)
		if err != nil {
			logrus.Errorf("packetBytes io.ReadFull err:%v", err)
			return err
		}

		var opt options.PacketOptions
		err = opt.Unmarshal(packetBytes)
		if err != nil {
			logrus.Errorf("opt.Unmarshal err:%v", err)
			return err
		}

		protocolAddr, err := address.UnMarshalAddr(stream)
		if err != nil {
			logrus.Errorf("address.UnMarshalAddr err:%v", err)
			return err
		}

		opt.Addr = protocolAddr

		cmd.Options = &opt

	case CmdDissociate:
		dissociateBytes := make([]byte, DissociateLen)
		_, err = io.ReadFull(stream, dissociateBytes)
		if err != nil {
			logrus.Errorf("dissociateBytes io.ReadFull err:%v", err)
			return err
		}

		var opt options.DissociateOptions
		err = opt.Unmarshal(dissociateBytes)
		if err != nil {
			logrus.Errorf("opt.Unmarshal err:%v", err)
			return err
		}

		cmd.Options = &opt

	case CmdHeartbeat:

	default:
		return errors.New("unknown command type")
	}

	return nil
}
