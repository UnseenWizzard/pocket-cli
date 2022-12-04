/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"fmt"
	"github.com/UnseenWizzard/pocket-cli/cmd"
	"os"
)

func main() {
	//cmd.Execute()
	p := cmd.GetProgram()
	if err := p.Start(); err != nil {
		fmt.Printf("Failed to start CLI: %v", err)
		os.Exit(1)
	}
}
