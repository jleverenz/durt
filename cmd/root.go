package cmd

import (
	"fmt"
	"os"
	"regexp"

	"github.com/jleverenz/durt/core"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "durt",
	Short:   "durt - disk utilization reporting tool",
	Version: "0.0.1",
	// TODO Long:
	Run: func(cmd *cobra.Command, args []string) {
		for _, exc := range exclusionStrings {
			core.GlobalOpts.Exclusions = append(core.GlobalOpts.Exclusions, regexp.MustCompile(exc))
		}
		core.MainAction(args)
	},
}

var exclusionStrings []string

func Execute() {
	rootCmd.PersistentFlags().BoolVar(&core.GlobalOpts.Head, "head", false, "display the top 20 records (default: false)")
	rootCmd.PersistentFlags().StringArrayVar(&exclusionStrings, "exclude", []string{}, "exclude paths by regex")
	rootCmd.Args = cobra.ArbitraryArgs
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
