package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// helloCmd represents the hello command
var selectCmd = &cobra.Command{
	Use:   "select",
	Short: "Cluster selection",
	Long:  `Select an openshift cluster`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Selecting ...")
	},
}

func init() {
	RootCmd.AddCommand(selectCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// helloCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// helloCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
