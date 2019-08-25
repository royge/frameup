package cmd

import (
	"fmt"
	"log"

	"github.com/royge/frameup/framer"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(cropCmd)
}

var cropCmd = &cobra.Command{
	Use:   "crop",
	Short: "Crop selected pictures inside the source directory.",
	Long:  "Crop selected pictures inside the source directory.",
	Run: func(cmd *cobra.Command, args []string) {
		scanner := framer.Scanner{Delay: delay}

		go func() {
			defer close(inDirChan)
			if err := scanner.ScanDir(&dirWg, src, inDirChan); err != nil {
				log.Fatalf("error scanning %s directory: %v", src, err)
			}
		}()

		for v := range inDirChan {
			c := make(chan string, 4)
			go func(dir string) {
				defer close(c)
				if err := scanner.Scan(&fileWg, dir, c, ext); err != nil {
					log.Fatalf("error scanning %s directory: %v", dir, err)
				}
			}(v)

			for w := range c {
				go func(file string) {

					for _, v := range dims {
						d, _ := framer.ParseDimension(v)
						func(d framer.Dimension) {
							if err := framer.Crop(file, dst, d.Width, d.Height); err != nil {
								fmt.Printf("error cropping picture file %s: %v", file, err)
							}
						}(*d)
					}

					fileWg.Done()
				}(w)
			}

			fileWg.Wait()

			// Done scanning 1 directory.
			dirWg.Done()
		}

		dirWg.Wait()
	},
}
