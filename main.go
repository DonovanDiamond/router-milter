package main

import (
	"log"
	"net"
	"os"
	routerConfig "router-milter/config"
	"router-milter/session"

	"github.com/DonovanDiamond/milter"
)

var version, commit string

func main() {
	log.Printf("milter-router version %s (commit %s)", version, commit)

	config, err := routerConfig.LoadConfig()
	if err != nil {
		log.Fatal(err.Error())
	}

	script, err := os.ReadFile(config.ScriptPath)
	if err != nil {
		log.Fatalf("Failed to read script from '%s': %v", config.ScriptPath, err)
	}

	// make sure socket does not exist
	if config.Protocol == "unix" {
		// ignore os.Remove errors
		_ = os.Remove(config.Address)
	}

	socket, err := net.Listen(config.Protocol, config.Address)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		_ = socket.Close()
	}()

	log.Printf("Starting milter on %s:%s", config.Protocol, config.Address)

	if config.Protocol == "unix" {
		// set mode 0660 for unix domain sockets
		if err := os.Chmod(config.Address, 0660); err != nil {
			log.Fatal(err)
		}
		// remove socket on exit
		defer func() {
			_ = os.Remove(config.Address)
		}()
	}

	init := func() (milter.Milter, milter.OptAction, milter.OptProtocol) {
		return session.NewSession(string(script)), session.OptAction, session.OptProtocol
	}
	if err := milter.RunServer(socket, init); err != nil {
		log.Fatal(err)
	}
}
