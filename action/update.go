package action

import (
	"errors"
	"fmt"
	"time"

	log "github.com/Sirupsen/logrus"
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/apis/batch"
	client "k8s.io/kubernetes/pkg/client/unversioned"
)

const (
	LabelTugbotEvents      = "tugbot.kubernetes.events"
	LabelTugbotCreatedFrom = "tugbot.created.from"
)

func UpdateJobs(kube client.JobInterface, event string) error {
	jobs, err := getTestJobs(kube, event)
	if err != nil {
		return err
	}
	updateJobs(kube, jobs)

	return nil
}

func getTestJobs(kube client.JobInterface, event string) ([]batch.Job, error) {
	jobs, err := kube.List(api.ListOptions{})
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to get list jobs (%v)", err))
	}

	var ret []batch.Job
	for _, currJob := range jobs.Items {
		if currJob.ObjectMeta.Labels[LabelTugbotEvents] == event &&
			currJob.ObjectMeta.Labels[LabelTugbotCreatedFrom] == "" {
			ret = append(ret, currJob)
		}
	}

	return ret, nil
}

func updateJobs(kube client.JobInterface, jobs []batch.Job) {
	for _, currJob := range jobs {
		newJob := createJobFrom(currJob)
		_, err := kube.Create(&newJob)
		if err != nil {
			log.Errorf("Failed to create rerun job (Original: %+v New: %+v). %v", currJob, newJob, err)
		}
	}
}

func createJobFrom(job batch.Job) batch.Job {
	from := job.Name
	job.Name = fmt.Sprintf("tugbot.%s.%s", from, time.Now())
	job.Labels[LabelTugbotCreatedFrom] = from

	return job
}
