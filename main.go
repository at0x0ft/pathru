/*
Copyright © 2023 at0x0ft <26642966+at0x0ft@users.noreply.github.com>
*/
package main

import (
    "github.com/at0x0ft/pathru/pkg/cmd"
	"log"
	"os"
)

func main() {
	if err := cmd.NewRootCommand().Execute(); err != nil {
		log.Fatalf("%v\n", err.Error())
		os.Exit(1)
	}
}
