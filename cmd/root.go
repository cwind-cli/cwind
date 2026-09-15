package cmd

import (
	"fmt"
	"os"

	"github.com/cwind-cli/cwind/internal/version"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "cwind",
	Short: "Cwind - Reportes de ampacidad térmica dinámica (DLR)",
	Long:  `Herramienta para medir la ampacidad de conductores de alta tensión.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Fprintln(cmd.OutOrStdout(), "Use 'cwind --help' para ver los comandos disponibles")
		return nil
	},
	Version: version.Value,
}

func Execute() error {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR] %v\n", err)
		return err
	}
	return nil
}
