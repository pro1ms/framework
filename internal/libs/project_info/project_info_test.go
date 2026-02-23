package project_info

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetProjectInfo(t *testing.T) {
	info, err := GetProjectInfo("./")
	assert.NoError(t, err)
	assert.NotNil(t, info)
	assert.Equal(t, "github.com/pro1ms/framework", info.ModuleName)

	info, err = GetProjectInfo("/Users/emris/www/")
	assert.Error(t, err)
	assert.Nil(t, info)

}
