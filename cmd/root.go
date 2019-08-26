package cmd

import (
	"fmt"
	"image"
	"log"
	"os"
	"sync"

	"github.com/spf13/cobra"
)

var (
	bg      string
	overlay string

	delay int64

	dims = []string{
		"1200x1800",
		"460x920",
		"460x880",
	}
	locations = map[string]image.Point{
		"1200x1800": image.Pt(-160, 0),
		"460x920":   image.Pt(0, -880),
		"460x880":   image.Pt(0, 0),
	}

	inDirChan  = make(chan string, 4)
	outDirChan = make(chan string, 4)

	dirWg  sync.WaitGroup
	fileWg sync.WaitGroup

	src string
	dst string
	ext string
)

var rootCmd = &cobra.Command{
	Use:   "frameup",
	Short: "Crop and create frame to selected pictures.",
	Long:  `Crop and create frame to selected pictures.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Welcome! run 'frameup --help' for usage.")
	},
}

// Execute commands.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&src, "source", "s", "", "Source directory.")
	rootCmd.PersistentFlags().StringVarP(&dst, "output", "o", "", "Output directory.")
	rootCmd.PersistentFlags().StringVarP(&ext, "ext", "e", ".jpg", "Picture files allowed extensions.")
	rootCmd.PersistentFlags().StringVarP(&bg, "bg", "b", "./assets/bg.jpg", "Background image.")
	rootCmd.PersistentFlags().StringVarP(&overlay, "overlay", "l", "./assets/overlay.jpg", "Overlay or frame image.")
	rootCmd.PersistentFlags().Int64VarP(&delay, "delay", "d", 100, "Delay.")

	rootCmd.MarkFlagRequired("source")
	rootCmd.MarkFlagRequired("output")
}
