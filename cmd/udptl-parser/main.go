package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"gortc.io/sdp"
)

type UDPTLStream struct {
	SrcIP   string
	SrcPort int
	DstIP   string
	DstPort int
	Codec   string
}
type Type rune

// Define constants for line types
const (
	Version          sdp.Type = 'v'
	Origin           sdp.Type = 'o'
	SessionName      sdp.Type = 's'
	ConnectionData   sdp.Type = 'c'
	Timing           sdp.Type = 't'
	MediaDescription sdp.Type = 'm'
	Attribute        sdp.Type = 'a'
)

type SDPPacket struct {
	Version        int
	Origin         string
	SessionName    string
	ConnectionData string
	Timing         string
	MediaDesc      []string
	Attributes     []string
}

func main() {
	//parse flags and check if filename is provided as an argument. If not - return an error
	// Define a string flag for filename with a default value and a usage description
	filename := flag.String("filename", "", "Path to the pcap file")
	//filterNumber := flag.String("number", "", "numbers to search")

	// Parse the command-line flags
	flag.Parse()

	// Check if filename is provided
	if *filename == "" {
		fmt.Println(`
Error: No filename provided. Use the -filename flag to specify the path to the PCAP file.
Example: udptl-parser -filename dumpfile.pcap
How to get pcap file: tcpdump -s0 -X -vvvnn -w dumpfile.pcap`)
		os.Exit(1)
	}

	// Open the PCAP file
	handle, err := pcap.OpenOffline(*filename)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()

	// List to store UDPTL sessions
	var uniqueSDPSessions []SDPPacket

	// Iterate over packets in the PCAP file
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range packetSource.Packets() {
		// Check if the packet is a SIP packet
		app := packet.ApplicationLayer()
		if app == nil {
			continue
		}

		// Check if the packet contains SDP
		payload := string(app.Payload())
		if !strings.Contains(payload, "m=") {
			continue
		}

		// Parse the SDP content
		session, err := sdp.DecodeSession([]byte(payload), nil)
		if err != nil {
			//log.Println("Error parsing SDP:", err)
			continue
		}

		//// debug
		//for k, v := range session {
		//	fmt.Println(k, v)
		//}

		// Parse the output
		packet := parseSDPPacket(session)
		if packet.Origin == "" || len(packet.MediaDesc) == 0 {
			continue
		}
		//check if the packet is unique. If yes - we add it to the list of streams
		isUnique := checkIfUnique(uniqueSDPSessions, packet)
		if isUnique {
			uniqueSDPSessions = append(uniqueSDPSessions, packet)
		}
	}

	// Print the found UDPTL streams
	fmt.Println("UDPTL Streams:")
	for _, pkt := range uniqueSDPSessions {
		//fmt.Println("Version:", pkt.Version)
		//fmt.Println("Origin:", pkt.Origin)
		//fmt.Println("Session Name:", pkt.SessionName)
		//fmt.Println("Connection Data:", pkt.ConnectionData)
		//fmt.Println("Timing:", pkt.Timing)
		//fmt.Println("Media Description:")
		//for _, desc := range pkt.MediaDesc {
		//	fmt.Println("  -", desc)
		//}
		//fmt.Println("Attributes:")
		//for _, attr := range pkt.Attributes {
		//	fmt.Println("  -", attr)
		//}

		fmt.Printf("From:%v,\t%v\t\n", pkt.Origin, pkt.MediaDesc[0])

	}

}

func parseSDPPacket(session sdp.Session) SDPPacket {
	/*
		0 version: 0
		1 origin: 3906-improcom 8000 8000 IN IP4 172.16.1.165
		2 session name: SIP Call
		3 connection data: IN IP4 172.16.1.165
		4 timing: 0 0
		5 media description: audio 5004 RTP/AVP 0 101
		6 attribute: sendrecv
		7 attribute: rtpmap:0 PCMU/8000
		8 attribute: ptime:20
		9 attribute: rtpmap:101 telephone-event/8000
		10 attribute: fmtp:101 0-16,32-36,54
	*/
	var currentPacket SDPPacket

	for _, line := range session {
		switch line.Type {
		case Version:
			currentPacket.Version = parseInt(string(line.Value))
		case Origin:
			currentPacket.Origin = string(line.Value)
		case SessionName:
			currentPacket.SessionName = string(line.Value)
		case ConnectionData:
			currentPacket.ConnectionData = string(line.Value)
		case Timing:
			currentPacket.Timing = string(line.Value)
		case MediaDescription:
			currentPacket.MediaDesc = append(currentPacket.MediaDesc, string(line.Value))
		case Attribute:
			currentPacket.Attributes = append(currentPacket.Attributes, string(line.Value))
		}
	}

	return currentPacket
}

// parseInt parses an integer from a string
func parseInt(s string) int {
	var num int
	_, err := fmt.Sscanf(s, "%d", &num)
	if err != nil {
		log.Fatal(err)
	}
	return num
}

func checkIfUnique(packets []SDPPacket, packet SDPPacket) bool {
	// Create a map to store unique packets based on Origin and Media Description
	packetMap := make(map[string]bool)

	// Convert slice to map for efficient lookup
	for _, pkt := range packets {
		key := pkt.Origin + "-" + fmt.Sprintf("%v", pkt.MediaDesc)
		packetMap[key] = true
	}

	// Check if the new packet's key already exists in the map
	key := packet.Origin + "-" + fmt.Sprintf("%v", packet.MediaDesc)
	if _, ok := packetMap[key]; ok {
		return false
	}
	return true
}
