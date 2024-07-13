package validator

import (
	"log"
)

var (
	gitSha = "unknown" // Populated with the last Git commit SHA (short) at build time
	semVer = "unknown" // Populated with semantic version at build time
)

func main() {

	// Define standards
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetPrefix("Validator: ")
	log.Printf("starting validator, Version %s, GIT sha %s", semVer, gitSha)

	// TODO: Implement the rest of the main function
}
