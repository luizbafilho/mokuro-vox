/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"errors"
	"fmt"

	"github.com/luizbafilho/mokuro-vox/mokurovox"
	"github.com/spf13/cobra"
)

// listSpeakersCmd represents the listSpeakers command
var listSpeakersCmd = &cobra.Command{
	Use:   "list-speakers",
	Short: "Lists all Speakers available in VoiceVox",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if err := checkVoicevoxConnection(); err != nil {
			return errors.New("VoiceVox is not running!")
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		speakers, err := mokurovox.ListSpeakers()
		if err != nil {
			return err
		}

		fmt.Printf("VoiceVox Speakers:\n\n")
		for _, s := range speakers {
			fmt.Printf("- Name: %s\n", s.Name)
			fmt.Printf("  Styles:\n")
			for _, style := range s.Styles {
				fmt.Printf("  - Name: %s\n", style.Name)
				fmt.Printf("    Id: %d\n", style.ID)
			}
			fmt.Println()
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(listSpeakersCmd)
}
