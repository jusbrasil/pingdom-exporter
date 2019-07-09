package main

import (
	"log"
	"os"

	"github.com/veepee-moc/pingdom_exporter/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		log.Print(err)
		os.Exit(1)
	}
}
