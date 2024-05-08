[![build](https://github.com/olegromanchuk/udptl-parser/actions/workflows/build.yml/badge.svg)](https://github.com/olegromanchuk/udptl-parser/actions/workflows/build.yml/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/olegromanchuk/udptl-parser)](https://goreportcard.com/report/github.com/olegromanchuk/udptl-parser)
[![Go Reference](https://pkg.go.dev/badge/github.com/olegromanchuk/udptl-parser.svg)](https://pkg.go.dev/github.com/olegromanchuk/udptl-parser)
[![GitHub release (release name instead of tag name)](https://img.shields.io/github/v/release/olegromanchuk/udptl-parser)](https://github.com/olegromanchuk/udptl-parser/releases)

[//]: # ([![Coverage Status]&#40;https://coveralls.io/repos/github/olegromanchuk/udptl-parser/badge.svg?branch=master&#41;]&#40;https://coveralls.io/github/olegromanchuk/udptl-parser?branch=master&#41;)

# Udptl / RTP stream parser
Helpful in fax troubleshooting.  
The util parses traffic dump .pcap file and displays udptl/rtp streams. Wireshark does not display udptl streams properly.  
You need to supply traffic dump file in pcap format.   
Use `tcpdump -s0 -X -vvvnn -w dumpfile.pcap>` - without any filters. Be careful on production systems - filesize will grow fast.  
Use `tcpdump \(host 1.1.1.1 or host 2.2.2.2 or net 5.5.5.0/24 \) -s0 -X -vvvnn -w dumpfile` with  host filters. Make sure to include media servers networks as they are usually different from signaling networks.

## Getting Started
* Install udptl-parser by downloading the latest release from the [releases](https://github.com/olegromanchuk/udptl-parser/releases) page.
* run `udptl-parser -filename dumpfile.pcap` to parse the dump file and display all udptl/t.38/rtp streams.
* To add filters by numbers use `udptl-parser -filename dumpfile.pcap -number 2125554444`

