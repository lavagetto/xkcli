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
		db, err := database.Open(dbPath)
		defer db.Close()
		if err != nil {
			panic(err)
		}
		searchResult, err := database.SearchStr(db, args[0], nil)
		if err != nil {
			panic(err)
		}
		fmt.Println("Your search results:")
		for pos, result := range searchResult.Hits {
			if result.Score > minScore {
				strip := database.NewStripFromDb(result)
				fmt.Printf("%d - (%.2f) %s", pos, result.Score, strip.Summary())
			}
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
}
