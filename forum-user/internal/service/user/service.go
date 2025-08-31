package user

import (
	"context"
	"errors"
	"forum-user/internal/model"
)

func New(
	userRepo model.IUserRepo,
	notificationClient model.INotificationClient,
) *ServiceUser {
	return &ServiceUser{
		userRepo:           userRepo,
		notificationClient: notificationClient,
	}
}

type ServiceUser struct {
	userRepo           model.IUserRepo
	notificationClient model.INotificationClient
}

func (s *ServiceUser) CreateUser(ctx context.Context,
	accountID int,
	login string,
) (int, error) {
	userID, err := s.userRepo.CreateUser(ctx, accountID, login)
	if err != nil {
		return 0, err
	}

	err = s.userRepo.CreateUserSearchIndex(ctx, accountID, login)
	if err != nil {
		return 0, err
	}
	return userID, nil
}

func (s *ServiceUser) BanUser(ctx context.Context,
	fromAccountID,
	toAccountID int,
) error {
	err := s.userRepo.CreateUserBan(ctx, fromAccountID, toAccountID)
	if err != nil {
		return err
	}

	return nil
}

func (s *ServiceUser) NewWarningFromAdmin(ctx context.Context,
	adminAccountID,
	toAccountID int,
	warningType,
	adminLogin string,
) error {
	warningText := model.WarningMap[warningType]
	err := s.userRepo.CreateWarningFromAdmin(
		ctx,
		adminAccountID,
		toAccountID,
		warningType,
		warningText,
		adminLogin,
	)
	if err != nil {
		return err
	}

	err = s.notificationClient.NewWarningFromAdminNotification(ctx,
		adminAccountID,
		toAccountID,
		warningText,
		adminLogin,
	)
	if err != nil {
		return err
	}

	return nil
}

func (s *ServiceUser) UploadAvatar(ctx context.Context,
	accountID int,
	avatar []byte,
) error {
	fileID, err := s.userRepo.UploadAvatar(ctx, accountID, avatar)
	if err != nil {
		return err
	}

	err = s.userRepo.UpdateAvatarUrl(ctx, accountID, fileID)
	if err != nil {
		return err
	}
	return nil
}

func (s *ServiceUser) UploadHeader(ctx context.Context,
	accountID int,
	header []byte,
) error {
	fileID, err := s.userRepo.UploadHeader(ctx, accountID, header)
	if err != nil {
		return err
	}

	err = s.userRepo.UpdateHeaderUrl(ctx, accountID, fileID)
	if err != nil {
		return err
	}
	return nil
}

func (s *ServiceUser) DownloadAvatar(ctx context.Context,
	fileID string,
) ([]byte, error) {
	avatar, err := s.userRepo.DownloadAvatar(ctx, fileID)
	if err != nil {
		return nil, err
	}
	return avatar, err
}

func (s *ServiceUser) DownloadHeader(ctx context.Context,
	fileID string,
) ([]byte, error) {
	header, err := s.userRepo.DownloadHeader(ctx, fileID)
	if err != nil {
		return nil, err
	}
	return header, err
}

func (s *ServiceUser) UserByAccountID(ctx context.Context,
	accountID int,
) (*model.User, error) {
	user, err := s.userRepo.UserByAccountID(ctx, accountID)
	if err != nil {
		return nil, err
	}
	return user[0], nil
}

func (s *ServiceUser) UserByLogin(ctx context.Context,
	login string,
) ([]*model.User, error) {
	user, err := s.userRepo.UserByLogin(ctx, login)
	if err != nil {
		return nil, err
	}
	if len(user) == 0 {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (s *ServiceUser) UsersByLoginSearch(ctx context.Context,
	login string,
) ([]*model.UserSearch, error) {
	users, err := s.userRepo.UsersByLoginSearch(ctx, login)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (s *ServiceUser) BanByAccountID(ctx context.Context,
	toAccountID int,
) ([]*model.UserBan, error) {
	userBans, err := s.userRepo.BanByAccountID(ctx, toAccountID)
	if err != nil {
		return nil, err
	}

	return userBans, nil
}

func (s *ServiceUser) AllWarningFromAdmin(ctx context.Context,
	toAccountID int,
) ([]*model.WarningFromAdmin, error) {
	userWarnings, err := s.userRepo.AllWarningFromAdmin(ctx, toAccountID)
	if err != nil {
		return nil, err
	}

	return userWarnings, nil
}

func (s *ServiceUser) UnbanUser(ctx context.Context,
	fromAccountID,
	toAccountID int,
) error {
	err := s.userRepo.DeleteUserBan(ctx, fromAccountID, toAccountID)
	if err != nil {
		return err
	}

	return nil
}
