package cmd

import (
	"image"
	"log"
	"os"
	"sync"
	"time"

	"github.com/spf13/cobra"
)

var (
	bg      = "../assets/bg.jpg"
	overlay = "../assets/overlay.png"
	delay   = 100 * time.Millisecond
	dims    = []string{
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
	dirWg      sync.WaitGroup
	fileWg     sync.WaitGroup

	src string
	dst string
	ext string
)

var rootCmd = &cobra.Command{
	Use:   "frameup",
	Short: "Crop and create frame to selected pictures.",
	Long:  `Crop and create frame to selected pictures.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("Welcome! type -h for usage instructions.")
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
	rootCmd.Flags().StringVarP(&src, "source", "s", "", "Source directory.")
	rootCmd.Flags().StringVarP(&dst, "output", "o", "", "Output directory.")
	rootCmd.Flags().StringVarP(&ext, "ext", "e", ".jpg", "Picture files allowed extensions.")

	rootCmd.MarkFlagRequired("source")
	rootCmd.MarkFlagRequired("output")
}
