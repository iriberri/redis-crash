#Redis Crash Tester for Heroku Redis
This small app is to demonstrate that Redis provided by Heroku can crash when multiple tasks are run against it. We started encountering the problem when we used Heroku Redis as a backend for Job-Queueing system and thus this sample program emulates that same scenario.

The repo contains two binaries that can be created. The first program is the _**worker**_ which looks for new jobs in the queue (which is held in the Redis server). The second program is an _**enqueuer**_ program which puts a specified number of jobs into the redis system at a very fast rate (using go routines).

# Installation
This app was compiled and tested against go 1.15.7 but should produce the same output with go 1.14.x and 1.13.x as well.

It uses the **go modules** feature.

1. Clone the repo using `git clone`. **_We are going to assume that the repo is being cloned into `~/go/src/redis-crash`_**
2. Launch terminal and `cd` to the cloned directory (`cd ~/go/src/redis-crash`)
3. Run the following command to compile and run the worker: `go build -o tmp/worker ./cmd/worker && REDIS_BG_URL=redis://127.0.0.1 tmp/worker`. This will compile the worker and run it. The worker will start looking for jobs in the redis installation pointed to by `REDIS_BG_URL`
4. Run the following command to compile and run the enqueuer: ` go build -o tmp/testJobEnqueuer ./cmd/testJobEnqueuer && JOB_COUNT=1000 tmp/testJobEnqueuer`. This will compile the enqueuer and launch it to enqueue 1000 jobs (governed by `JOB_COUNT` variable) in the redis instance pointed to by the `REDIS_BG_URL`.

Change the `REDIS_BG_URL` value to test different Redis installations.

# How to run the programs
**NOTE**: Please make sure that both the worker and the enqueuer are using the same Redis installation.

You should run one worker (to consume) and 5+ enqueuers at the same time (don't just increase the job count, increase the number of **simultaneous enqueuers**).

**NOTE 2**: Either that or you can increase the number of connections in the pool near the max connections limit and run it.

In fact, the worker is optional; the only purpose it would serve for this demo is to drain the jobs so that Redis server does not fill up. We are not encountering any issues with the worker. It's the enqueueing that fails.

The error messages that we see (mostly on the 3rd, 4th, 5th terminal run of the enqueuer) are like these:

```
E#7X6OQ - Could not enqueue: EOF
E#7X6OQ - Could not enqueue: read tcp 192.168.29.50:56565->52.22.235.152:9390: read: connection reset by peer
```

In your own run, the source IP address (`192.168.29.50`) should be something different. 

This never happened with any other Redis installation we tested. SSL/TLS connection does not seem to be the issue, neither does it look like it is because of the client connections count (it never crossed 16 from what I can see in the Heroku Dashboard).

# Metrics

Heroku dashboard never shows us hitting the connection limit though!
