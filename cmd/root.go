/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"riedmann.dev/pocket-cli/cmd/list"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "pocket-cli",
	Short: "A simple CLI for accessing and adding articles to pocket",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand( &cobra.Command{
		Use: "list",
		Short: "Lists your pocket articles", 
		Run: func(cmd *cobra.Command, args []string) {
			list.ListArticles()
		},
	})
}


