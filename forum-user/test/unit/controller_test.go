package unit

import (
	"forum-user/internal/model"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestControllerCreateUser(t *testing.T) {
	testConfig.PrepareDB()

	err := testConfig.userClient.CreateUser(1, "123456")
	assert.NoError(t, err)

	user, err := testConfig.userClient.GetUserByAccountID(1)
	assert.NoError(t, err)

	assert.Equal(t, 1, user.AccountID)
	assert.Equal(t, "123456", user.Login)
	assert.Equal(t, "", user.HeaderUrl)
	assert.Equal(t, "", user.AvatarUrl)
}
func TestControllerBanUser(t *testing.T) {
	testConfig.PrepareDB()

	err := testConfig.userClient.CreateUser(1, "user1")
	assert.NoError(t, err)

	err = testConfig.userClient.CreateUser(2, "user2")
	assert.NoError(t, err)

	err = testConfig.userClient.BanUser(1, "user2")
	assert.NoError(t, err)

	userBans, err := testConfig.userClient.BanByAccountID("user1")
	assert.NoError(t, err)

	assert.Equal(t, 2, userBans.UserBans[0].FromAccountID)
	assert.Equal(t, 1, userBans.UserBans[0].ToAccountID)
}

func TestControllerUnbanUser(t *testing.T) {
	testConfig.PrepareDB()

	err := testConfig.userClient.CreateUser(1, "user1")
	assert.NoError(t, err)

	err = testConfig.userClient.CreateUser(2, "user2")
	assert.NoError(t, err)

	err = testConfig.userClient.BanUser(1, "user2")
	assert.NoError(t, err)

	err = testConfig.userClient.UnbanUser(1, "user2")

	userBans, err := testConfig.userClient.BanByAccountID("user1")
	assert.NoError(t, err)

	assert.Equal(t, 0, len(userBans.UserBans))
}

func TestControllerNewWarningFromAdmin(t *testing.T) {
	err := testConfig.userClient.CreateUser(1, "user1")
	assert.NoError(t, err)

	err = testConfig.userClient.NewWarningFromAdmin(1, "BULLING", "admin1", "admin")
	assert.NoError(t, err)

	userWarnings, err := testConfig.userClient.AllWarningFromAdmin(1)
	assert.NoError(t, err)

	assert.Equal(t, 3, userWarnings.UserWarnings[0].AdminAccountID)
	assert.Equal(t, "admin", userWarnings.UserWarnings[0].AdminLogin)
	assert.Equal(t, 1, userWarnings.UserWarnings[0].ToAccountID)
	assert.Equal(t, model.BullingWarningType, userWarnings.UserWarnings[0].WarningType)
	assert.Equal(t, model.WarningMap[userWarnings.UserWarnings[0].WarningType], userWarnings.UserWarnings[0].WarningText)
}

func TestControllerUploadAvatar(t *testing.T) {
	testConfig.PrepareDB()

	err := testConfig.userClient.CreateUser(1, "123456")
	assert.NoError(t, err)

	avatarFile := []byte{1, 2}
	err = testConfig.userClient.UploadAvatar(avatarFile, 1, "user1")

	user, err := testConfig.userClient.GetUserByAccountID(1)
	assert.NoError(t, err)

	assert.Equal(t, true, user.AvatarUrl != "")
	assert.Equal(t, "", user.HeaderUrl)
}

func TestControllerUploadHeader(t *testing.T) {
	testConfig.PrepareDB()

	err := testConfig.userClient.CreateUser(1, "123456")
	assert.NoError(t, err)

	headerFile := []byte{1, 2}
	err = testConfig.userClient.UploadHeader(headerFile, 1, "user1")

	user, err := testConfig.userClient.GetUserByAccountID(1)
	assert.NoError(t, err)

	assert.Equal(t, true, user.HeaderUrl != "")
	assert.Equal(t, "", user.AvatarUrl)
}

func TestControllerDownloadHeader(t *testing.T) {
	testConfig.PrepareDB()

	err := testConfig.userClient.CreateUser(1, "123456")
	assert.NoError(t, err)

	headerFile := []byte{1, 2}
	err = testConfig.userClient.UploadHeader(headerFile, 1, "user1")

	user, err := testConfig.userClient.GetUserByAccountID(1)
	assert.NoError(t, err)

	headerFileResponse, err := testConfig.userClient.DownloadHeader(user.HeaderUrl)
	assert.NoError(t, err)

	assert.Equal(t, headerFile, headerFileResponse)
}

func TestControllerDownloadAvatar(t *testing.T) {
	testConfig.PrepareDB()

	err := testConfig.userClient.CreateUser(1, "123456")
	assert.NoError(t, err)

	avatarFile := []byte{1, 2}
	err = testConfig.userClient.UploadAvatar(avatarFile, 1, "user1")

	user, err := testConfig.userClient.GetUserByAccountID(1)
	assert.NoError(t, err)

	avatarFileResponse, err := testConfig.userClient.DownloadAvatar(user.AvatarUrl)
	assert.NoError(t, err)

	assert.Equal(t, avatarFile, avatarFileResponse)
}

func TestControllerUsersByLoginSearch(t *testing.T) {
	testConfig.PrepareDB()

	err := testConfig.userClient.CreateUser(1, "приветик")
	assert.NoError(t, err)

	err = testConfig.userClient.CreateUser(2, "приветствую")
	assert.NoError(t, err)

	err = testConfig.userClient.CreateUser(3, "приветули")
	assert.NoError(t, err)

	err = testConfig.userClient.CreateUser(4, "парапара")
	assert.NoError(t, err)

	err = testConfig.userClient.CreateUser(5, "привет")
	assert.NoError(t, err)

	time.Sleep(15000)

	users, err := testConfig.userClient.UsersByLoginSearch("привет")
	assert.NoError(t, err)

	assert.Equal(t, 3, len(users.Users))
}

func TestControllerUsersByLogin(t *testing.T) {
	testConfig.PrepareDB()

	err := testConfig.userClient.CreateUser(1, "привет")
	assert.NoError(t, err)

	users, err := testConfig.userClient.UserByLogin("привет")
	assert.NoError(t, err)

	assert.Equal(t, "привет", users.User[0].Login)
	assert.Equal(t, 1, users.User[0].AccountID)
}
