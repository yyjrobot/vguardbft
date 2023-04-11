<h1 align="center"> V-Guard with Database Integration (MongoDB) </h1>

This project is build based on [V-Guard](https://github.com/vguardbc/vguardbft) by Edward (Gengrui) Zhang.

## About V-Guard with Database Integration (MongoDB) 

V-Guard is a recent paper that proposed a design for a new permissioned blockchain system dedicated to achieve consensus for vehicular data under changing memberships. The goal of V-Guard is grant full vehicular data access to consumers to avoid the centralized data monopoly of manufacturers. However, current V-Guard design only supports storing the consensus log locally, which could become an issue if the vehicle get involved in an accident and caused the storage component to be broken. In such case a remote backup of the data could be vital for law purposes. Hence this project aims to extend current V-Guard implementation with database to store the consensus log in addition to directly storing on vehicles to avoid the issue that could caused by data loss during a vehicle accident. We will also evaluate the performance of connecting V-Guard to a database comparing to V-Guard BFT without any modification.
This project extends the above V-Guard consensus protocol, with an integration to a cloud based database (MongoDB). 

### Solution
- Vehicular Data
After the consensus process of vanilla V-Guard has been completed, the Vehicular Data is sent to MongoDB to be stored with a timestamp. The dataset that is being pushed to MongoDB are [Passive Vehicular Sensors Dataset (PVS)](https://www.kaggle.com/datasets/jefmenegazzo/pvs-passive-vehicular-sensors-datasets?resource=download).

## Real Life Applications
### Determine Vehicle Collision Responsiblities
- By providing a backup of dataset that is not stored locally, this data can be used later as evidence for any legal proceedings.

### Providing Transparency
- The dataset provides transparency to all vehicles on the road, on the make, model, and vehicular data.

### Prevention of Vehicle Theft
- By extending V-Guard with a cloud databse, it allows better tracking of stolen vehicles to be tracked down and retrieved.


## Instructions (Based on V-Guard instructions)
### Install dependencies
GoLang should have been properly installed with `GOPATH` and `GOROOT`. The GoLang version should be at least `go1.17.6`. In addition, three external packages were used (check out `go.mod`).

Install the latest version of docker container.

    // threshold signatures
    go get go.dedis.ch/kyber
    // logging
    go get github.com/sirupsen/logrus
    // some math packages
    go get gonum.org/v1/gonum/
    // mongodb
    go get go.mongodb.org/mongo-driver

### Run MongoDB instances locally
Below shows an example of running a mongodb database locally. 

    $ docker pull mongo
    $ docker run -d -p 27017:27017 --name mongo mongo:latest

### Run V-Guard instances locally
Below shows an example of running a V-Guard instance with a booth of size 4 and 6 initial available connections. The quorum size of a booths of size 4 is 3, so the threshold is set to 2, as the proposer is always included.
    
    // Assume the downloaded folder is called "vguardbft"
    // First, move the keyGen folder outside of vguardbft.
    mv keyGen ../
    
    // Then, go to "keyGen" and generate keys
    cd ../keyGen
    go build generator.go
    
    // Keys are private and public keys for producing and 
    // validatoring threshold signatures where t is the threshold
    // and n is the number of participants
    ./generator -t=2 -n=6
    
    // A "keys" folder should be generated with 6 private keys and 1 public key
    // Privates keys: pri_#id.dupe
    // Public key: vguard_pub.dupe
    // Now copy the "keys" folder into the "vguardbft" folder
    cp -r keys ../vguardbft/
    
    // Compile the code in "vguardbft" using the build script
    cd ../vguardbft
    ./scripts/build.sh

    // Next, create a log folder
    mkdir logs
    
    // Finally, we can start running a V-Guard instance by starting
    // a proposer, which always has an ID of 0. The script takes 
    // two parameters: $1=ID; $2=role (proposer: 0; validator: 1)
    ./scripts/run.sh 0 0 // this starts a proposer

    // run 5 validators
    ./scripts/run.sh 1 1 // this starts a validator whose ID=1
    ./scripts/run.sh 2 1 // this starts a validator whose ID=2
    ./scripts/run.sh 3 1
    ./scripts/run.sh 4 1
    ./scripts/run.sh 5 1


Check out `parameters.go` for further parameters tuning.

After running V-Guard instances locally, you should be able to see the vehicle_data table in mongodb by running
    $ show dbs
      vehicle_data

The latency data can be revisited via `/logs` folder.
