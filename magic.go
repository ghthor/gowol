package wol

import (
	"encoding/hex"
	"errors"
	"net"
	"strings"
)

type MagicPacket []byte

// Use a MAC address to form a magic packet
// macAddr form 12:34:56:78:9a:bc
func NewMagicPacket(macAddr string) (MagicPacket, error) {
	if len(macAddr) != (6*2 + 5) {
		return nil, errors.New("Invalid MAC Address String: " + macAddr)
	}

	macBytes, err := hex.DecodeString(strings.Join(strings.Split(macAddr, ":"), ""))
	if err != nil {
		return nil, err
	}

	b := []uint8{255, 255, 255, 255, 255, 255}
	for i := 0; i < 16; i++ {
		b = append(b, macBytes...)
	}

	return MagicPacket(b), nil
}

// Send a Magic Packet to an broadcast class IP address via UDP
func (p MagicPacket) Send(bcastAddr string) error {
	a, err := net.ResolveUDPAddr("udp", net.JoinHostPort(bcastAddr, "7"))
	if err != nil {
		return err
	}

	c, err := net.DialUDP("udp", nil, a)
	if err != nil {
		return err
	}

	written, err := c.Write(p)
	c.Close()

	// Packet must be 102 bytes in length
	if written != 102 {
		return err
	}

	return nil
}

// Constructs a packet and Sends it
func MagicWake(macAddr string, bcastAddr string) error {
	packet, err := NewMagicPacket(macAddr)
	if err != nil {
		return err
	}

	return packet.Send(bcastAddr)
}
