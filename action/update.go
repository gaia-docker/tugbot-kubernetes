package action

import (
	"errors"
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/gaia-docker/tugbot-kubernetes/common"
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/unversioned"
	"k8s.io/kubernetes/pkg/apis/batch"
	client "k8s.io/kubernetes/pkg/client/unversioned"
	"strings"
	"time"
)

const (
	LabelTugbotEvents        = "tugbot.kubernetes.events"
	LabelTugbotTriggerBy     = "tugbot.trigger.by"
	LabelTugbotTriggerByName = "tugbot.trigger.by.name"
)

func UpdateJobs(kube client.JobInterface, event *api.Event) error {
	if event == nil || event.InvolvedObject.Kind == "" || event.Reason == "" {
		return errors.New("Invalid event")
	}

	log.Debugf("Event: %s:%s, %s", event.InvolvedObject.Kind, event.Reason, event.Message)
	jobs, err := getTestJobs(kube, event)
	if err != nil {
		return err
	}
	updateJobs(kube, jobs, event)

	return nil
}

func getTestJobs(kube client.JobInterface, event *api.Event) ([]batch.Job, error) {
	jobs, err := kube.List(api.ListOptions{})
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to get list jobs (%v)", err))
	}

	var ret []batch.Job
	for _, currJob := range jobs.Items {
		if currJob.Labels[LabelTugbotTriggerBy] == "" &&
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
	return fmt.Sprintf("%s.%s", event.InvolvedObject.Kind, event.Reason)
}

func updateJobs(kube client.JobInterface, jobs []batch.Job, event *api.Event) {
	for _, currJob := range jobs {
		newJob := createJobFrom(currJob, event)
		log.Infof("Creating a job %s (from %s)...", newJob.Name, currJob.Name)
		_, err := kube.Create(&newJob)
		if err != nil {
			log.Errorf("Job creation failed: %v (Original: %+v New: %+v)", err, currJob, newJob)
		} else {
			log.Debug("New job %s created (Original: %+v New: %+v)", newJob.Name, currJob, newJob)
		}
	}
}

func createJobFrom(job batch.Job, event *api.Event) batch.Job {
	ret := job
	ret.Name = fmt.Sprintf("tugbot.%s.%d", job.Name, time.Now().UnixNano())
	if ret.Labels == nil {
		ret.Labels = make(map[string]string)
	}
	ret.Labels[LabelTugbotTriggerBy] = toString(event)
	ret.Labels[LabelTugbotTriggerByName] = event.Name
	if ret.Spec.Template.Labels != nil {
		ret.Spec.Template.Labels = make(map[string]string)
	}
	if ret.Spec.Selector != nil && ret.Spec.Selector.MatchLabels != nil {
		ret.Spec.Selector.MatchLabels = make(map[string]string)
	}
	ret.Status = batch.JobStatus{}
	ret.SelfLink = ""
	ret.ResourceVersion = ""
	ret.CreationTimestamp = unversioned.Time{}
	ret.DeletionTimestamp = &unversioned.Time{}

	return ret
}
