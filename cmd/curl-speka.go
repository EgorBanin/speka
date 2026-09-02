package cmd

import (
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/egorbanin/speka/speka"
	"github.com/hjson/hjson-go/v4"
	"github.com/spf13/cobra"
)

const (
	defaultMethod = "POST"
)

var (
	domain  string
	urlPath string
)

func init() {
	rootCmd.AddCommand(curlSpeka)
	curlSpeka.PersistentFlags().StringVar(&domain, "domain", "http://localhost:8080", "request domain address")
	curlSpeka.PersistentFlags().StringVar(&urlPath, "url-path", "/", "request url path")
}

var curlSpeka = &cobra.Command{
	Use: "curl",
	RunE: func(cmd *cobra.Command, args []string) error {
		b, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("os.ReadFile %s: %w", path, err)
		}

		var s speka.Speka
		if err := hjson.Unmarshal(b, &s); err != nil {
			return fmt.Errorf("hjson.Unmarshal: %w", err)
		}

		u, err := url.Parse(domain)
		if err != nil {
			return fmt.Errorf("url.Parse: %w", err)
		}
		u = u.JoinPath(urlPath)

		method, ok := s.Methods[urlPath]
		if !ok {
			return fmt.Errorf("urlPath %s not found", urlPath)
		}
		jsonRq, err := method.Rq.MarshalJSON()

		var strBuilder strings.Builder
		strBuilder.WriteString(fmt.Sprintf("curl -X %s --location \"%s\" \\\n", defaultMethod, u.String()))
		strBuilder.WriteString(fmt.Sprintf("\t-H \"Content-Type: application/json\" \\\n"))
		strBuilder.WriteString(fmt.Sprintf("\t-d '%s'\n", string(jsonRq)))

		if _, err := fmt.Fprintf(os.Stdout, strBuilder.String()); err != nil {
			return err
		}

		return nil
	},
}
