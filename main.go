package main

import (
	"log"
	"os"

	"github.com/nobelsmith/go-fence/cmd"
)

func main() {
	log.SetOutput(os.Stdout)
	cmd.Execute()
}
