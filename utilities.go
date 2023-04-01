package main

import (
	"encoding/json"
	"math/rand"
	"time"
	//"fmt"
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

// TODO: Change here to add some vehicle driving data
// txGenerator enqueues mock data entries to all message queues
func txGenerator(len int) {
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"



	for i := 0; i < NumOfValidators; i++ {
		q := make(chan *Proposal, MaxQueue)

		for i := int64(0); i < MsgLoad; i++ {
			rand_lat := randomLatitude()
			rand_lon := randomLongitude()
			log.Infof("lat: %v, lon: %v", rand_lat, rand_lon)
			q <- &Proposal{
				Timestamp:   time.Now().UnixMicro(),
				lat: rand_lat,
				lon: rand_lon,
				speed: rand.Int63(),
				// lat: rand_lat,
				// lon: rand_lon,
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
