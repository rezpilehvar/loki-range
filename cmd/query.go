package cmd

import (
	"fmt"
	lokirange "github.com/rezpilehvar/loki-range/pkg/lokirange"
	"github.com/rezpilehvar/loki-range/utils"
	"github.com/spf13/cobra"
	"log"
)

var (
	lokiUrl   string
	limit     int
	start     string
	end       string
	timeRange string
	query     string
)

var queryCmd = &cobra.Command{
	Use:   "query",
	Short: "Run Loki queries with limit",
	Long:  `N/A`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		query = args[0]
		fmt.Println(fmt.Sprintf("query: %s", query))
		fmt.Println(fmt.Sprintf("limit: %d", limit))

		if len(timeRange) > 0 {
			cStart, cEnd, err := utils.CalculateTimeRange(timeRange)
			if err != nil {
				log.Fatal(err)
			}
			start = utils.Format(cStart)
			end = utils.Format(cEnd)
		}

		collectedLogs, queryType, err := lokirange.Query(lokiUrl, query, limit, start, end)
		if err != nil {
			log.Fatal(err)
		}

		err = lokirange.WriteToCsv(collectedLogs, "export", queryType)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(queryCmd)

	queryCmd.Flags().IntVar(
		&limit,
		"limit",
		5000,
		"limitation of entries",
	)

	queryCmd.Flags().StringVar(
		&lokiUrl,
		"url",
		"",
		"base url of loki gateway",
	)
	queryCmd.MarkFlagRequired("url")

	queryCmd.Flags().StringVar(
		&start,
		"start",
		"",
		"start time of entities with RFC3339 format",
	)
	queryCmd.Flags().StringVar(
		&end,
		"end",
		"",
		"end time of entities with RFC3339 format",
	)
	queryCmd.MarkFlagsRequiredTogether("start", "end")

	queryCmd.Flags().StringVar(
		&timeRange,
		"range",
		"",
		"supported formats: today, yesterday, {x}d, {x}h, {x}m",
	)
	queryCmd.MarkFlagsMutuallyExclusive("range", "start", "end")
	// TODO add oneFlagRequired after cobra release!
}
