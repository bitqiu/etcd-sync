package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var sourceHost, targetHost, username, password string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "etcd-sync",
	Short: "sync etcd",
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&sourceHost, "source", "s", "", "Source Etcd host (required)")
	rootCmd.PersistentFlags().StringVarP(&targetHost, "target", "t", "", "Target Etcd host (required)")
	rootCmd.PersistentFlags().StringVarP(&username, "username", "u", "", "Username for authentication")
	rootCmd.PersistentFlags().StringVarP(&password, "password", "p", "", "Password for authentication")
}
