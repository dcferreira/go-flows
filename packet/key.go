package packet

import (
	"bytes"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"pm.cn.tuwien.ac.at/ipfix/go-flows/flows"
)

// src 4 dst 4 proto 1 src 2 dst 2
type fiveTuple4 [13]byte

func (t fiveTuple4) SrcIP() []byte   { return t[0:4] }
func (t fiveTuple4) DstIP() []byte   { return t[4:8] }
func (t fiveTuple4) Proto() byte     { return t[8] }
func (t fiveTuple4) SrcPort() []byte { return t[9:11] }
func (t fiveTuple4) DstPort() []byte { return t[11:13] }
func (t fiveTuple4) Hash() uint64    { return fnvHash(t[:]) }

// src 16 dst 16 proto 1 src 2 dst 2
type fiveTuple6 [37]byte

func (t fiveTuple6) SrcIP() []byte   { return t[0:16] }
func (t fiveTuple6) DstIP() []byte   { return t[16:32] }
func (t fiveTuple6) Proto() byte     { return t[32] }
func (t fiveTuple6) SrcPort() []byte { return t[33:35] }
func (t fiveTuple6) DstPort() []byte { return t[35:37] }
func (t fiveTuple6) Hash() uint64    { return fnvHash(t[:]) }

var emptyPort = make([]byte, 2)

func fivetuple(packet gopacket.Packet) (flows.FlowKey, bool) {
	network := packet.NetworkLayer()
	if network == nil {
		return nil, false
	}
	transport := packet.TransportLayer()
	if transport == nil {
		return nil, false
	}
	srcPort, dstPort := transport.TransportFlow().Endpoints()
	srcPortR := srcPort.Raw()
	dstPortR := dstPort.Raw()
	proto := transport.LayerType()
	srcIP, dstIP := network.NetworkFlow().Endpoints()
	forward := true
	if dstIP.LessThan(srcIP) {
		forward = false
		srcIP, dstIP = dstIP, srcIP
		if !layers.LayerClassIPControl.Contains(proto) {
			srcPortR, dstPortR = dstPortR, srcPortR
		}
	} else if bytes.Compare(srcIP.Raw(), dstIP.Raw()) == 0 {
		if srcPort.LessThan(dstPort) {
			forward = false
			srcIP, dstIP = dstIP, srcIP
			if !layers.LayerClassIPControl.Contains(proto) {
				srcPortR, dstPortR = dstPortR, srcPortR
			}
		}
	}
	var protoB byte
	switch proto {
	case layers.LayerTypeTCP:
		protoB = byte(layers.IPProtocolTCP)
	case layers.LayerTypeUDP:
		protoB = byte(layers.IPProtocolUDP)
	case layers.LayerTypeICMPv4:
		protoB = byte(layers.IPProtocolICMPv4)
	case layers.LayerTypeICMPv6:
		protoB = byte(layers.IPProtocolICMPv6)
	}
	srcIPR := srcIP.Raw()
	dstIPR := dstIP.Raw()

	if len(srcIPR) == 4 {
		ret := fiveTuple4{}
		copy(ret[0:4], srcIPR)
		copy(ret[4:8], dstIPR)
		ret[8] = protoB
		copy(ret[9:11], srcPortR)
		copy(ret[11:13], dstPortR)
		return ret, forward
	}
	if len(srcIPR) == 16 {
		ret := fiveTuple6{}
		copy(ret[0:16], srcIPR)
		copy(ret[16:32], dstIPR)
		ret[32] = protoB
		copy(ret[33:35], srcPortR)
		copy(ret[35:37], dstPortR)
		return ret, forward
	}
	return nil, false
}