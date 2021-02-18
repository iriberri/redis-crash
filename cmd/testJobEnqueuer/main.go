package main

import (
	"fmt"
	"os"
	"redis-crash/cmd/worker/jobs"
	"strconv"
	"sync"
	"time"

	"github.com/gocraft/work"
	"github.com/gomodule/redigo/redis"
)

var wg sync.WaitGroup

type envVars struct {
	redisUrl             string // Redis URL for the worker
	redisWorkerNamespace string // Worker namespace
	messageCount         int    // Number of messages to send
}

func getEnvVars() envVars {
	vars := envVars{
		redisUrl:             getEnvValueOrDefault("REDIS_BG_URL", "redis://127.0.0.1:6379"),
		redisWorkerNamespace: getEnvValueOrDefault("REDIS_BG_WORKER_NS", "bgWorkerNamespace"),
	}

	messageCountStr := getEnvValueOrDefault("JOB_COUNT", "1000")
	intMessageCount, err := strconv.Atoi(messageCountStr)
	if err != nil {
		intMessageCount = 1000
	}

	vars.messageCount = intMessageCount

	return vars
}

func getEnvValueOrDefault(envVarName string, defaultValue string) string {
	returnValue := os.Getenv(envVarName)
	if returnValue == "" {
		returnValue = defaultValue
	}
	return returnValue
}

func main() {
	vars := getEnvVars()

	fmt.Println("L#7XIVM - REDIS_BG_URL is      :", vars.redisUrl)
	fmt.Println("L#84UPJ - REDIS_BG_WORKER_NS is:", vars.redisWorkerNamespace)
	fmt.Println("L#84UQA - JOB_COUNT is         :", vars.messageCount)

	redisPool := &redis.Pool{
		MaxActive: 5,
		MaxIdle:   5,
		Wait:      true,
		Dial: func() (redis.Conn, error) {
			return redis.DialURL(vars.redisUrl, redis.DialTLSSkipVerify(true))
		},
	}
	// Make an enqueuer with a particular namespace
	var enqueuer = work.NewEnqueuer(vars.redisWorkerNamespace, redisPool)

	testJob := &jobs.TestJob{
		RedisEnqueuer: enqueuer,
	}

	startTime := time.Now().UTC()

	fmt.Println("L#7X63X - Test Job Enqueuer Starts at", time.Now().UTC())
	wg.Add(vars.messageCount)

	for i := 0; i < vars.messageCount; i++ {
		idx := i
		go func() {
			time.Sleep(1 * time.Second)

			err := testJob.Enqueue(&jobs.TestJobParams{
				ContentString:  fmt.Sprintf("Index: %v RunID: %v", idx, strconv.FormatInt(startTime.UnixNano(), 36)),
				TimeOfCreation: time.Now().UTC().Unix(),
			})

			if err != nil {
				fmt.Println("E#7X6OQ - Could not enqueue:", err)
			}
			wg.Done()
		}()

		//time.Sleep(5 * time.Millisecond)
	}

	wg.Wait()
	fmt.Println("---------------------")
	fmt.Println("L#7X647 - Test Job Enqueuer Ends at", time.Now().UTC())
}
