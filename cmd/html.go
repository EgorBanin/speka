package cmd

import (
	"fmt"
	"os"

	"github.com/egorbanin/speka/speka"
	"github.com/egorbanin/speka/speka/html"
	"github.com/hjson/hjson-go/v4"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(htmlCmd)
}

var htmlCmd = &cobra.Command{
	Use: "html",
	RunE: func(cmd *cobra.Command, args []string) error {
		b, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("os.ReadFile %s: %w", path, err)
		}

		var s speka.Speka
		if err := hjson.Unmarshal(b, &s); err != nil {
			return fmt.Errorf("hjson.Unmarshal: %w", err)
		}

		g, err := html.CreateGenerator()
		if err != nil {
			return fmt.Errorf("html.CreateGenerator: %w", err)
		}

		if err := g.Html(s, os.Stdout); err != nil {
			return fmt.Errorf("g.Html: %w", err)
		}

		return nil
	},
}
