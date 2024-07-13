package api_handler

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
	log.SetPrefix("Server: ")
	log.Printf("starting http server, Version %s, GIT sha %s", semVer, gitSha)

	// TODO: Implement the rest of the main function
}
