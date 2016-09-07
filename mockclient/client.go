package mockclient

import (
	"errors"

	"github.com/stretchr/testify/mock"
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/apis/batch"
	"k8s.io/kubernetes/pkg/watch"
)

type MockClient struct {
	mock.Mock
}

func NewMockClient() *MockClient {
	return &MockClient{}
}

// k8s.io/kubernetes/pkg/client/unversioned -> JobInterface

func (m *MockClient) List(opts api.ListOptions) (*batch.JobList, error) {
	args := m.Called(opts)
	return args.Get(0).(*batch.JobList), args.Error(1)
}

func (m *MockClient) Get(name string) (*batch.Job, error) {
	return nil, errors.New("Not Implemented")
}

func (m *MockClient) Create(job *batch.Job) (*batch.Job, error) {
	args := m.Called(job)
	return args.Get(0).(*batch.Job), args.Error(1)
}

func (m *MockClient) Update(job *batch.Job) (*batch.Job, error) {
	return nil, errors.New("Not Implemented")
}

func (m *MockClient) Delete(name string, options *api.DeleteOptions) error {
	return errors.New("Not Implemented")
}

func (m *MockClient) Watch(opts api.ListOptions) (watch.Interface, error) {
	return nil, errors.New("Not Implemented")
}

func (m *MockClient) UpdateStatus(job *batch.Job) (*batch.Job, error) {
	return nil, errors.New("Not Implemented")
}
