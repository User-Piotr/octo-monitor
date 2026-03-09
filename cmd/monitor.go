package cmd

import (
	"context"
	"fmt"
	"sync"
	"time"

	monitor "octo-monitor/internal"

	"github.com/spf13/cobra"
)

var timeout int

func worker(ctx context.Context, jobs <-chan string, results chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()
	for url := range jobs {
		status := monitor.RunMonitoring(ctx, url)
		results <- fmt.Sprintf("%s: %s", status, url)
	}
}

var monitorCmd = &cobra.Command{
	Use:   "monitor",
	Short: "Health check URLs",
	Args:  cobra.RangeArgs(1, 5),
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
		defer cancel()

		tasks := make(chan string, len(args))
		results := make(chan string, len(args))
		numWorkers := len(args)

		for _, url := range args {
			if url != "" {
				tasks <- url
			}
		}
		close(tasks)

		var wg sync.WaitGroup

		for range numWorkers {
			wg.Add(1)
			go worker(ctx, tasks, results, &wg)
		}

		var printerWg sync.WaitGroup
		printerWg.Go(func() {
			for result := range results {
				fmt.Println(result)
			}
		})

		wg.Wait()
		close(results)
		printerWg.Wait()
	},
}

func init() {
	rootCmd.AddCommand(monitorCmd)
	monitorCmd.Flags().IntVar(&timeout, "timeout", 5, "timeout in seconds")
}
