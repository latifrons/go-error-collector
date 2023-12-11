package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{Use: "app"}

func init() {
	//rootCmd.PersistentFlags().StringP("chain-rpc", "u", "http://localhost:8545", "RPC URL")
	//_ = viper.BindPFlag("chain.rpc", rootCmd.PersistentFlags().Lookup("chain-rpc"))
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
