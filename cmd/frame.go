package cmd

import (
	"log"
	"sync"

	"github.com/royge/frameup/framer"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(frameCmd)
}

var frameCmd = &cobra.Command{
	Use:   "frame",
	Short: "Create frame selected pictures.",
	Long:  "Create frame selected pictures.",
	Run: func(cmd *cobra.Command, args []string) {
		scnr := framer.Scanner{Delay: delay}
		fmr := framer.Framer{
			Dims:      dims,
			Locations: locations,
			Bg:        bg,
			Overlay:   overlay,
		}

		go func() {
			defer close(outDirChan)
			if err := scnr.ScanDir(&dirWg, dst, outDirChan); err != nil {
				log.Fatalf("error scanning %s directory: %v", src, err)
			}
		}()

		for v := range outDirChan {
			c := make(chan string, 1)
			files := []string{}
			mu := sync.Mutex{}

			go func(dir string) {
				defer close(c)
				if err := scnr.Scan(&fileWg, dir, c, ext); err != nil {
					log.Fatalf("error scanning %s directory: %v", dir, err)
				}
			}(v)

			for f := range c {
				mu.Lock()
				files = append(files, f)
				mu.Unlock()
				fileWg.Done()
			}

			m := framer.Classify(files, dims)
			err := fmr.Frame(m, dst)
			if err != nil {
				log.Fatalf("error creating frame: %v", err)
			}

			fileWg.Wait()

			// Done scanning 1 directory.
			dirWg.Done()
		}

		dirWg.Wait()
	},
}
