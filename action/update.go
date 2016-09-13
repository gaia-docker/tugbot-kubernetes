package action

import (
	"errors"
	"fmt"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/gaia-docker/tugbot-kubernetes/common"
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/apis/batch"
	client "k8s.io/kubernetes/pkg/client/unversioned"
	"strings"
)

const (
	LabelTugbotEvents      = "tugbot.kubernetes.events"
	LabelTugbotCreatedFrom = "tugbot.created.from"
)

func UpdateJobs(kube client.JobInterface, event *api.Event) error {
	if event == nil {
		return nil
	}

	log.Infof("Event: %s:%s, %s", event.InvolvedObject.Kind, event.Reason, event.Message)
	jobs, err := getTestJobs(kube, event)
	if err != nil {
		return err
	}
	updateJobs(kube, jobs)

	return nil
}

func getTestJobs(kube client.JobInterface, event *api.Event) ([]batch.Job, error) {
	jobs, err := kube.List(api.ListOptions{})
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to get list jobs (%v)", err))
	}

	var ret []batch.Job
	for _, currJob := range jobs.Items {
		if currJob.Labels[LabelTugbotCreatedFrom] == "" &&
			isJobContainsEvent(currJob, event) {
			ret = append(ret, currJob)
		}
	}

	return ret, nil
}

func isJobContainsEvent(job batch.Job, event *api.Event) bool {
	ret := false
	if job.Labels != nil {
		jobEvents, ok := job.Labels[LabelTugbotEvents]
		ret = ok && common.SliceContains(toString(event), strings.Split(jobEvents, ","))
	}

	return ret
}

func toString(event *api.Event) string {
	return fmt.Sprintf("%s:%s", event.InvolvedObject.Kind, event.Reason)
}

func updateJobs(kube client.JobInterface, jobs []batch.Job) {
	for _, currJob := range jobs {
		newJob := createJobFrom(currJob)
		log.Infof("Creating job... (Original: %+v New: %+v)", currJob, newJob)
		_, err := kube.Create(&newJob)
		if err != nil {
			log.Errorf("Update job failed (Original: %+v New: %+v). %v", currJob, newJob, err)
		}
	}
}

func createJobFrom(job batch.Job) batch.Job {
	from := job.Name
	job.Name = fmt.Sprintf("tugbot.%s.%d", from, time.Now().UnixNano())
	if job.Labels == nil {
		job.Labels = make(map[string]string)
	}
	job.Labels[LabelTugbotCreatedFrom] = from

	return job
}
