package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var path string

var rootCmd = &cobra.Command{
	Use: "speka",
}

func init() {
	rootCmd.PersistentFlags().StringVar(&path, "file", "speka.json5", "path to speka file")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err.Error())
	}
}
