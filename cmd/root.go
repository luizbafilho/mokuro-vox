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
	volumeFile   string
	speaker      string
	overrideHtml bool
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
		return mokurovox.GenerateAudio(volumeFile, speaker, overrideHtml)
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
	rootCmd.PersistentFlags().StringVar(&volumeFile, "volume-file", "", "Volume html file")
	rootCmd.PersistentFlags().StringVar(&speaker, "speaker", "", "Speaker ID")
	rootCmd.PersistentFlags().BoolVar(&overrideHtml, "override-html", false, "Updates Volume file")

	rootCmd.MarkPersistentFlagRequired("volume-file")
	rootCmd.MarkPersistentFlagRequired("speaker")
}
