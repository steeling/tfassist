// The plananizer command analyzes terraform json plan files.
// Currently only supports checking if resources are deleted.
package main

import (
	"flag"
	"log"
	"os"

	"github.com/steeling/tfassist/internal/capture"
)

var (
	stateFile = flag.String("state", "", "terraform json state file.")
	out       = flag.String("output", "", "path to file to write capture non sensitive terraform output")
	overwrite = flag.Bool("overwrite", false, "whether to overwrite if the output file exists")
)

func main() {
	flag.Parse()
	f, err := os.Open(*stateFile)
	if err != nil {
		log.Fatal(err)
	}
	c, err := capture.NewFromState(f)
	if err != nil {
		log.Fatal(err)
	}

	fileFlag := os.O_CREATE
	if !*overwrite {
		perm |= os.O_EXCL
	}
	f, err := os.OpenFile(*out, fileFlag, 0600)
	if err != nil {
		log.Fatal(err)
	}

	err := c.CaptureOutputs(f)
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("successfully wrote outputs to: %s", *out)
}
