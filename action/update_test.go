package action

import (
	"errors"
	"testing"

	"github.com/gaia-docker/tugbot-kubernetes/mockclient"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"k8s.io/kubernetes/pkg/apis/batch"
)

func TestUpdateJobs(t *testing.T) {
	kube := mockclient.NewMockClient()
	kube.On("List", mock.Anything).Return(&batch.JobList{}, errors.New("Expected"))
	err := UpdateJobs(kube, "")
	assert.Error(t, err)
	kube.AssertExpectations(t)
}
