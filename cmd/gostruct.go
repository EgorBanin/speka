package cmd

import (
	"fmt"
	"github.com/egorbanin/speka/speka"
	"github.com/egorbanin/speka/speka/generator"
	"os"

	"github.com/hjson/hjson-go/v4"
	"github.com/spf13/cobra"
)

var pckg string

func init() {
	rootCmd.AddCommand(goStruct)
	rootCmd.PersistentFlags().StringVar(&pckg, "package", "", "package name")
}

var goStruct = &cobra.Command{
	Use: "gostruct",
	RunE: func(cmd *cobra.Command, args []string) error {
		b, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("os.ReadFile %s: %w", path, err)
		}

		var s speka.Speka
		if err := hjson.Unmarshal(b, &s); err != nil {
			return fmt.Errorf("hjson.Unmarshal: %w", err)
		}

		p := s.Name
		if pckg != "" {
			p = pckg
		}

		gen := generator.NewGoStruct(p)
		for name, m := range s.Methods {
			p, err := speka.ParseProperty(fmt.Sprintf("%s_rq", name), m.Rq)
			if err != nil {
				return fmt.Errorf("speka.ParseProperty: %w", err)
			}

			if err := gen.Generate(p, os.Stdout, generator.GoStructOpts{
				Validator: true,
			}); err != nil {
				return fmt.Errorf("gen.Generate: %w", err)
			}

			p, err = speka.ParseProperty(fmt.Sprintf("%s_rs", name), m.Rs)
			if err != nil {
				return fmt.Errorf("speka.ParseProperty: %w", err)
			}

			if err := gen.Generate(p, os.Stdout, generator.GoStructOpts{}); err != nil {
				return fmt.Errorf("gen.Generate: %w", err)
			}
		}

		return nil
	},
}
