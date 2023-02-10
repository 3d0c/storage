package node

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	version string
	date    string
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "full version of storage proxy component",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("version=v%s, date=%s\n", version, date)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
