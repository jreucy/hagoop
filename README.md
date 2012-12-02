HaGoop
======

A mapreduce framework built in Go with an emphasis on high churn

#### How to Run

    $ make
    $ ./server <port>
    $ ./worker localhost:<port>
    $ ./request localhost:<port> <src> <output> ./main

#### Phase 1

Assumptions:  
1 worker client  
1 request client  
1 server  
