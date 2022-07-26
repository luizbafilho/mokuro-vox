/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"errors"
	"net/http"

	"github.com/luizbafilho/mokuro-vox/mokurovox"
	"github.com/spf13/cobra"
)

var (
	volumeFile   string
	speaker      string
	overrideFile bool
)

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if err := checkVoicevoxConnection(); err != nil {
			return errors.New("VoiceVox is not running!")
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return mokurovox.GenerateAudio(volumeFile, speaker, overrideFile)
	},
}

func checkVoicevoxConnection() error {
	_, err := http.Get("https://localhost/50021/version")

	return err
}

func init() {
	rootCmd.AddCommand(generateCmd)

	generateCmd.PersistentFlags().StringVar(&volumeFile, "volume-file", "", "Volume html file")
	generateCmd.PersistentFlags().StringVar(&speaker, "speaker", "", "Speaker ID")
	generateCmd.PersistentFlags().BoolVar(&overrideFile, "override-file", false, "Updates Volume file")

	generateCmd.MarkPersistentFlagRequired("volume-file")
	generateCmd.MarkPersistentFlagRequired("speaker")
}
