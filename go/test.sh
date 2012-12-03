#!/bin/bash
cd ..
make clean > /dev/null
make > /dev/null
cd - > /dev/null

# Set variables
# Pick random line num between [500, 10000)
LINES=$(((RANDOM % 500) + 500))
INPUT="wc"

rm .m*
rm -rf ${INPUT}
mkdir ${INPUT}
mkdir ${INPUT}/wc2
for i in `seq 0 $((LINES/3))`
do
	echo "cat dog rat" >> ${INPUT}/wc1.txt
done

for i in `seq $(((LINES/3) + 1)) $(((2*LINES)/3))`
do
	echo "cat dog rat" >> ${INPUT}/wc2.txt
done

for i in `seq $((((2*LINES)/3) + 1)) $((LINES-1))`
do
	echo "cat dog rat" >> ${INPUT}/wc2/wc3.txt
done

OUTPUT="log"
SERVER=./bin/server
MAIN=./bin/wc
REQUEST=./bin/request
EVIL=./bin/evil-worker
DEAD=./bin/dead-worker
LONG=./bin/long-worker
LOOP=./bin/loop-worker
WORKER=./bin/worker
REQS=0

# Pick random port between [10000, 20000)
PORT=$(((RANDOM % 10000) + 10000))

function startServer {
	${SERVER} ${PORT} 2> /dev/null &
	SERVER_PID=$!
}

function stopServer {
	kill -9 ${SERVER_PID}
    wait ${SERVER_PID} 2> /dev/null
}

function startEvilWorkers {
	N=$1
    for i in `seq 0 $((N-1))`
    do
        ${EVIL} localhost:${PORT} 2> /dev/null &
    done
}

function startDeadWorkers {
	N=$1
    for i in `seq 0 $((N-1))`
    do
        ${DEAD} localhost:${PORT} 2> /dev/null &
    done
}

function startLongWorker {
    ${LONG} localhost:${PORT} 2> /dev/null &
	LONG_PID=$!
}

function waitLongWorker {
    wait ${LONG_PID} 2> /dev/null
}

function startLoopWorker {
    ${LOOP} localhost:${PORT} 2> /dev/null &
	LOOP_PID=$!
}

function stopLoopWorker {
	kill -9 ${LOOP_PID}
    wait ${LOOP_PID} 2> /dev/null
}

function startWorkers {
	N=$1
    for i in `seq 0 $((N-1))`
    do
        ${WORKER} localhost:${PORT} $2 2> /dev/null &
    done
}

function startRequests {
	N=$((REQS + $1))
    for i in `seq ${REQS} $((N-1))`
    do
        ${REQUEST} localhost:${PORT} ${INPUT} ${OUTPUT}:${i} ${MAIN} > /dev/null &
        REQUEST_PID[$i]=$!
    done
}

function testResults {
	N=$1
	PASSED=0
	for i in `seq 0 $((N-1))`
    do
    	wait ${REQUEST_PID[$i]} 2> /dev/null
        PASS=`cat ${OUTPUT}:$i | grep ${LINES} | wc -l`
	    if [ "$PASS" -eq 3 ]
	    then
	    	PASSED=$((PASSED + 1))
	   	fi
	    rm -rf ${OUTPUT}:$i
    done
    if [ "$PASSED" -eq $N ]
    then
    	echo "PASS"
        PASS_COUNT=$((PASS_COUNT + 1))
    else
        echo "FAIL"
        FAIL_COUNT=$((FAIL_COUNT + 1))
    fi
}

function testOneWorkerOneRequest {
	echo "Testing: 1 worker then 1 request"
	startServer
	startWorkers 1 0
	startRequests 1
	testResults 1
	stopServer	
}

function testThreeWorkerOneRequest {
	echo "Testing: 3 workers then 1 request"
	startServer
	startWorkers 3 0
	startRequests 1
	testResults 1
	stopServer
}

function testOneRequestOneWorker {
	echo "Testing: 1 request then 1 worker"
	startServer
	startRequests 1
	startWorkers 1 0
	testResults 1
	stopServer
}

function testOneRequestThreeWorker {
	echo "Testing: 1 request then 3 workers"
	startServer
	startRequests 1
	startWorkers 3 0
	testResults 1
	stopServer
}

function testThreeWorkerThreeRequest {
	echo "Testing: 3 workers then 3 requests"
	startServer
	startWorkers 3 0
	startRequests 3
	testResults 3
	stopServer
}

function testThreeEqualWorkers {
	echo "Testing: 3 workers, same connection time"
	startServer
	startWorkers 3 0
	sleep 3
	startRequests 1
	testResults 1
	stopServer
}

function testThreeIncreasedWorkers {
	echo "Testing: 3 workers, increasing connection time"
	startServer
	startWorkers 1 0
	sleep 1
	startWorkers 1 0
	sleep 1
	startWorkers 1 0
	sleep 1
	startRequests 1
	testResults 1
	stopServer
}

function testTenBadWorkers {
	echo "Testing: 10 workers, failure rate of 25%"
	startServer
	startWorkers 10 25
	sleep .1
	startRequests 1
	testResults 1
	stopServer
}

function testHundredTerribleWorkers {
	echo "Testing: 100 workers, failure rate of 75%"
	startServer
	startWorkers 100 75
	startRequests 1
	testResults 1
	stopServer
}

function testMixedWorkers {
	echo "Testing: 6 requests, 60 workers, various failure rates and connection times"
	startServer
	startWorkers 10 0
	sleep 1
	startWorkers 10 25
	sleep 1
	startRequests 3
	startWorkers 10 45
	sleep 1
	startWorkers 10 65
	sleep 1
	startWorkers 10 85
	startEvilWorkers 10
	REQS=3
	startRequests 3
	testResults 6
	startWorkers 
	stopServer
	REQS=0
}

function testEvilWorkers {
	echo "Testing: 10 evil-workers, request, 10 evil workers, good worker"
	startServer
	startEvilWorkers 10
	sleep 3
	startRequests 1
	startEvilWorkers 10
	startWorkers 1 0
	testResults 1
	stopServer
}

function testTimeoutLongWorker {
	echo "Testing: worker timeout but returns"
	startServer
	startRequests 1
	startLongWorker
	waitLongWorker
	startWorkers 10 0
	testResults 1
	stopServer
}

function testTimeoutInfWorker {
	echo "Testing: worker timeout no return"
	startServer
	startLoopWorker
	startRequests 1
	sleep 10
	startWorkers 10 0
	testResults 1
	stopLoopWorker
	stopServer
}

# Run tests
PASS_COUNT=0
FAIL_COUNT=0
echo "Running tests with input file of length ${LINES}"
echo ""
echo "Running sanity tests"
testOneWorkerOneRequest
testThreeWorkerOneRequest
testOneRequestOneWorker
testOneRequestThreeWorker
testThreeWorkerThreeRequest
echo ""
echo "Running reliable worker tests"
testThreeEqualWorkers
testThreeIncreasedWorkers
testTenBadWorkers
testHundredTerribleWorkers
testMixedWorkers
testEvilWorkers
testTimeoutLongWorker
testTimeoutInfWorker

rm -rf ${INPUT}
echo "Passed (${PASS_COUNT}/$((PASS_COUNT + FAIL_COUNT))) tests"