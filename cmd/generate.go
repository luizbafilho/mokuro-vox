/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"errors"

	"github.com/luizbafilho/mokuro-vox/mokurovox"
	"github.com/spf13/cobra"
)

var (
	volumeFile   string
	speaker      string
	overrideFile bool
	speed        float64
)

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:          "generate",
	Short:        "Generates audio files for Mokuro text boxes.",
	SilenceUsage: true,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if err := checkVoicevoxConnection(); err != nil {
			return errors.New("VoiceVox is not running!")
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return mokurovox.GenerateAudio(volumeFile, speaker, overrideFile, speed)
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)

	generateCmd.PersistentFlags().StringVar(&volumeFile, "volume-file", "", "Volume html file (Required)")
	generateCmd.PersistentFlags().StringVar(&speaker, "speaker", "", "Speaker ID (Required)")
	generateCmd.PersistentFlags().Float64Var(&speed, "speed", 1, "Controls the playback speed")
	generateCmd.PersistentFlags().BoolVar(&overrideFile, "override-file", false, "Updates Volume file")

	generateCmd.MarkPersistentFlagRequired("volume-file")
	generateCmd.MarkPersistentFlagRequired("speaker")
}
