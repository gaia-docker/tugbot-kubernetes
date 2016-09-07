package action

import (
	"errors"
	"testing"

	"github.com/gaia-docker/tugbot-kubernetes/mockclient"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/apis/batch"
)

func TestUpdateJobs_FailedToGetListJobs(t *testing.T) {
	kube := mockclient.NewMockClient()
	kube.On("List", mock.Anything).Return(&batch.JobList{}, errors.New("Expected"))
	err := UpdateJobs(kube, "")
	assert.Error(t, err)
	kube.AssertExpectations(t)
}

func TestUpdateJobsNoTugbotEventDefined(t *testing.T) {
	kube := mockclient.NewMockClient()
	kube.On("List", mock.Anything).Return(
		&batch.JobList{Items: []batch.Job{batch.Job{}}}, nil)
	err := UpdateJobs(kube, "")
	assert.NoError(t, err)
	kube.AssertExpectations(t)
}

func TestUpdateJobsTugbotEventNotMatch(t *testing.T) {
	kube := mockclient.NewMockClient()
	kube.On("List", mock.Anything).Return(
		&batch.JobList{Items: []batch.Job{
			batch.Job{ObjectMeta: api.ObjectMeta{Labels: map[string]string{LabelTugbotEvents: "create-pod"}}}}},
		nil)
	err := UpdateJobs(kube, "deployment")
	assert.NoError(t, err)
	kube.AssertExpectations(t)
}
