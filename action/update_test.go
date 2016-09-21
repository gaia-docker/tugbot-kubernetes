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
	assert.Error(t, err)
}

func TestUpdateJobsEmptyEvent(t *testing.T) {
	err := UpdateJobs(mockclient.NewMockClient(), &api.Event{})
	assert.Error(t, err)
}

func TestUpdateJobsEventWithNoReason(t *testing.T) {
	err := UpdateJobs(mockclient.NewMockClient(), &api.Event{InvolvedObject: api.ObjectReference{Kind: "ReplicaSet"}})
	assert.Error(t, err)
}

func TestUpdateJobsEventWithNoKind(t *testing.T) {
	err := UpdateJobs(mockclient.NewMockClient(), &api.Event{InvolvedObject: api.ObjectReference{Kind: ""}, Reason: "Created"})
	assert.Error(t, err)
}

func TestUpdateJobsFailedToGetListJobs(t *testing.T) {
	kube := mockclient.NewMockClient()
	kube.On("List", mock.Anything).Return(&batch.JobList{}, errors.New("Expected")).Once()
	err := UpdateJobs(kube, &api.Event{InvolvedObject: api.ObjectReference{Kind: "ReplicaSet"}, Reason: "SuccessfulCreate"})
	assert.Error(t, err)
	kube.AssertExpectations(t)
}

func TestUpdateJobsNoTugbotEventLabelDefinedOnJob(t *testing.T) {
	kube := mockclient.NewMockClient()
	kube.On("List", mock.Anything).Return(
		&batch.JobList{Items: []batch.Job{batch.Job{}}}, nil).Once()
	err := UpdateJobs(kube, &api.Event{InvolvedObject: api.ObjectReference{Kind: "ReplicaSet"}, Reason: "SuccessfulCreate"})
	assert.NoError(t, err)
	kube.AssertExpectations(t)
}

func TestUpdateJobsTugbotEventNotMatch(t *testing.T) {
	kube := mockclient.NewMockClient()
	kube.On("List", mock.Anything).Return(
		&batch.JobList{Items: []batch.Job{
			batch.Job{ObjectMeta: api.ObjectMeta{Labels: map[string]string{LabelTugbotEvents: "Node.NodeHasSufficientDisk"}}}}},
		nil).Once()
	err := UpdateJobs(kube, &api.Event{InvolvedObject: api.ObjectReference{Kind: "ReplicaSet"}, Reason: "SuccessfulCreate"})
	assert.NoError(t, err)
	kube.AssertExpectations(t)
}

func TestUpdateJobsCreatedByTugbot(t *testing.T) {
	kube := mockclient.NewMockClient()
	kube.On("List", mock.Anything).Return(
		&batch.JobList{Items: []batch.Job{
			batch.Job{ObjectMeta: api.ObjectMeta{Labels: map[string]string{
				LabelTugbotEvents: "ReplicaSet.SuccessfulCreate,Node.NodeHasSufficientDisk", LabelTugbotTriggerBy: "ReplicaSet.SuccessfulCreate"}}}}},
		nil).Once()
	err := UpdateJobs(kube, &api.Event{InvolvedObject: api.ObjectReference{Kind: "ReplicaSet"}, Reason: "SuccessfulCreate"})
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
				Labels: map[string]string{LabelTugbotEvents: "Node.NodeHasSufficientDisk,ReplicaSet.SuccessfulCreate"}}}}},
		nil).Once()
	kube.On("Create", mock.Anything).Run(func(args mock.Arguments) {
		assert.True(t, strings.HasPrefix(args.Get(0).(*batch.Job).Name,
			fmt.Sprintf("tugbot.%s.", name)))
	}).Return(&batch.Job{}, nil).Once()
	err := UpdateJobs(kube, &api.Event{InvolvedObject: api.ObjectReference{Kind: "ReplicaSet"}, Reason: "SuccessfulCreate"})
	assert.NoError(t, err)
	kube.AssertExpectations(t)
}

func TestUpdateJobsErrorCreatingJob(t *testing.T) {
	kube := mockclient.NewMockClient()
	kube.On("List", mock.Anything).Return(
		&batch.JobList{Items: []batch.Job{
			batch.Job{ObjectMeta: api.ObjectMeta{Labels: map[string]string{
				LabelTugbotEvents: "ReplicaSet.SuccessfulCreate"}}}}},
		nil).Once()
	kube.On("Create", mock.Anything).Return(&batch.Job{}, errors.New("Expected")).Once()
	err := UpdateJobs(kube, &api.Event{InvolvedObject: api.ObjectReference{Kind: "ReplicaSet"}, Reason: "SuccessfulCreate"})
	assert.NoError(t, err)
	kube.AssertExpectations(t)
}
