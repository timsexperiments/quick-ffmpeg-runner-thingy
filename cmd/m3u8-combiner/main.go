package main

import (
	"os"

	"github.com/timsexperiments/quick-ffmpeg-runner-thingy/internal/log"

	"github.com/spf13/cobra"
	"github.com/timsexperiments/quick-ffmpeg-runner-thingy/internal/combiner"
)

var (
	sequence     int
	outputFormat string
	output       string
	verbose      bool
	audio        bool
)

var rootCmd = &cobra.Command{
	Use:   "m3u8-combiner [flags] <m3u8_file>",
	Short: "Combine M3U8 files into a single video file",
	Long:  "m3u8-combiner is a CLI tool that combines M3U8 files into a single video file. It allows you to specify the output name, sequence number, and video format.",
	Args:  cobra.ExactArgs(1),
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		log.InitLogger(verbose)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		path := args[0]

		cfg := combiner.CombinerConfig{
			Path:         path,
			Sequence:     sequence,
			OutputFormat: outputFormat,
			Output:       output,
			Audio:        audio,
		}

		log.L().Debug("Starting combiner",
			"path", cfg.Path,
			"sequence", cfg.Sequence,
			"format", cfg.OutputFormat,
			"output", cfg.Output,
		)

		return combiner.RunCombiner(cfg)
	},
	Version: "0.0.2",
}

func init() {
	rootCmd.Flags().StringVarP(&output, "output", "o", "combined", "The base name for the output file")
	rootCmd.Flags().IntVarP(&sequence, "sequence", "s", 1, "Sequence number to append to the output file")
	rootCmd.Flags().StringVarP(&outputFormat, "format", "f", "mp4", "Output video format (e.g., mp4, mkv)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")
	rootCmd.Flags().BoolVarP(&audio, "audio", "a", false, "Extract a separate audio file in the output")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.L().Error("Error executing command", "error", err)
		os.Exit(1)
	}
}
