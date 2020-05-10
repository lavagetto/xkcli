/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"sync"

	"github.com/lavagetto/xkcli/database"
	"github.com/lavagetto/xkcli/download"
	"github.com/spf13/cobra"
)

// refreshCmd represents the refresh command
var refreshCmd = &cobra.Command{
	Use:   "refresh",
	Short: "Download the info about the missing strips.",
	Long: `xkcli refresh will refresh the local database of strips, 
fetching all the  relative metadata.`,
	Run: func(cmd *cobra.Command, args []string) {
		logger := setupLogging(debugLog).Sugar()
		defer logger.Sync()
		download.SetLogger(logger)
		logger.Debug("Showing logs at debug level")
		db, err := database.Open("xkcd.bleve", logger)
		if err != nil {
			logger.Fatalw("Unable to open the database", "error", err)
		}
		defer db.Close()
		// Setup the download manager
		c, _ := cmd.Flags().GetInt("concurrency")
		bus := make(chan struct{}, c)
		ua, _ := cmd.Flags().GetString("userAgent")
		mgr := download.Manager{
			Bus: bus,
			Ua:  ua,
		}
		defer mgr.Close()
		logger.Debug("Fetching the latest ID")
		latest := mgr.GetLatestID()
		logger.Debugf("Max id is %d", latest)
		var wg sync.WaitGroup
		for id := 1; id <= latest; id++ {
			logger.Debugf("Scheduling download of id %d", id)
			wg.Add(1)
			go func(i int, wg *sync.WaitGroup) {
				defer wg.Done()
				w := mgr.Get(i)
				if w != nil {
					err := w.Index(db)
					if err == nil {
						logger.Infof("Indexed strip %s", w.Summary())
					}
				}
			}(id, &wg)
		}
		wg.Wait()
	},
}

func init() {
	rootCmd.AddCommand(refreshCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// refreshCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// refreshCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	refreshCmd.Flags().IntP("concurrency", "c", 1, "Number of parallel threads to launch to download missing strips")
	refreshCmd.Flags().StringP("userAgent", "u", "XKCD-cli Crawler/1.0.0", "The user-agent to use when downloading the contents.")
}
