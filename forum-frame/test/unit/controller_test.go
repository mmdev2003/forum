package unit

import (
	"forum-frame/internal/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestControllerCreatePaymentForFrame(t *testing.T) {
	testConfig.PrepareDB()

	frameFile := []byte{1, 2}
	err := testConfig.frameClient.AddNewFrame(frameFile, 1.23, 2.34, "frame1")

	_, err = testConfig.frameClient.CreatePaymentForFrame(1, 1, "btc", "user1")
	assert.NoError(t, err)

	accountFrames, err := testConfig.frameClient.FramesByAccountID(1)
	assert.NoError(t, err)

	assert.Equal(t, model.Pending, accountFrames.Frames[0].PaymentStatus)
	assert.Equal(t, 1, accountFrames.Frames[0].FrameID)
	assert.Equal(t, 1, accountFrames.Frames[0].AccountID)
	assert.Equal(t, 1, accountFrames.CurrentFrame.DbFrameID)
}

func TestControllerConfirmPaymentForFrame(t *testing.T) {
	testConfig.PrepareDB()

	frameFile := []byte{1, 2}
	err := testConfig.frameClient.AddNewFrame(frameFile, 1.23, 2.34, "frame1")

	_, err = testConfig.frameClient.CreatePaymentForFrame(1, 1, "btc", "user1")
	assert.NoError(t, err)

	err = testConfig.frameClient.ConfirmPaymentForFrame(1, testConfig.interServerSecretKey)
	assert.NoError(t, err)

	accountFrames, err := testConfig.frameClient.FramesByAccountID(1)
	assert.NoError(t, err)

	assert.Equal(t, model.Confirmed, accountFrames.Frames[0].PaymentStatus)
}

func TestControllerChangeCurrentFrame(t *testing.T) {
	testConfig.PrepareDB()

	frameFile1 := []byte{1, 2}
	err := testConfig.frameClient.AddNewFrame(frameFile1, 1.23, 2.34, "frame1")

	frameFile2 := []byte{1, 2, 3}
	err = testConfig.frameClient.AddNewFrame(frameFile2, 1.23, 2.34, "frame2")
	assert.NoError(t, err)

	_, err = testConfig.frameClient.CreatePaymentForFrame(1, 1, "btc", "user1")
	assert.NoError(t, err)

	_, err = testConfig.frameClient.CreatePaymentForFrame(2, 1, "btc", "user1")
	assert.NoError(t, err)

	err = testConfig.frameClient.ChangeCurrentFrame(2, "user1")
	assert.NoError(t, err)

	accountFrames, err := testConfig.frameClient.FramesByAccountID(1)
	assert.NoError(t, err)

	assert.Equal(t, 2, len(accountFrames.Frames))
	assert.Equal(t, 2, accountFrames.CurrentFrame.DbFrameID)
}

func TestControllerAllFrame(t *testing.T) {
	testConfig.PrepareDB()

	_, err := testConfig.frameClient.AllFrame()
	assert.NoError(t, err)
}

func TestControllerAddNewFrame(t *testing.T) {
	testConfig.PrepareDB()

	frameFile1 := []byte{1, 2}
	err := testConfig.frameClient.AddNewFrame(frameFile1, 1.23, 2.34, "frame1")

	frameFile2 := []byte{1, 2, 3}
	err = testConfig.frameClient.AddNewFrame(frameFile2, 1.23, 2.34, "frame2")

	frames, err := testConfig.frameClient.AllFrame()
	assert.NoError(t, err)

	assert.Equal(t, 2, len(frames.Frames))

}

func TestControllerDownloadFrame(t *testing.T) {
	testConfig.PrepareDB()

	frameFile1 := []byte{1, 2}
	err := testConfig.frameClient.AddNewFrame(frameFile1, 1.23, 2.34, "frame1")

	frameFile2 := []byte{1, 2, 3}
	err = testConfig.frameClient.AddNewFrame(frameFile2, 1.23, 2.34, "frame2")

	frame, err := testConfig.frameClient.DownloadFrame(1)
	assert.NoError(t, err)
	assert.Equal(t, frameFile1, frame)

	frame, err = testConfig.frameClient.DownloadFrame(2)
	assert.NoError(t, err)
	assert.Equal(t, frameFile2, frame)
}
