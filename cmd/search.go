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
	"fmt"
	"os"

	"github.com/blevesearch/bleve/search"
	"github.com/spf13/viper"

	"github.com/lavagetto/xkcli/database"
	"github.com/spf13/cobra"
)

// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search XKCD strips",
	Long: `Search the xkcd database for a string:

Don't forget to quote your query on the shell.`,
	Run: func(cmd *cobra.Command, args []string) {
		dbPath := viper.GetString("dbPath")
		minScore := viper.GetFloat64("minScore")
		doFeelLucky, err := cmd.Flags().GetBool("lucky")
		if err != nil {
			panic(err)
		}
		db, err := database.Open(dbPath)
		if err != nil {
			panic(err)
		}
		defer db.Close()
		searchResult, err := database.SearchStr(db, args[0], nil)
		if err != nil {
			panic(err)
		}
		if searchResult.MaxScore < minScore {
			fmt.Println("No result matching the query abouve minimum score")
			os.Exit(1)
		}
		if doFeelLucky {
			strip := database.NewStripFromDb(searchResult.Hits[0])
			fmt.Println(strip.URL())
			os.Exit(0)
		}
		skipped := 0
		var validResults []*search.DocumentMatch
		for _, result := range searchResult.Hits {
			if result.Score >= minScore {
				validResults = append(validResults, result)
			} else {
				skipped++
			}
		}
		fmt.Println("Your search results:")
		for pos, result := range validResults {
			strip := database.NewStripFromDb(result)
			fmt.Printf("%d - (%.2f) %s", pos, result.Score, strip.Summary())
		}
		if skipped > 0 {
			fmt.Printf("We also found %d results below the threshold (%.2f)\n", skipped, minScore)
		}
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// searchCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// searchCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	searchCmd.Flags().BoolP("lucky", "l", false, "Return just the url of one result (for IRC usage).")
}
