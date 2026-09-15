package cmd

import (
	"fmt"
	"os"
	"sort"
	"text/tabwriter"

	"github.com/cwind-cli/cwind/internal/conductors"
	"github.com/spf13/cobra"
)

var conductorsCmd = &cobra.Command{
	Use:   "conductors",
	Short: "Gestiona la base de datos de conductores",
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Lista todos los conductores disponibles",
	Run: func(cmd *cobra.Command, args []string) {
		list := conductors.List()
		sort.Slice(list, func(i, j int) bool {
			return list[i].ID < list[j].ID
		})

		writer := tabwriter.NewWriter(os.Stdout, 0, 4, 2, ' ', 0)
		fmt.Fprintln(writer, "ID\tCONDUCTOR\tAMPACIDAD ESTÁTICA")
		for _, c := range list {
			fmt.Fprintf(writer, "%s\t%s\t%.0f A\n", c.ID, c.Name, c.StaticRating)
		}
		_ = writer.Flush()
	},
}

func init() {
	rootCmd.AddCommand(conductorsCmd)
	conductorsCmd.AddCommand(listCmd)
}
