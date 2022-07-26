/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "mokuro-vox",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// RunE: func(cmd *cobra.Command, args []string) error {
	// 	return mokurovox.GenerateAudio(volumeFile, speaker, overrideHtml)
	// },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func checkVoicevoxConnection() error {
	_, err := http.Get("http://localhost:50021/version")

	return err
}
