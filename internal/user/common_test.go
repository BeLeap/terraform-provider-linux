package user

import (
	"terraform-provider-linux/internal/util"
	"testing"

	"gotest.tools/assert"
	is "gotest.tools/assert/cmp"
)

func TestGet(t *testing.T) {
	desired := &LinuxUser{
		Username: "root",
		Uid:      0,
		Gid:      0,
	}

	linuxContext := util.GetLinuxContextForTest(t)
	username := "root"
	user, err := Get(linuxContext, username)

	if err != nil {
		t.Fatalf("Failed to get user '%s': %v", username, err)
	}

	assert.DeepEqual(t, desired, user)
}

func TestGetInvalidUser(t *testing.T) {
	linuxContext := util.GetLinuxContextForTest(t)
	username := "user_not_exists"
	user, err := Get(linuxContext, username)

	assert.Assert(t, is.Nil(user))
	assert.Assert(t, is.Nil(err))
}
