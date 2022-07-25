/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"os"

	"github.com/luizbafilho/mokuro-vox/mokurovox"

	"github.com/spf13/cobra"
)

var (
	htmlPath string
	speaker  string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "mokuro-vox",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	RunE: func(cmd *cobra.Command, args []string) error {
		return mokurovox.GenerateAudio(htmlPath, speaker)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&htmlPath, "html-path", "", "Volume html file")
	rootCmd.PersistentFlags().StringVar(&speaker, "speaker", "", "Speaker ID")
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
}
