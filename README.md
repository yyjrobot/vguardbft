<h1 align="center"> V-Guard: An Efficient Permissioned Blockchain for Achieving Consensus under Dynamic Memberships in V2X Networks </h1>


## About V-Guard

V-Guard achieves high performance operating under dynamically changing memberships, targeting the problem of vehicles' arbitrary connectivity on the roads. When dealing with membership changes, traditional BFT algorithms (e.g., PBFT and HotStuff) must stop any ongoing consensus process to update their system configuration using additional membership management approaches, thereby suffering from severe performance degradation.

In contrast, V-Guard integrates the consensus of membership management into the consensus of data transactions. Compared to traditional BFT algorithms, V-Guard not only achieves consensus for data transactions but also for their residing membership configurations.

    # Consensus target of tradition BFT algorithms: <data transactions>
    # Consensus target of V-Guard: <data transactions, membership profiles>

The integration makes V-Guard achieve consensus seamlessly under changing members (e.g., joining or leaving vehicles) and produces an immutable ledger recording traceable data with corresponding membership profiles.

### Features and Use Cases
#### Membership Management Unit (MMU)
V-Guard uses Membership Management Unit (MMU) to keep track of available vehicle connections and manage them into sets of membership profiles. The MMU describes a membership profile that contains a set of vehicles as a **booth**. The management of booths with a size of four is illustrated below.

![](./docs/booths.gif)

### A use case example
In below example, the yellow car is the proposer, and the MMU is managing 10 members. The consensus has two targets: data (in blue font) and its membership (in red font), where O is the orderingID (sequence #), B is the data batch, V is the booth, Q is the quorum, and Sigma is the threshold signature.

![](./docs/mmu-ordering.png)

Ordering instances can take place in different booths. E.g., when the proposer conducts consensus for data entry 1, the ordering takes place in Booth 1. When Booth 1 is still available, the next ordering instance can reuse the booth (data entry 2). However, when Booth 1 is no longer available (e.g., some vehicles go offline), the MMU will provide a new booth (i.e., Booth 3) for a new ordering instance.

![](./docs/mmu-consensus.png)

Consensus instances are executed periodically. They operate as "shuttle buses" with a sole purpose of committing the entries appended on the total order log. A consensus instance can operate in a different booth from ordering instances, where the new members will scrutinize the signatures of each entry if they have not participated in the corresponding ordering instance.

## Try the Current Version

### Install dependencies
GoLang should have been properly installed with `GOPATH` and `GOROOT`. The GoLang version should be at least `go1.17.6`. In addition, three external packages were used (check out `go.mod`).

    # threshold signatures
    go get go.dedis.ch/kyber
    # logging
    go get github.com/sirupsen/logrus
    # some math packages
    go get gonum.org/v1/gonum/

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

    # run validators
    ./scripts/run.sh 1 1 // this starts a validator whose ID=1
    ./scripts/run.sh 2 1 // this starts a validator whose ID=2
    ./scripts/run.sh 3 1
    ./scripts/run.sh 4 1
    ./scripts/run.sh 5 1


Check out `parameters.go` for further parameters tuning.

## Deployment on clusters
The project is under a double-blind review process. We temporarily redacted the deployment details to preserve anonymity.
