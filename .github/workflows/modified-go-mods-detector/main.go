package main

import (
	"modified-go-mods-detector/cmd"
	"os"

	log "github.com/sirupsen/logrus"
)

func main() {
	err := cmd.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
