package vote_test

import (
	"context"
	"log"

	"github.com/Jguer/votar/pkg/vote"
)

// Vote for a single package
func Example() {
	client, err := vote.NewClient()
	if err != nil {
		log.Println("Failed to create client")
	}

	client.SetCredentials("user", "password")

	if err = client.Vote(context.Background(), "package"); err != nil {
		log.Println("Failed to vote for", "package")
	}
}
