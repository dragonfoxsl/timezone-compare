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

// Set the struct for the processed data
type processedTimeZoneData struct {
	timezone     string
	variationPST float64
	variationEST float64
	comparison   string
}

func main() {

	// Get the primary timezones
	pstTZ, _ := time.LoadLocation("US/Pacific")
	pstTimezone := time.Date(time.Now().Year(), time.Now().Month(), 0, 0, 0, 0, 0, pstTZ)
	fmt.Printf("PST Time: %s\n", pstTimezone)

	estTZ, _ := time.LoadLocation("US/Eastern")
	estTimezone := time.Date(time.Now().Year(), time.Now().Month(), 0, 0, 0, 0, 0, estTZ)
	fmt.Printf("EST Time: %s\n", estTimezone)

	// Read the source data CSV file
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

		// Compare the timezones for each row in the CSV file
		// Add the processed data to the array
		csv_processed_data = compareTimezones(rec, pstTimezone, estTimezone, csv_processed_data)
	}

	// Call the function to write the processed data to a new CSV file
	writeProcessedData(csv_processed_data)
}

// Write processed data to a new CSV file
func writeProcessedData(csv_processed_data []processedTimeZoneData) {
	file, err := os.Create("./processed_data/timezone_data.csv")
	if err != nil {
		log.Fatalln("failed to open file", err)
		panic(err)
	}
	defer file.Close()

	csvFileWriter := csv.NewWriter(file)
	defer csvFileWriter.Flush()
	// # Write the header row
	csvHeader := []string{"Timezone", "Variation to PST", "Variation to EST", "Comparison"}
	var csvData [][]string
	csvData = append(csvData, csvHeader)

	for _, record := range csv_processed_data {
		csvData = append(csvData, []string{record.timezone, fmt.Sprint(record.variationPST), fmt.Sprint(record.variationEST), record.comparison})
	}
	csvFileWriter.WriteAll(csvData)
}

// Compare timezone data
func compareTimezones(rec []string, pstTimezone time.Time, estTimezone time.Time, csv_processed_data []processedTimeZoneData) []processedTimeZoneData {
	// Load the timezone to check
	tzToCheck, _ := time.LoadLocation(rec[0])
	tzCompareLocalized := time.Date(time.Now().Year(), time.Now().Month(), 0, 0, 0, 0, 0, tzToCheck)
	fmt.Printf("Compared Timezone Time: %s\n", tzCompareLocalized)

	// Compare against PST
	comparePST := math.Abs(float64((pstTimezone.Sub(tzCompareLocalized).Minutes())))
	fmt.Printf("Variation to PST Time: %v\n", comparePST)

	// Compare against EST
	compareEST := math.Abs(float64(estTimezone.Sub(tzCompareLocalized).Minutes()))
	fmt.Printf("Variation to EST Time: %v\n", compareEST)

	// Compare the PST vs EST
	comparison := ""
	if comparePST > compareEST {
		comparison = "Closer to EST"
	} else if comparePST < compareEST {
		comparison = "Closer to PST"
	} else if comparePST == compareEST {
		comparison = "No difference can be used either in EST or PST"
	}

	// Add to the slice created by the struct
	csv_processed_data = append(csv_processed_data, processedTimeZoneData{
		rec[0], comparePST, compareEST, comparison})
	return csv_processed_data
}
