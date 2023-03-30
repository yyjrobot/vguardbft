package main

import (
	"encoding/json"
	"math/rand"
	"time"
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
	return rand.Float64()*(90-(-90)) + (-90)
}

func randomLongitude() float64 {
	return rand.Float64()*(180-(-180)) + (-180)
}


func mockGpsData(length int) []byte {
	b := make([]byte, length)

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
			q <- &Proposal{
				Timestamp:   time.Now().UnixMicro(),
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
