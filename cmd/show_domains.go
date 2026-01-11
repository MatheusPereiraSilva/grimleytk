package cmd

import (
	"fmt"

	"grimleytk/internal/config"

	"github.com/spf13/cobra"
)

var showDomainsCmd = &cobra.Command{
	Use:   "domains",
	Short: "Show defined domains",
	Run:   showDomains,
}

func init() {
	showCmd.AddCommand(showDomainsCmd)
}

func showDomains(cmd *cobra.Command, args []string) {
	cfg, err := config.Load("grimley.yaml")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Domains:")
	for name, domain := range cfg.Domains {
		fmt.Printf("- %s\n", name)
		fmt.Printf("  Schema: %s\n", domain.Schema)
		fmt.Printf("  Owner : %s\n", domain.Owner)

		if domain.Owns != nil && len(domain.Owns.Tables) > 0 {
			fmt.Printf("  Owns  : %d table(s)\n", len(domain.Owns.Tables))
		} else {
			fmt.Printf("  Owns  : none\n")
		}

		if len(domain.Reads) > 0 {
			fmt.Printf("  Reads : %d read model(s)\n", len(domain.Reads))
		} else {
			fmt.Printf("  Reads : none\n")
		}

		fmt.Println()
	}
}
