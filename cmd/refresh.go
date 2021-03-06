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
	"github.com/spf13/viper"
)

var idToSkip = map[int]string{
	404: "This strip is not found on purpose",
}

// refreshCmd represents the refresh command
var refreshCmd = &cobra.Command{
	Use:   "refresh",
	Short: "Download the info about the missing strips.",
	Long: `xkcli refresh will refresh the local database of strips, 
fetching all the  relative metadata.`,
	Run: func(cmd *cobra.Command, args []string) {
		dbPath := viper.GetString("dbPath")
		logger := setupLogging(debugLog).Sugar()
		defer logger.Sync()
		download.SetLogger(logger)
		database.SetLogger(logger)
		// the waitgroup is used to wait for all the goroutines to be done.
		var wg sync.WaitGroup
		logger.Debug("Showing logs at debug level")
		db, err := database.Open(dbPath)
		if err != nil {
			logger.Fatalw("Unable to open the database", "path", dbPath, "error", err)
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
		// Get the max number of records to download
		maxRecords, _ := cmd.Flags().GetInt("maxRecords")

		// Determine which strips to download. We will start from the highest-id
		// strip we have, and add maxRecords new strips.
		logger.Debug("Fetching the most recent ID in the database.")
		lastInDb := database.GetLatestID(db)
		logger.Debugf("Maximum stored ID found: %d", lastInDb)
		logger.Debug("Fetching the latest ID")
		latest := mgr.GetLatestID()
		if maxRecords == 0 {
			maxRecords = latest
		}
		logger.Debugf("Max id is %d", latest)
		toDownload := make([]int, 0)
		// Now search for missing strips in the database
		existingIDs := database.GetAllIDs(db, latest)
		for i := 1; i <= latest; i++ {
			if _, ok := existingIDs[i]; ok {
				continue
			}
			if reason, ok := idToSkip[i]; ok {
				logger.Debugw("Skipping strip", "id", i, "reason", reason)
				continue
			}
			toDownload = append(toDownload, i)
			if len(toDownload) >= maxRecords {
				break
			}
		}
		if len(toDownload) == 0 {
			logger.Info("Nothing to download")
			return
		}
		logger.Info("Downloading strips")

		// download and index data
		for _, id := range toDownload {
			logger.Debugf("Scheduling download of id %d", id)
			wg.Add(1)
			go func(i int, wg *sync.WaitGroup) {
				defer wg.Done()
				w := mgr.Get(i)
				if w != nil {
					doc := database.NewStrip(w)
					err := doc.Index(db)
					if err == nil {
						logger.Infof("Indexed strip %s", doc.Summary())
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
	refreshCmd.Flags().IntP("maxRecords", "m", 0, "Maximum number of records to retreive. By default unbounded.")
}
