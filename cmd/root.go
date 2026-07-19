package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// Version is set at build time via -ldflags.
var Version = "dev"

var rootCmd = &cobra.Command{
	Use:   "meclaw",
	Short: "IM → Agent gateway (Agent Infra)",
	Long: `meclaw is an IM-to-Agent runtime.

Product name: meclaw (claw). Category: Agent Infra / Agent SaaS.
Narrative: IM Agent infrastructure → a work assistant inside WeChat.`,
	Version: Version,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("meclaw scaffold — see docs/ and README_CN.md")
		fmt.Println("next: implement gateway + agent routing (scenario A)")
		return nil
	},
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(Version)
	},
}
