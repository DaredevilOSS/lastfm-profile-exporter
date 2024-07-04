package main

import (
	"github.com/alexflint/go-arg"
	"lastfm-profile-exporter/internal"
)

type args struct {
	Usernames []string `arg:"--usernames,required" help:"usernames to export data for"`
	Since     string   `arg:"--since" help:"a date filter when applicable"`
}

func main() {
	var args args
	arg.MustParse(&args)

	for _, username := range args.Usernames {
		internal.Collect(username)
	}
}
