package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"time"
)

type processedTimeZoneData struct {
	timezone     string
	variationPST float64
	variationEST float64
	comparison   string
}

func main() {

	pstTZ, _ := time.LoadLocation("US/Pacific")
	pstTimezone := time.Date(time.Now().Year(), time.Now().Month(), 0, 0, 0, 0, 0, pstTZ)
	fmt.Printf("PST Time: %s\n", pstTimezone)

	estTZ, _ := time.LoadLocation("US/Eastern")
	estTimezone := time.Date(time.Now().Year(), time.Now().Month(), 0, 0, 0, 0, 0, estTZ)
	fmt.Printf("EST Time: %s\n", estTimezone)

	f, err := os.Open("./data/timedata.csv")
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	csv_processed_data := []processedTimeZoneData{}

	csvReader := csv.NewReader(f)
	for {
		rec, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		csv_processed_data = compareTimezones(rec, pstTimezone, estTimezone, csv_processed_data)
	}

	writeProcessedData(csv_processed_data)
}

func writeProcessedData(csv_processed_data []processedTimeZoneData) {
	file, err := os.Create("./processed_data/timezone_data.csv")
	if err != nil {
		log.Fatalln("failed to open file", err)
		panic(err)
	}
	defer file.Close()

	w := csv.NewWriter(file)
	defer w.Flush()

	csvHeader := []string{"Timezone", "Variation to PST", "Variation to EST", "Comparison"}
	var csvData [][]string
	csvData = append(csvData, csvHeader)

	for _, record := range csv_processed_data {
		csvData = append(csvData, []string{record.timezone, fmt.Sprint(record.variationPST), fmt.Sprint(record.variationEST), record.comparison})
	}
	w.WriteAll(csvData)
}

func compareTimezones(rec []string, pstTimezone time.Time, estTimezone time.Time, csv_processed_data []processedTimeZoneData) []processedTimeZoneData {
	tzToCheck, _ := time.LoadLocation(rec[0])
	tzCompareLocalized := time.Date(time.Now().Year(), time.Now().Month(), 0, 0, 0, 0, 0, tzToCheck)
	fmt.Printf("Compared Time: %s\n", tzCompareLocalized)

	comparePST := math.Abs(float64((pstTimezone.Sub(tzCompareLocalized).Minutes())))
	compareEST := math.Abs(float64(estTimezone.Sub(tzCompareLocalized).Minutes()))

	fmt.Printf("PST: %v\n", comparePST)
	fmt.Printf("EST: %v\n", compareEST)

	comparison := ""
	if comparePST > compareEST {
		comparison = "Closer to EST"
	} else if comparePST < compareEST {
		comparison = "Closer to PST"
	} else if comparePST == compareEST {
		comparison = "No difference can be used either in EST or PST"
	}

	csv_processed_data = append(csv_processed_data, processedTimeZoneData{
		rec[0], comparePST, compareEST, comparison})
	return csv_processed_data
}
