package admin

import (
	"forum-admin/internal/model"
	"golang.org/x/net/context"
)

func New(
	adminRepo model.IAdminRepo,
) *ServiceAdmin {
	return &ServiceAdmin{
		adminRepo: adminRepo,
	}
}

type ServiceAdmin struct {
	adminRepo model.IAdminRepo
}

func (s *ServiceAdmin) CreateAdmin(ctx context.Context,
	accountID int,
) (int, error) {
	adminID, err := s.adminRepo.CreateAdmin(ctx, accountID)
	if err != nil {
		return 0, err
	}

	return adminID, nil
}

func (s *ServiceAdmin) AllAdmin(ctx context.Context,
) ([]*model.Admin, error) {
	admins, err := s.adminRepo.AllAdmin(ctx)
	if err != nil {
		return nil, err
	}
	return admins, nil
}
