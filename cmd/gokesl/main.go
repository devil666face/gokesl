package main

import (
	"log"
	"os"

	_ "embed"

	"gokesl/internal/gokesl"
)

//go:embed klnagent64_15.1.0-20748_amd64.deb
var AgentBin []byte

//go:embed kesl-astra_11.1.0-3013.mod_amd64.deb
var KeslBin []byte

//go:embed ws.key
var KeyFile []byte

var (
	// Os        = gokesl.Redhat
	Os        = gokesl.Debian
	KscIP     = ""
	UpdateURI = ""
)

func main() {
	_gokesl, err := gokesl.New(
		&AgentBin, &KeslBin, &KeyFile,
		Os,
		KscIP,
		UpdateURI,
	)
	if err != nil {
		log.Fatalln(err)
	}
	if err := _gokesl.Install(); err != nil {
		log.Fatalln(err)
	}
	os.Exit(0)
}
