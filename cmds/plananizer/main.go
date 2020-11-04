// The plananizer command analyzes terraform json plan files.
// Currently only supports checking if resources are deleted.
package main

import (
	"flag"
	"log"
	"os"

	"github.com/steeling/tfassist/internal/tfplananizer"
)

var (
	planFile      = flag.String("plan", "", "terraform json plan file.")
	failOnDeletes = flag.Bool("exit_on_deletes", false, "error out if deletes are detected")
)

func main() {
	flag.Parse()
	f, err := os.Open(*planFile)
	if err != nil {
		log.Fatal(err)
	}
	tfp, err := tfplananizer.NewFromPlan(f)
	if err != nil {
		log.Fatal(err)
	}

	if destroys, r := tfp.Destroys(); *failOnDeletes && destroys {
		log.Fatalf("plan: %s will destroy the following resources: %v", *planFile, r)
	}
	log.Print("no deletions detected")
}
