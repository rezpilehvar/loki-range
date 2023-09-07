package lokirange

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/rezpilehvar/loki-range/utils"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

func Query(queryURL string, query string, limit int, timeRange string, start, end string) ([]LogItem, error) {

	if len(timeRange) > 0 {
		now := time.Now()
		switch {
		case timeRange == "today":
			{
				start = utils.BeginningOfDay(now).Format(time.RFC3339)
				end = now.Format(time.RFC3339)
			}
		case timeRange == "yesterday":
			{
				yesterday := now.AddDate(0, 0, -1)
				start = utils.BeginningOfDay(yesterday).Format(time.RFC3339)
				end = utils.EndOfDay(yesterday).Format(time.RFC3339)
			}
		case strings.HasSuffix(timeRange, "d"):
			{
				daysStr, _ := strings.CutSuffix(timeRange, "d")
				days, err := strconv.Atoi(daysStr)
				if err != nil {
					return nil, errors.New("invalid range format")
				}

				fromDate := now.AddDate(0, 0, -days)
				start = utils.BeginningOfDay(fromDate).Format(time.RFC3339)
				end = now.Format(time.RFC3339)
			}
		case strings.HasSuffix(timeRange, "h"):
			{
				hoursStr, _ := strings.CutSuffix(timeRange, "h")
				hours, err := strconv.Atoi(hoursStr)
				if err != nil {
					return nil, errors.New("invalid range format")
				}

				fromDate := now.Add(time.Duration(-hours) * time.Hour)
				start = fromDate.Format(time.RFC3339)
				end = now.Format(time.RFC3339)
			}
		case strings.HasSuffix(timeRange, "m"):
			{
				minutesStr, _ := strings.CutSuffix(timeRange, "m")
				minutes, err := strconv.Atoi(minutesStr)
				if err != nil {
					return nil, errors.New("invalid range format")
				}

				fromDate := now.Add(time.Duration(-minutes) * time.Minute)
				start = fromDate.Format(time.RFC3339)
				end = now.Format(time.RFC3339)
			}
		default:
			{
				return nil, errors.New("invalid range format")
			}
		}
	}

	fmt.Println(fmt.Sprintf("input start: %s", start))
	fmt.Println(fmt.Sprintf("input end: %s", end))

	chunk := 1
	fetchStart := start
	fetchEnd := end
	var collectedLogItems []LogItem
	for {
		fmt.Println(fmt.Sprintf("--------loading chunk #%d--------", chunk))
		res, err := fetchData(queryURL, query, fetchStart, fetchEnd, limit)
		if err != nil {
			return nil, err
		}

		fmt.Println(fmt.Sprintf("chunk #%d-> query exec time: %f", chunk, res.Data.Stats.Summary.ExecTime))
		fmt.Println(fmt.Sprintf("chunk #%d-> entries returned: %d", chunk, res.Data.Stats.Summary.TotalEntriesReturned))

		collectedLogItems = append(collectedLogItems, res.Data.Result...)
		if len(res.Data.Result) > 0 && res.Data.Stats.Summary.TotalEntriesReturned == limit {
			lastItemTimeNano, _ := strconv.ParseInt(res.Data.Result[len(res.Data.Result)-1].Values[0][0], 10, 64)
			fetchStart = time.Unix(0, lastItemTimeNano).Format(time.RFC3339)
			chunk++
		} else {
			break
		}
	}
	fmt.Println(fmt.Sprintf("data collected in %d chunks", chunk))
	fmt.Println(fmt.Sprintf("total retreived entries: %d", len(collectedLogItems)))

	return collectedLogItems, nil
}

func WriteToCsv(collectedLogItems []LogItem, reportName string) error {

	headers := make(map[string]int)
	for _, record := range collectedLogItems {
		for key := range record.Stream {
			if _, ok := headers[key]; !ok {
				headers[key] = len(headers)
			}
		}
	}
	headersArr := make([]string, len(headers))

	for key, value := range headers {
		headersArr[value] = key
	}

	rows := make([][]string, len(collectedLogItems))

	for idx, record := range collectedLogItems {
		rows[idx] = make([]string, len(headers))

		for key, value := range record.Stream {
			rows[idx][headers[key]] = value
		}
	}

	csvFile, err := os.Create(fmt.Sprintf("%s.csv", reportName))
	if err != nil {
		return fmt.Errorf("failed creating csv file: %s", err)
	}
	defer func(csvFile *os.File) {
		err := csvFile.Close()
		if err != nil {
			log.Fatal(fmt.Errorf("failed to close csv file: %s", err))
		}
	}(csvFile)

	csvWriter := csv.NewWriter(csvFile)
	err = csvWriter.Write(headersArr)
	if err != nil {
		return err
	}
	err = csvWriter.WriteAll(rows)
	if err != nil {
		return err
	}

	csvWriter.Flush()
	return nil
}

func fetchData(lokiQueryURL string, query string, start string, end string, limit int) (*LokiQueryResponse, error) {
	params := url.Values{}
	params.Add("limit", strconv.Itoa(limit))
	params.Add("start", start)
	params.Add("end", end)
	params.Add("query", query)
	params.Add("direction", "forward")

	baseURL, _ := url.Parse(lokiQueryURL)
	baseURL.RawQuery = params.Encode()

	resp, err := http.Get(baseURL.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		fmt.Println("OK, items are collected.")

		fmt.Println("unpacking data...")

		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}

		var response LokiQueryResponse
		err = json.Unmarshal(bodyBytes, &response)
		if err != nil {
			return nil, err
		}

		fmt.Println("sorting logs...")

		sort.Slice(response.Data.Result, func(i, j int) bool {
			iNanoTime, _ := strconv.ParseInt(response.Data.Result[i].Values[0][0], 10, 64)
			jNanoTime, _ := strconv.ParseInt(response.Data.Result[j].Values[0][0], 10, 64)
			return iNanoTime < jNanoTime
		})

		fmt.Println("sort done.")
		return &response, nil
	} else {
		bodyBytes, _ := io.ReadAll(resp.Body)

		return nil, errors.New(fmt.Sprintf("api error response, code: %d, body: %s", resp.StatusCode, string(bodyBytes)))
	}
}

type LokiQueryResponse struct {
	Status string `json:"status"`
	Data   struct {
		Result []LogItem `json:"result"`
		Stats  struct {
			Summary struct {
				ExecTime             float32 `json:"execTime"`
				TotalEntriesReturned int     `json:"totalEntriesReturned"`
			} `json:"summary"`
		} `json:"stats"`
	}
}

type LogItem struct {
	Stream map[string]string `json:"stream"`
	Values [][]string        `json:"values"`
}
