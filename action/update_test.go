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
	kube.On("List", mock.Anything).Return(&batch.JobList{}, errors.New("Expected")).Once()
	err := UpdateJobs(kube, "")
	assert.Error(t, err)
	kube.AssertExpectations(t)
}

func TestUpdateJobsNoTugbotEventDefined(t *testing.T) {
	kube := mockclient.NewMockClient()
	kube.On("List", mock.Anything).Return(
		&batch.JobList{Items: []batch.Job{batch.Job{}}}, nil).Once()
	err := UpdateJobs(kube, "")
	assert.NoError(t, err)
	kube.AssertExpectations(t)
}

func TestUpdateJobsTugbotEventNotMatch(t *testing.T) {
	kube := mockclient.NewMockClient()
	kube.On("List", mock.Anything).Return(
		&batch.JobList{Items: []batch.Job{
			batch.Job{ObjectMeta: api.ObjectMeta{Labels: map[string]string{LabelTugbotEvents: "create-pod"}}}}},
		nil).Once()
	err := UpdateJobs(kube, "deployment")
	assert.NoError(t, err)
	kube.AssertExpectations(t)
}

func TestUpdateJobsCreatedByTugbot(t *testing.T) {
	kube := mockclient.NewMockClient()
	kube.On("List", mock.Anything).Return(
		&batch.JobList{Items: []batch.Job{
			batch.Job{ObjectMeta: api.ObjectMeta{Labels: map[string]string{
				LabelTugbotEvents: "deployment", LabelTugbotCreatedFrom: "testing-pi"}}}}},
		nil).Once()
	err := UpdateJobs(kube, "deployment")
	assert.NoError(t, err)
	kube.AssertExpectations(t)
}

func TestUpdateJobs(t *testing.T) {
	kube := mockclient.NewMockClient()
	kube.On("List", mock.Anything).Return(
		&batch.JobList{Items: []batch.Job{
			batch.Job{ObjectMeta: api.ObjectMeta{Labels: map[string]string{
				LabelTugbotEvents: "deployment"}}}}},
		nil).Once()
	kube.On("Create", mock.Anything).Return(&batch.Job{}, nil).Once()
	err := UpdateJobs(kube, "deployment")
	assert.NoError(t, err)
	kube.AssertExpectations(t)
}

func TestUpdateJobsErrorCreatingJob(t *testing.T) {
	kube := mockclient.NewMockClient()
	kube.On("List", mock.Anything).Return(
		&batch.JobList{Items: []batch.Job{
			batch.Job{ObjectMeta: api.ObjectMeta{Labels: map[string]string{
				LabelTugbotEvents: "deployment"}}}}},
		nil).Once()
	kube.On("Create", mock.Anything).Return(&batch.Job{}, errors.New("Expected")).Once()
	err := UpdateJobs(kube, "deployment")
	assert.NoError(t, err)
	kube.AssertExpectations(t)
}
