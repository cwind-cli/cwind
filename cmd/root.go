package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "cwind",
	Short: "Cwind - Reportes de ampacidad térmica dinámica (DLR)",
	Long:  `Herramienta para medir la ampacidad de conductores de alta tensión.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Use 'cwind --help' para ver los comandos disponibles")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR] %v\n", err)
		os.Exit(1)
	}
}
