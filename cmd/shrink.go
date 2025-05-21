package cmd

import (
	"github.com/gkwa/poppingpenguin/internal/app"
	"github.com/gkwa/poppingpenguin/internal/logging"
	"github.com/spf13/cobra"
)

var (
	compressionLevel int
	concurrencyLevel int
	shrinkCmd        = &cobra.Command{
		Use:   "shrink [files...]",
		Short: "Shrink image files",
		Long:  `Shrink image files and display the original size, new size, and shrink percentage.`,
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			logger := logging.NewLogger(verbose)
			shrinker := app.NewImageShrinker(compressionLevel, concurrencyLevel, logger)
			return shrinker.ShrinkImages(args)
		},
	}
)

func init() {
	rootCmd.AddCommand(shrinkCmd)
	shrinkCmd.Flags().IntVarP(&compressionLevel, "level", "l", 80, "compression level (1-100, lower means smaller file)")
	shrinkCmd.Flags().IntVarP(&concurrencyLevel, "concurrency", "c", 4, "number of images to process concurrently")
}
