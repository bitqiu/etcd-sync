package cmd

import (
	"etcd-sync/config"
	"etcd-sync/pkg/etcd"
	"log"

	"github.com/spf13/cobra"
)

// syncCmd represents the sync command
var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "etcd sync",
	Run: func(cmd *cobra.Command, args []string) {
		sourceCfg := config.Config{
			Host:     sourceHost,
			Username: username,
			Password: password,
		}

		targetCfg := config.Config{
			Host:     targetHost,
			Username: username,
			Password: password,
		}

		etcdCli := etcd.NewEtcd(sourceCfg, targetCfg)
		err := etcdCli.Sync("/")
		if err != nil {
			log.Fatal("err: ", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)
}
