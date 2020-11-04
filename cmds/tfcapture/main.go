// The plananizer command analyzes terraform json plan files.
// Currently only supports checking if resources are deleted.
package main

import (
	"flag"
	"log"
	"os"

	"github.com/steeling/tfassist/internal/tfcapture"
)

var (
	stateFile = flag.String("state", "", "terraform json state file.")
	out       = flag.String("out", "", "path to file to write capture non sensitive terraform output")
	overwrite = flag.Bool("overwrite", false, "whether to overwrite if the output file exists")
)

func main() {
	flag.Parse()
	f, err := os.Open(*stateFile)
	if err != nil {
		log.Fatal(err)
	}
	c, err := tfcapture.NewFromState(f)
	if err != nil {
		log.Fatal(err)
	}

	fileFlag := os.O_CREATE | os.O_WRONLY
	if !*overwrite {
		fileFlag |= os.O_EXCL
	}
	outF, err := os.OpenFile(*out, fileFlag, 0600)
	if err != nil {
		log.Fatal(err)
	}

	n, err := c.CaptureOutputs(outF)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("successfully wrote %d outputs to: %s", n, *out)
}
