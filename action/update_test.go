package action

import (
	"errors"
	"testing"

	"fmt"
	"github.com/gaia-docker/tugbot-kubernetes/mockclient"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/apis/batch"
	"strings"
)

func TestUpdateJobsNilEvent(t *testing.T) {
	err := UpdateJobs(mockclient.NewMockClient(), nil)
	assert.NoError(t, err)
}

func TestUpdateJobsFailedToGetListJobs(t *testing.T) {
	kube := mockclient.NewMockClient()
	kube.On("List", mock.Anything).Return(&batch.JobList{}, errors.New("Expected")).Once()
	err := UpdateJobs(kube, &api.Event{})
	assert.Error(t, err)
	kube.AssertExpectations(t)
}

func TestUpdateJobsNoTugbotEventDefined(t *testing.T) {
	kube := mockclient.NewMockClient()
	kube.On("List", mock.Anything).Return(
		&batch.JobList{Items: []batch.Job{batch.Job{}}}, nil).Once()
	err := UpdateJobs(kube, &api.Event{})
	assert.NoError(t, err)
	kube.AssertExpectations(t)
}

func TestUpdateJobsEmptyEvent(t *testing.T) {
	kube := mockclient.NewMockClient()
	kube.On("List", mock.Anything).Return(
		&batch.JobList{Items: []batch.Job{
			batch.Job{ObjectMeta: api.ObjectMeta{Labels: map[string]string{LabelTugbotEvents: "NodeReady"}}}}},
		nil).Once()
	err := UpdateJobs(kube, &api.Event{})
	assert.NoError(t, err)
	kube.AssertExpectations(t)
}

func TestUpdateJobsTugbotEventNotMatch(t *testing.T) {
	kube := mockclient.NewMockClient()
	kube.On("List", mock.Anything).Return(
		&batch.JobList{Items: []batch.Job{
			batch.Job{ObjectMeta: api.ObjectMeta{Labels: map[string]string{LabelTugbotEvents: "create-pod"}}}}},
		nil).Once()
	err := UpdateJobs(kube, &api.Event{Reason: "NodeReady"})
	assert.NoError(t, err)
	kube.AssertExpectations(t)
}

func TestUpdateJobsCreatedByTugbot(t *testing.T) {
	kube := mockclient.NewMockClient()
	kube.On("List", mock.Anything).Return(
		&batch.JobList{Items: []batch.Job{
			batch.Job{ObjectMeta: api.ObjectMeta{Labels: map[string]string{
				LabelTugbotEvents: "NodeReady", LabelTugbotCreatedFrom: "testing-pi"}}}}},
		nil).Once()
	err := UpdateJobs(kube, &api.Event{Reason: "NodeReady"})
	assert.NoError(t, err)
	kube.AssertExpectations(t)
}

func TestUpdateJobs(t *testing.T) {
	const name = "test-job"
	kube := mockclient.NewMockClient()
	kube.On("List", mock.Anything).Return(
		&batch.JobList{Items: []batch.Job{
			batch.Job{ObjectMeta: api.ObjectMeta{
				Name:   name,
				Labels: map[string]string{LabelTugbotEvents: "NodeReady,NodeHasSufficientDisk"}}}}},
		nil).Once()
	kube.On("Create", mock.Anything).Run(func(args mock.Arguments) {
		assert.True(t, strings.HasPrefix(args.Get(0).(*batch.Job).Name,
			fmt.Sprintf("tugbot.%s.", name)))
	}).Return(&batch.Job{}, nil).Once()
	err := UpdateJobs(kube, &api.Event{Reason: "NodeHasSufficientDisk"})
	assert.NoError(t, err)
	kube.AssertExpectations(t)
}

func TestUpdateJobsErrorCreatingJob(t *testing.T) {
	kube := mockclient.NewMockClient()
	kube.On("List", mock.Anything).Return(
		&batch.JobList{Items: []batch.Job{
			batch.Job{ObjectMeta: api.ObjectMeta{Labels: map[string]string{
				LabelTugbotEvents: "NodeReady"}}}}},
		nil).Once()
	kube.On("Create", mock.Anything).Return(&batch.Job{}, errors.New("Expected")).Once()
	err := UpdateJobs(kube, &api.Event{Reason: "NodeReady"})
	assert.NoError(t, err)
	kube.AssertExpectations(t)
}
