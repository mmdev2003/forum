package unit

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestControllerCreateAdmin(t *testing.T) {
	testConfig.PrepareDB()

	err := testConfig.adminClient.CreateAdmin(1)
	assert.NoError(t, err)

	err = testConfig.adminClient.CreateAdmin(2)
	assert.NoError(t, err)

	admins, err := testConfig.adminClient.AllAdmin()
	assert.NoError(t, err)

	assert.Equal(t, 2, len(admins.Admins))
}
