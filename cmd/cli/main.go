package main

import (
	"github.com/alexflint/go-arg"
	"log"
)

type args struct {
	Since string `arg:"--since" help:"a date filter when applicable"`
}

func main() {
	var args args
	arg.MustParse(&args)

	err := authenticate()
	if err != nil {
		log.Fatalln(err)
		return
	}
}
