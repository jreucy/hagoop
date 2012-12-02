#!/bin/bash
cd ..
make > /dev/null
cd - > /dev/null

# Set variables
LINES=500
INPUT="wc.txt"

touch ${INPUT}
for i in `seq 0 $((LINES-1))`
do
	echo "cat" >> ${INPUT}
	echo "dog" >> ${INPUT}
	echo "rat" >> ${INPUT}
done

OUTPUT="log"
SERVER=./bin/server
MAIN=./bin/wc
REQUEST=./bin/request
EVIL=./bin/evil-worker
DEAD=./bin/dead-worker
LONG=./bin/long-worker
WORKER=./bin/worker

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
        EVIL_PID[$i]=$!
    done
}

function stopEvilWorkers {
	N=${#EVIL_PID[@]}
    for i in `seq 0 $((N-1))`
    do
        kill -9 ${EVIL_PID[$i]}
        wait ${EVIL_PID[$i]} 2> /dev/null
    done
}

function startDeadWorkers {
	N=$1
    for i in `seq 0 $((N-1))`
    do
        ${DEAD} localhost:${PORT} 2> /dev/null &
        DEAD_PID[$i]=$!
    done
}

function stopDeadWorkers {
	N=${#DEAD_PID[@]}
    for i in `seq 0 $((N-1))`
    do
        kill -9 ${DEAD_PID[$i]}
        wait ${DEAD_PID[$i]} 2> /dev/null
    done
}

function startWorkers {
	N=$1
    for i in `seq 0 $((N-1))`
    do
        ${WORKER} localhost:${PORT} 2> /dev/null &
        WORKER_PID[$i]=$!
    done
}

function stopWorkers {
	N=${#WORKER_PID[@]}
    for i in `seq 0 $((N-1))`
    do
        kill -9 ${WORKER_PID[$i]}
        wait ${WORKER_PID[$i]} 2> /dev/null
    done
}

function startRequest {
    ${REQUEST} localhost:${PORT} ${INPUT} ${OUTPUT} ${MAIN} > /dev/null
    PASS=`cat ${OUTPUT} | grep ${LINES} | wc -l`
    if [ "$PASS" -eq 3 ]
    then
        echo "PASS"
        PASS_COUNT=$((PASS_COUNT + 1))
    else
        echo "FAIL"
        FAIL_COUNT=$((FAIL_COUNT + 1))
    fi
    rm -rf ${OUTPUT}
}

function testOneWorkerOneRequest {
	echo "Testing: 1 worker then 1 request"
	sleep 1
	startServer
	startWorkers 1
	startRequest
	stopWorkers
	stopServer
}

function testThreeWorkerOneRequest {
	echo "Testing: 3 workers then 1 request"
	sleep 1
	startServer
	startWorkers 3
	startRequest
	stopWorkers
	stopServer
}

# Run tests
PASS_COUNT=0
FAIL_COUNT=0
testOneWorkerOneRequest
testThreeWorkerOneRequest

rm -rf ${INPUT}
echo "Passed (${PASS_COUNT}/$((PASS_COUNT + FAIL_COUNT))) tests"