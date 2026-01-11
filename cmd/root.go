package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "grimleytk",
	Short: "GrimleyTK - Declarative Data Architecture Toolkit",
	Long:  "GrimleyTK is a CLI toolkit for validating and evolving data architectures.",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// aqui futuramente entram flags globais
}
