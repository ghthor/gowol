package wol

import (
	"encoding/hex"
	"log"
	"net"
	"os"
	"strings"
)

// macAddr form 12:34:56:78:9a:bc
func SendMagicPacket(macAddr string, bcastAddr string) os.Error {

	if len(macAddr) != (6*2 + 5) {
		return os.NewError("Invalid MAC Address String: " + macAddr)
	}

	packet, err := constructMagicPacket(macAddr)
	if err != nil {
		return err
	}

	a, err := net.ResolveUDPAddr("udp", bcastAddr+":7")
	if err != nil {
		return err
	}

	c, err := net.DialUDP("udp", nil, a)
	if err != nil {
		return err
	}

	written, err := c.Write(packet)
	c.Close()

	// Packet must be 102 bytes in length
	if written != 102 {
		return err
	}

	return nil
}

func constructMagicPacket(macAddr string) ([]byte, os.Error) {
	macBytes, err := hex.DecodeString(strings.Join(strings.Split(macAddr, ":"), ""))
	if err != nil {
		log.Fatalln("Error Hex Decoding:", err)
		return nil, err
	}

	b := []uint8{255, 255, 255, 255, 255, 255}
	for i := 0; i < 16; i++ {
		b = append(b, macBytes...)
	}
	return b, err
}
