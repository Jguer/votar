package main

import (
	"fmt"
	"log"

	"github.com/Jguer/votar/pkg/vote"
	"github.com/alexflint/go-arg"
)

func main() {
	var args struct {
		Vote     []string `arg:"-v" help:"packages to vote for"`
		Unvote   []string `arg:"-u" help:"packages to unvote for"`
		User     string   `arg:"env:AUR_USER"`
		Password string   `arg:"env:AUR_PASSWORD"`
	}
	arg.MustParse(&args)
	fmt.Println(args.Vote, args.Unvote)

	client, err := vote.NewClient(nil, nil)
	if err != nil {
		log.Println("Failed to create client")
	}
	if len(args.Vote) == 0 && len(args.Unvote) == 0 {
		log.Println("Nothing to do.")
	}

	client.SetCredentials(args.User, args.Password)
	for _, v := range args.Vote {
		err := client.Vote(v)
		if err != nil {
			log.Println("Unable to vote for", v, err)
		}
	}

	for _, v := range args.Unvote {
		err := client.Unvote(v)
		if err != nil {
			log.Println("Unable to unvote for", v)
		}
	}
}
