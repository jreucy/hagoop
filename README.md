# HaGoop : A mapreduce framework

HaGoop is a 2 week project to create a mapreduce framework in Go with an emphasis on high churn.  Specifically, it is a mapreduce framework with an emphasis on high churn in the worker nodes.  We are anticipating our worker nodes to have a large range of failure rates.  Rather than expecting multiple instances of a single type of machine (and thus a single failure rate), we expect a range of machines to serve as our workers - some more reliable than others.  We are also expecting both worker clients and request clients to be entering and exiting our network at all times.  Thus, we need our framework to be robust to even higher failure rates than considered “normal”.  

One possible application would be a network that utilizes home computers, which have very different hardware specifications and constantly log in and out of the application.  This contrasts against the typical model of mapreduce, which is run against a datacenter of machines of the same type and a constant rate of failure is expected.

# Run Demo

Type into command line in the go directory:

	$ ./test.sh

This will run our tests, based on a simple word count implementation of up to 500k words, with nodes that are good, bad (shuts down right away), or evil (hangs indefinitely), with different chances of "crashing".

# Structure

### Assumptions

Since we only had two weeks to implement HaGoop, we emphasized the more interesting problem of scheduling while leaving certain other implementation features out.

- Distributed File System: We run our mapreduce framework on a single computer (as different processes).  This allows us to avoid implementing a distributed file system, which is another problem in it of itself.  Of course, our implementation still runs as if the processes were completely separate, to simulate the environment of running HaGoop on multiple machines.
- Replicated Server Node: We assume that the server node is super-reliable (a.k.a. cannot fail).
- File Splitting:  Our mappers take line numbers and a file as input.  Thus, we assume that clients will not want any keys that split a single line in any way.

### Components

We use TCP as our method of communication between different components.  We implemented 3 different components:

- Request: A mapreduce request TCP client. This is an implementation of our framework. The user defines and compiles map and reduce methods consistent with an interface given with the framework. The TCP client then sends a request to the server specifying the location of the compiled binary and the directory on which to perform the map-reduce operations.  It also receives the server’s response (the solution to their mapreduce request), prints the response, and exits.

- Worker: The worker continually accepts map and reduce requests from the server, executes map and reduce binaries with parameters supplied from the server, and responds to the server with the result of its job completion.

- Server:  The server keeps track of all the workers and requests.  It handles mapreduce requests, splits the jobs based on the reliabilities of its current workers, and farms out these jobs to available workers.  It returns the final result back to the client that send the request.  

# Job Scheduling Algorithm

We are using an “experience” based method of scheduling.  In short, the server treats new nodes as having the least experience, and thus gives them smaller jobs.  Nodes that have been connected to the server for longer periods of time and have completed smaller tasks have more experience and are given bigger jobs.  We’re using this algorithm because we are unsure about the reliability of each node.  Each node has a chance to “prove” itself to be reliable, after which the server will give them bigger jobs.  Otherwise, they could be very unreliable; hence, the smaller jobs.

We also have an upper bound on our job sizes to prevent a worker from receiving a job too large to complete efficiently.  Likewise, we have a lower bound on our job size to make sure the work to process small jobs isn’t trumped by network communications.

## Authors
* [Russell Cullen](https://github.com/coolbrow)
* [Jonathan Hsu](https://github.com/jreucy)