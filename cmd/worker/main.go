package main

import (
	"fmt"
	"os"
	"os/signal"
	"redis-crash/cmd/worker/jobs"
	"syscall"
	"time"

	"github.com/gocraft/work"
	"github.com/gomodule/redigo/redis"
)

type envVars struct {
	redisUrl             string
	redisWorkerNamespace string
}

func getEnvVars() envVars {
	vars := envVars{
		redisUrl:             getEnvValueOrDefault("REDIS_BG_URL", "redis://127.0.0.1:6379"),
		redisWorkerNamespace: getEnvValueOrDefault("REDIS_BG_WORKER_NS", "bgWorkerNamespace"),
	}
	return vars
}

func getEnvValueOrDefault(envVarName string, defaultValue string) string {
	returnValue := os.Getenv(envVarName)
	if returnValue == "" {
		returnValue = defaultValue
	}
	return returnValue
}

// Global variables for the worker
var (
	redisPool     *redis.Pool
	redisEnqueuer *work.Enqueuer
	vars          envVars
	//NOTE: We have to use a struct and we don't have one, so we are using an empty struct
	redisWorkerContext struct{}
)

func init() {
	vars = getEnvVars()

	// Init a redis pool
	redisPool = &redis.Pool{
		MaxActive: 5,
		MaxIdle:   5,
		Wait:      true,
		Dial: func() (redis.Conn, error) {
			return redis.DialURL(vars.redisUrl, redis.DialTLSSkipVerify(true))
		},
	}

	// We need a context with type struct. So we set a blank one
	redisWorkerContext = struct{}{}

	// Make an enqueuer with the namespace
	redisEnqueuer = work.NewEnqueuer(vars.redisWorkerNamespace, redisPool)
}

func main() {
	//vars := getEnvVars()
	fmt.Println("L#84USJ - REDIS_BG_URL is      :", vars.redisUrl)
	fmt.Println("L#84USL - REDIS_BG_WORKER_NS is:", vars.redisWorkerNamespace)

	testJob := &jobs.TestJob{
		RedisEnqueuer: redisEnqueuer,
	}

	// ----------- REDIS WORKER POOL AND JOBS -----------
	redisWorkerPool := work.NewWorkerPool(redisWorkerContext, 10, vars.redisWorkerNamespace, redisPool)

	jobOptions := work.JobOptions{
		Priority:       100,
		MaxFails:       10,
		SkipDead:       false,
		MaxConcurrency: 0,
		Backoff:        nil,
	}

	// Test Job, for testing if the worker is working or not
	// (usually for use on the local development machine)
	redisWorkerPool.JobWithOptions(jobs.Test, jobOptions, testJob.Perform)

	fmt.Printf("L#84UT3 - Starting Redis based worker at %v\n", time.Now())
	redisWorkerPool.Start()

	// ----------- Wait for interrupt -----------
	// Wait for a signal to quit:
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan

	// ----------- Quit on interrupt -----------
	fmt.Println("\nL#84UTN - Worker shutting down...")

	// Stop the Redis based redisWorkerPool
	fmt.Println("L#84UTT - Stopping Redis based redisWorkerPool")
	redisWorkerPool.Stop()
	fmt.Println("L#84UTX - Redis based redisWorkerPool stopped")

	fmt.Println("L#84UU1 - Worker shutdown complete")
}
