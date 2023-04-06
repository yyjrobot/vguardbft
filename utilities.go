package main

import (
	"encoding/json"
	"encoding/csv"
	"math/rand"
	"time"
	"os"
	"strconv"
)

var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))

func mockRandomBytes(length int, charset string) []byte {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return b
}

func randomLatitude() float64 {
	lat_num := rand.Float64()*(90-(-90)) + (-90)
	//lat := fmt.Sprintf("%.6f", lat_num)
	return lat_num
}

func randomLongitude() float64 {
	lon_num := rand.Float64()*(180-(-180)) + (-180)
	//lon := fmt.Sprintf("%.6f", lon_num)
	return lon_num
}


func mockGpsData(length int) []byte {
	b := make([]byte, length)
	//charset := "0"
	rand_lat := randomLatitude()
	rand_lon := randomLongitude()

	// Set the remaining bytes of b to zero
	// for i := 0; i < length; i++ {
	// 	b[i] = charset[0]
	// }

	// Copy rand_lat and rand_lon into b
	//length := len(rand_lat)
	copy(b, []byte(strconv.FormatFloat(rand_lat, 'f', 6, 64)))
	copy(b[16:], []byte(strconv.FormatFloat(rand_lon, 'f', 6, 64)))
	
	
	log.Infof("%v is generated nsg load", string(b))
	// for i := range b {
	// 	b[i] = 
	// }
	return b
}

// https://stackoverflow.com/questions/24999079/reading-csv-file-in-go
func readCsvFile(filePath string) [][]string {
    f, err := os.Open(filePath)
    if err != nil {
        log.Fatal("Unable to read input file " + filePath, err)
    }
    defer f.Close()

    csvReader := csv.NewReader(f)
    records, err := csvReader.ReadAll()
    if err != nil {
        log.Fatal("Unable to parse file as CSV for " + filePath, err)
    }

    return records
}

// TODO: Change here to add some vehicle driving data
// txGenerator enqueues mock data entries to all message queues
func txGenerator(len int) {
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	
	// Hardcode vehiculer data file here
	gps_data := readCsvFile("dataset_gps.csv")


	for i := 0; i < NumOfValidators; i++ {
		q := make(chan *Proposal, MaxQueue)

		for i := int64(0); i < MsgLoad; i++ {
			// rand_lat := randomLatitude()
			// rand_lon := randomLongitude()
			// log.Infof("lat: %v, lon: %v", rand_lat, rand_lon)
			
			// Take modulo here since we don't have that much data in our dataset
			// +1 to skip the first row
			latitude, _ := strconv.ParseFloat(gps_data[i % 1000 + 1][1], 64)
			longitude, _ := strconv.ParseFloat(gps_data[i % 1000 + 1][2], 64)
			speed_meters_per_second, _ := strconv.ParseFloat(gps_data[i % 1000 + 1][6], 64)

			latitudeStr := strconv.FormatFloat(latitude, 'f', 8, 64)
			longitudeStr := strconv.FormatFloat(longitude, 'f', 8, 64)
			speedStr := strconv.FormatFloat(speed_meters_per_second, 'f', 8, 64)

			
			q <- &Proposal{
				Timestamp:   time.Now().UnixMicro(),
				Lat: latitudeStr,
				Lon: longitudeStr,
				Speed: speedStr,
				Transaction: mockRandomBytes(len, charset),
				
			}
		}
		requestQueue = append(requestQueue, q)
	}

	log.Infof("%d request queue(s) loaded with %d requests of size %d bytes", NumOfValidators, MsgLoad, MsgSize)
}

func serialization(m interface{}) ([]byte, error) {
	return json.Marshal(m)
}

func deserialization(b []byte, m interface{}) error {
	return json.Unmarshal(b, m)
}
