package action

import (
	"io/ioutil"

	"errors"
	"fmt"
	"github.com/ghodss/yaml"
	"k8s.io/kubernetes/pkg/apis/batch"
)

func ToJob(jobYaml string) (*batch.Job, error) {
	file, err := ioutil.ReadFile(jobYaml)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to read file (%v)", err))
	}

	var ret batch.Job
	if err := yaml.Unmarshal(file, &ret); err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to unmarshal Job (%v)", err))
	}

	return &ret, nil
}
