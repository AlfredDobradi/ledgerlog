package main

import "log"

type FindKeyCmd struct {
	Key string `arg:"" help:"Pattern to scan for" required:""`
}

func (cmd *FindKeyCmd) Run(ctx *Context) error {
	log.Println(cmd.Key)
	return nil
}
