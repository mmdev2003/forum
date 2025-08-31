package frame

import (
	"context"
	"forum-frame/internal/model"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"time"
)

func New(
	db model.IDatabase,
	weedFS model.IWeedFS,
) *RepoFrame {
	return &RepoFrame{db, weedFS}
}

type RepoFrame struct {
	db     model.IDatabase
	weedFS model.IWeedFS
}

func (frameRepo *RepoFrame) CreateFrame(ctx context.Context,
	frameID,
	accountID,
	paymentID int,
	expirationAt time.Time,
) (int, error) {
	args1 := pgx.NamedArgs{
		"frame_id":       frameID,
		"account_id":     accountID,
		"payment_id":     paymentID,
		"expiration_at":  expirationAt,
		"payment_status": model.Pending,
	}
	dbFrameID, err := frameRepo.db.Insert(ctx, CreateFrame, args1)
	if err != nil {
		return 0, err
	}

	args2 := pgx.NamedArgs{
		"db_frame_id": dbFrameID,
		"account_id":  accountID,
	}

	_, err = frameRepo.db.Insert(ctx, CreateCurrentFrame, args2)
	if err != nil {
		return 0, err
	}

	return dbFrameID, nil
}

func (frameRepo *RepoFrame) CreateFrameData(ctx context.Context,
	monthlyPrice,
	foreverPrice float32,
	name,
	fileID string,
) error {
	args := pgx.NamedArgs{
		"monthly_price": monthlyPrice,
		"forever_price": foreverPrice,
		"name":          name,
		"file_id":       fileID,
	}
	_, err := frameRepo.db.Insert(ctx, CreateFrameData, args)
	if err != nil {
		return err
	}
	return nil
}

func (frameRepo *RepoFrame) ConfirmPaymentForFrame(ctx context.Context,
	paymentID int,
) error {
	args := pgx.NamedArgs{
		"payment_id":     paymentID,
		"payment_status": model.Confirmed,
	}
	err := frameRepo.db.Update(ctx, ConfirmPaymentForFrame, args)
	if err != nil {
		return err
	}

	return nil
}

func (frameRepo *RepoFrame) ChangeCurrentFrame(ctx context.Context,
	dbFrameID int,
	accountID int,
) error {
	args := pgx.NamedArgs{
		"db_frame_id": dbFrameID,
		"account_id":  accountID,
	}
	err := frameRepo.db.Update(ctx, ChangeCurrentFrame, args)
	if err != nil {
		return err
	}
	return nil
}

func (frameRepo *RepoFrame) FramesByAccountID(ctx context.Context,
	accountID int,
) ([]*model.Frame, error) {
	args := pgx.NamedArgs{
		"account_id": accountID,
	}
	rows, err := frameRepo.db.Select(ctx, FramesByAccountID, args)
	if err != nil {
		return nil, err
	}

	var frames []*model.Frame
	err = pgxscan.ScanAll(&frames, rows)
	if err != nil {
		return nil, err
	}

	return frames, nil
}

func (frameRepo *RepoFrame) UploadFrame(frameFile []byte, name string) (string, error) {
	fileID, err := frameRepo.weedFS.Upload(frameFile, name, int64(len(frameFile)), "")
	if err != nil {
		return "", err
	}
	return fileID, nil
}

func (frameRepo *RepoFrame) AllFrame(ctx context.Context) ([]*model.FrameData, error) {
	rows, err := frameRepo.db.Select(ctx, AllFrames, pgx.NamedArgs{})
	if err != nil {
		return nil, err
	}
	var framesData []*model.FrameData
	err = pgxscan.ScanAll(&framesData, rows)
	if err != nil {
		return nil, err
	}

	return framesData, nil
}

func (frameRepo *RepoFrame) FrameDataByID(ctx context.Context,
	frameID int,
) ([]*model.FrameData, error) {
	args := pgx.NamedArgs{
		"frame_id": frameID,
	}
	rows, err := frameRepo.db.Select(ctx, GetFrameDataByID, args)
	if err != nil {
		return nil, err
	}

	var frameData []*model.FrameData
	err = pgxscan.ScanAll(&frameData, rows)
	if err != nil {
		return nil, err
	}
	return frameData, nil
}

func (frameRepo *RepoFrame) CurrentFrameByAccountID(ctx context.Context,
	accountID int,
) ([]*model.CurrentFrame, error) {
	args := pgx.NamedArgs{
		"account_id": accountID,
	}
	rows, err := frameRepo.db.Select(ctx, CurrentFrameByAccountID, args)
	if err != nil {
		return nil, err
	}

	var currentFrame []*model.CurrentFrame
	err = pgxscan.ScanAll(&currentFrame, rows)
	if err != nil {
		return nil, err
	}

	return currentFrame, nil
}

func (frameRepo *RepoFrame) DownloadFrame(
	ctx context.Context,
	frameID int,
) ([]byte, error) {
	frameData, err := frameRepo.FrameDataByID(ctx, frameID)
	if err != nil {
		return nil, err
	}
	file, err := frameRepo.weedFS.Download(frameData[0].FileID)
	if err != nil {
		return nil, err
	}

	return file, nil
}
