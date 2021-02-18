package jobs

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/gocraft/work"
)

type TestJob struct {
	RedisEnqueuer *work.Enqueuer
}

type TestJobParams struct {
	ContentString  string
	TimeOfCreation int64
}

const (
	// Name of the job
	Test = "TestJob"
	// Work Struct Key into which the JSON data will be put
	TestJobParamsKey = "testJobParams"
)

// Enqueue adds a job to the worker queue.
func (d *TestJob) Enqueue(params *TestJobParams) error {
	args, err := json.Marshal(params)
	if err != nil {
		log.Println("E#681J3 - Error when marshalling:", err)
		return err
	}

	strArgs := string(args)

	if params.TimeOfCreation%2 == 0 {
		// Even timestamp ones are run immediately
		_, err = d.RedisEnqueuer.Enqueue(Test, work.Q{TestJobParamsKey: strArgs})
	} else {
		// Odd timestamp ones are scheduled for one minute later
		_, err = d.RedisEnqueuer.EnqueueIn(Test, 60, work.Q{TestJobParamsKey: strArgs})
	}

	if err != nil {
		log.Println("E#681JX - Error when enqueuing Test job:", err)
		return err
	}

	log.Println("L#681K9 - Enqueued the job:", strArgs)

	return nil
}

func (d *TestJob) Perform(j *work.Job) error {
	var params *TestJobParams
	var performErrString string

	strParams, ok := j.Args[TestJobParamsKey].(string)
	if !ok {
		performErrString = "E#681L0 - Could not convert data to string"
		log.Println(performErrString)
		return errors.New(performErrString)
	}

	if err := json.Unmarshal([]byte(strParams), &params); err != nil {
		performErrString = fmt.Sprintf("E#681L4 - Could not unmarshall data: %v", err)
		log.Println(performErrString)
		return errors.New(performErrString)
	}

	fmt.Printf("L#683NU - Test Job Content: '%v', Time of Creation (UTC Unix Timestamp): %v\n", params.ContentString, params.TimeOfCreation)

	return nil
}
