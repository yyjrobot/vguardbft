<h1 align="center"> V-Guard: An Efficient Permissioned Blockchain for Achieving Consensus under Dynamic Memberships in V2X Networks </h1>


## About V-Guard with Database Integration (MongoDB Blockchain DB) 

V-Guard achieves high performance operating under dynamically changing memberships, targeting the problem of vehicles' arbitrary connectivity on the roads. When membership changes occur, traditional BFT algorithms (e.g., PBFT and HotStuff) must stop to update system configurations using additional membership management approaches, thereby suffering from severe performance degradation.


The current V-Guard design only is limited to storing the consensus log locally. This project extends the above V-Guard consensus protocol, with an integration to a cloud based database (MongoDB Blockchain DB).

## Use Case
V-Guard is a flexible blockchain platform that allows users to define their own message types. This platform enables vehicles to reach a consensus on the decisions made by their autonomous driving software. The messages can include various data, such as GPS location, speed, direction, acceleration, bearing, and more (similar categories to the data set of
[Passive Vehicular Sensors](https://www.kaggle.com/datasets/jefmenegazzo/pvs-passive-vehicular-sensors-datasets?resource=download-directory)).


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

## Try the Current Version

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

## Original V-Guard Project
| Paper         |   Authors                  |                                         
|:-------------------------------------------------------|:----------------------------|
| [VGuardDB](https://github.com/vguardbc/vguardbft) | Edward (Gengrui) Zhang