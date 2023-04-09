package main

/*
The implementation of V-Guard follows a "cache more, lock less" policy. This design
reduces lock overhead and contention while storing more intermediate results.
Intermediate consensus information for data batches are stored separately in the
ordering and consensus phases.
*/

import (
	"context"
	"encoding/hex"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ordSnapshot stores consensus information for each block in the ordering phase
// the map stores <blockID, blockSnapshot>
var ordSnapshot = struct {
	m map[int64]*blockSnapshot
	sync.RWMutex
}{m: make(map[int64]*blockSnapshot)}

// cmtSnapshot stores consensus information for each block in the consensus phase
// the map stores <blockID, blockSnapshot>
var cmtSnapshot = struct {
	m map[int64]*blockSnapshot
	sync.RWMutex
}{m: make(map[int64]*blockSnapshot)}

type blockSnapshot struct {
	sync.RWMutex
	// The hash of the block
	hash []byte
	// The data entries
	entries map[int]Entry
	// The signatures collected from validators to be converted to a threshold signature
	sigs [][]byte
	// rcvSig is the threshold signature of this block
	tSig []byte
	// The booth of this block
	booth Booth
}

var vgTxMeta = struct {
	sync.RWMutex
	sigs     map[int][][]byte // <rangeId, sigs[]>
	hash     map[int][]byte
	blockIDs map[int][]int64 // <rangeId, []blockIDs>
}{
	sigs:     make(map[int][][]byte),
	hash:     make(map[int][]byte),
	blockIDs: make(map[int][]int64),
}

var vgTxData = struct {
	sync.RWMutex
	tx  map[int]map[string][][]Entry // map<consInstID, map<orderingBooth, []entry>>
	boo map[int]Booth                //<consInstID, Booth>
}{
	tx:  make(map[int]map[string][][]Entry),
	boo: make(map[int]Booth),
}

func storeChunkToDB(chunk []Entry) (ids []interface{}) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://127.0.0.1:27017"))
    if err != nil {
        log.Fatal(err)
    }
    collection := client.Database("Mongo").Collection("vehicle_data")
    var documents []interface{}
    for _, e := range chunk {
        documents = append(documents, e)
    }
    res, err := collection.InsertMany(ctx, documents)
    if err != nil {
        log.Fatal(err)
    }
    for _, id := range res.InsertedIDs {
        ids = append(ids, id)
    }
    return
}

func storeToDB(e Entry) (id interface{}) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://127.0.0.1:27017"))
	if err != nil {
		log.Fatal(err)
	}
	collection := client.Database("Mongo").Collection("vehicle_data")
	res, err := collection.InsertOne(ctx, e)
	id = res.InsertedID
	if err != nil {
		log.Fatal(err)
	}
	return
}




func storeVgTx(consInstID int) {
	vgTxData.RLock()
	ordBoo := vgTxData.tx[consInstID]  //ordering booth?
	cmtBoo := vgTxData.boo[consInstID] //commit booth
	vgTxData.RUnlock()

	log.Infof("VGTX %d in Cmt Booth: %v | total # of tx: %d", consInstID, cmtBoo.Indices, vgrec.GetLastIdx()*BatchSize)

	entryCount := int64(0)
	chunkCount := int64(0)
	mongoStoreTime := int64(0)

	for key, chunk := range ordBoo { //map<boo, [][]entries>
		log.Infof("ordering booth: %v | len(ordBoo[%v]): %v", key, key, len(chunk))
		
		for _, entries := range chunk {
			chunkCount++
			// Store entry to mongo
			startTime := time.Now().UnixMilli()
			storeChunkToDB(entries)
			endTime := time.Now().UnixMilli()
			mongoStoreTime += (endTime - startTime)
			for _, e := range entries {
				log.Infof("ts: %v; tx: %v, lat: %v, lon: %v, speed: %v", e.TimeStamp, hex.EncodeToString(e.Tx), e.Lat, e.Lon, e.Speed)
				//log.Infof("With out encoded: ts: %v; tx: %v", e.TimeStamp, e.Tx)
				entryCount++

				

			}
		}
		
	}
	log.Infof("Store finished, avg latency for each entry: %v", mongoStoreTime/entryCount)
	log.Infof("Avg latency for each entries in chunk: %v", mongoStoreTime/chunkCount)
}


