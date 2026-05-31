package main

import (
	"flag"
	"log"
	"net"
	"os"
	"strings"

	"github.com/DonovanDiamond/milter"
)

var version, commit string

var protocol = flag.String("proto", "unix", "Protocol family (unix or tcp)")
var address = flag.String("addr", "/var/spool/postfix/milter/router.sock", "Bind to address or unix domain socket")
var rejectFrom = flag.String("reject-from", "", "Comma-separated list of rejected sender addresses")
var rejectTo = flag.String("reject-to", "", "Comma-separated list of rejected recipient addresses")
var configPath = flag.String("config", "", "Path to configuration file (yaml)")

func main() {
	log.Printf("milter-router version %s (commit %s)", version, commit)
	// parse commandline arguments
	flag.Parse()

	rejectedFrom := make(map[string]bool)
	rejectedTo := make(map[string]bool)

	// load config file if specified
	if *configPath != "" {
		cfg, err := LoadConfig(*configPath)
		if err != nil {
			log.Fatalf("error loading config: %v", err)
		}
		if cfg.Protocol != "" {
			*protocol = cfg.Protocol
		}
		if cfg.Address != "" {
			*address = cfg.Address
		}
		for _, addr := range cfg.RejectFrom {
			rejectedFrom[strings.Trim(addr, "<>")] = true
		}
		for _, addr := range cfg.RejectTo {
			rejectedTo[strings.Trim(addr, "<>")] = true
		}
	}

	if *rejectFrom != "" {
		for addr := range strings.SplitSeq(*rejectFrom, ",") {
			addr = strings.TrimSpace(addr)
			addr = strings.Trim(addr, "<>")
			rejectedFrom[addr] = true
		}
	}
	if *rejectTo != "" {
		for addr := range strings.SplitSeq(*rejectTo, ",") {
			addr = strings.TrimSpace(addr)
			addr = strings.Trim(addr, "<>")
			rejectedTo[addr] = true
		}
	}

	// make sure the specified protocol is either unix or tcp
	if *protocol != "unix" && *protocol != "tcp" {
		log.Fatal("invalid protocol name")
	}

	// make sure socket does not exist
	if *protocol == "unix" {
		// ignore os.Remove errors
		_ = os.Remove(*address)
	}

	// bind to listening address
	socket, err := net.Listen(*protocol, *address)
	if err != nil {
		log.Fatal(err)
	}
	defer socket.Close()

	log.Printf("Starting milter on %s:%s", *protocol, *address)

	if *protocol == "unix" {
		// set mode 0660 for unix domain sockets
		if err := os.Chmod(*address, 0660); err != nil {
			log.Fatal(err)
		}
		// remove socket on exit
		defer os.Remove(*address)
	}

	init := func() (milter.Milter, milter.OptAction, milter.OptProtocol) {
		return &RouterMilter{
				rejectedFrom: rejectedFrom,
				rejectedTo:   rejectedTo,
			},
			milter.OptAddHeader | milter.OptChangeHeader,
			milter.OptNoBody | milter.OptNoHeaders | milter.OptNoEOH
	}
	if err := milter.RunServer(socket, init); err != nil {
		log.Fatal(err)
	}
}
