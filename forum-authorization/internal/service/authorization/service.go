package authorization

import (
	"context"
	"errors"
	"forum-authorization/internal/model"
	"github.com/golang-jwt/jwt/v4"
	"log/slog"
	"time"
)

func New(
	authorizationRepo model.IAuthorizationRepo,
	jwtSecretKey string,
) *ServiceAuth {
	return &ServiceAuth{
		authorizationRepo: authorizationRepo,
		jwtSecretKey:      jwtSecretKey,
	}
}

type ServiceAuth struct {
	authorizationRepo model.IAuthorizationRepo
	jwtSecretKey      string
}

func (authorizationService *ServiceAuth) CreateTokens(
	ctx context.Context,
	accountID int,
	role string,
	twoFaStatus bool,
) (*model.JWTTokens, error) {
	account, err := authorizationService.authorizationRepo.AccountByID(ctx, accountID)
	if err != nil {
		return nil, err
	}

	if len(account) == 0 {
		err = authorizationService.authorizationRepo.SetAccount(ctx, accountID)
		if err != nil {
			return nil, err
		}
	}
	accessToken, err := authorizationService.createJWTToken(accountID, role, twoFaStatus, 15)
	if err != nil {
		return nil, err
	}

	refreshToken, err := authorizationService.createJWTToken(accountID, role, twoFaStatus, 30)
	if err != nil {
		return nil, err
	}

	err = authorizationService.authorizationRepo.UpdateRefreshToken(ctx, accountID, refreshToken)
	if err != nil {
		return nil, err
	}

	return &model.JWTTokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (authorizationService *ServiceAuth) CheckToken(
	token string,
) (*model.TokenPayload, error) {
	_token, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, model.ErrTokenInvalid
		}
		return []byte(authorizationService.jwtSecretKey), nil
	})

	if err != nil {
		var errVal *jwt.ValidationError
		if errors.As(err, &errVal) {
			if errVal.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, model.ErrTokenExpired
			}
		}
		return nil, model.ErrTokenInvalid
	}

	if claims, ok := _token.Claims.(jwt.MapClaims); ok && _token.Valid {
		tokenPayload := &model.TokenPayload{
			AccountID:   int(claims["accountID"].(float64)),
			Role:        claims["role"].(string),
			TwoFaStatus: claims["twoFaStatus"].(bool),
			Exp:         int64(claims["exp"].(float64)),
		}
		return tokenPayload, nil
	}

	return nil, model.ErrTokenInvalid
}

func (authorizationService *ServiceAuth) RefreshTokens(ctx context.Context,
	refreshToken string,
) (*model.JWTTokens, error) {
	account, err := authorizationService.authorizationRepo.AccountByRefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, err
	}
	if len(account) == 0 {
		return nil, model.ErrAccountNotFound
	}

	tokenPayload, err := authorizationService.CheckToken(refreshToken)
	if err != nil {
		return nil, err
	}

	jwtTokens, err := authorizationService.CreateTokens(ctx,
		tokenPayload.AccountID,
		tokenPayload.Role,
		tokenPayload.TwoFaStatus,
	)
	if err != nil {
		return nil, err
	}

	return jwtTokens, nil
}

func (authorizationService *ServiceAuth) createJWTToken(
	accountID int,
	role string,
	twoFaStatus bool,
	minutes int,
) (string, error) {
	expirationTime := time.Now().Add(time.Duration(minutes) * time.Minute)
	claims := jwt.MapClaims{
		"accountID":   accountID,
		"role":        role,
		"twoFaStatus": twoFaStatus,
		"exp":         expirationTime.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwtToken, err := token.SignedString([]byte(authorizationService.jwtSecretKey))
	if err != nil {
		slog.Error(err.Error())
		return "", err
	}

	return jwtToken, nil
}
