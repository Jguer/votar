package main

import (
	"context"
	"log"

	"github.com/alexflint/go-arg"

	"github.com/Jguer/votar/pkg/vote"
)

func main() {
	var args struct {
		Vote     []string `arg:"-v" help:"packages to vote for"`
		Unvote   []string `arg:"-u" help:"packages to unvote for"`
		User     string   `arg:"env:AUR_USER"`
		Password string   `arg:"env:AUR_PASSWORD"`
	}
	arg.MustParse(&args)

	client, err := vote.NewClient()
	if err != nil {
		log.Println("Failed to create client")
	}
	if len(args.Vote) == 0 && len(args.Unvote) == 0 {
		log.Println("Nothing to do.")
	}

	client.SetCredentials(args.User, args.Password)
	for _, v := range args.Vote {
		err := client.Vote(context.Background(), v)
		if err != nil {
			log.Println("Unable to vote for:", v, "\nerr:", err)
		}
	}

	for _, v := range args.Unvote {
		err := client.Unvote(context.Background(), v)
		if err != nil {
			log.Println("Unable to unvote for:", v, "\nerr:", err)
		}
	}
}
