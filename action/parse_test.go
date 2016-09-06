package action

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestToJob(t *testing.T) {
	job, err := ToJob("test-job.yaml")
	assert.NoError(t, err)
	assert.NotEqual(t, "", job.ObjectMeta.GetLabels()[LabelTugbotEvents])
}
