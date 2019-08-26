package cmd

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	bg      string
	overlay string

	delay int64

	dims = dimensions{}
	locs = locations{}

	inDirChan  = make(chan string, 4)
	outDirChan = make(chan string, 4)

	dirWg  sync.WaitGroup
	fileWg sync.WaitGroup

	src string
	dst string
	ext string
)

type dimensions []string
type locations map[string]location

type location struct {
	Top  float64
	Left float64
}

var rootCmd = &cobra.Command{
	Use:   "frameup",
	Short: "Crop and create frame to selected pictures.",
	Long:  `Crop and create frame to selected pictures.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Welcome! run 'frameup --help' for usage.")
	},
}

func configure() {
	d := viper.Get("dimensions").([]interface{})
	dims = make(dimensions, len(d))
	for _, v := range d {
		dims = append(dims, v.(string))
	}

	l := viper.Get("locations").(map[string]interface{})
	for k, v := range l {
		loc := v.(map[string]interface{})
		locs[k] = location{
			Top:  loc["top"].(float64),
			Left: loc["left"].(float64),
		}
	}
}

// Execute commands.
func Execute() {
	configure()

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
