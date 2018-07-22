// Copied and Modified for NEO from: https://github.com/decred/dcrd/blob/master/wire/VersionMessage.go

package wire

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"math/rand"
	"net"
	"time"
)

const (
	minMsgVersionSize = 28
)

type VersionMessage struct {
	w           *bytes.Buffer
	Version     ProtocolVersion
	Timestamp   time.Time
	Services    ServiceFlag
	IP          net.IP
	Port        uint16
	Nonce       uint32
	UserAgent   []byte
	StartHeight uint32
	Relay       bool
}

var ErrInvalidNetAddr = errors.New("provided net.Addr is not a net.TCPAddr")

func NewVersionMessage(addr net.Addr, startHeight uint32, relay bool, pver ProtocolVersion) (*VersionMessage, error) {
	tcpAddr, ok := addr.(*net.TCPAddr)
	if !ok {
		return nil, ErrInvalidNetAddr
	}
	version := &VersionMessage{
		new(bytes.Buffer),
		pver,
		time.Now(),
		NodePeerService,
		tcpAddr.IP,
		uint16(tcpAddr.Port),
		rand.Uint32(),
		[]byte(UserAgent),
		startHeight,
		relay,
	}

	// saves a buffer of version in version
	if err := version.EncodePayload(version.w); err != nil {
		return nil, err
	}
	return version, nil
}

// Implements Messager interface
func (v *VersionMessage) DecodePayload(r io.Reader) error {
	// Decode into v from reader
	if err := binary.Read(r, binary.LittleEndian, &v.Version); err != nil {
		return err
	}
	if err := binary.Read(r, binary.LittleEndian, &v.Services); err != nil {
		return err
	}
	if err := binary.Read(r, binary.LittleEndian, &v.Timestamp); err != nil {
		return err
	}
	if err := binary.Read(r, binary.LittleEndian, &v.Port); err != nil {
		return err
	}
	if err := binary.Read(r, binary.LittleEndian, &v.Nonce); err != nil {
		return err
	}

	var lenUA uint8
	if err := binary.Read(r, binary.LittleEndian, &lenUA); err != nil {
		return err
	}
	v.UserAgent = make([]byte, lenUA)
	if err := binary.Read(r, binary.LittleEndian, &v.UserAgent); err != nil {
		return err
	}
	if err := binary.Read(r, binary.LittleEndian, &v.StartHeight); err != nil {
		return err
	}
	return binary.Read(r, binary.LittleEndian, &v.Relay)

}

// Implements messager interface
func (v *VersionMessage) EncodePayload(w io.Writer) error {
	// encode into w from v

	bw := &binWriter{w: w}
	bw.Write(v.Version)
	bw.Write(v.Services)
	bw.Write(uint32(v.Timestamp.Unix()))
	bw.Write(v.Port)
	bw.Write(v.Nonce)
	bw.Write(uint8(len(v.UserAgent)))
	bw.Write(v.UserAgent)
	bw.Write(v.StartHeight)
	bw.Write(v.Relay)
	return bw.err
}

// Implements messager interface
func (v *VersionMessage) PayloadLength() uint32 {
	return calculatePayloadLength(v.w)
}

// Implements messager interface
func (v *VersionMessage) Checksum() uint32 {
	return calculateCheckSum(v.w)
}

// Implements messager interface
func (v *VersionMessage) Command() CommandType {
	return CMDVersion
}
