package auth

import (
	"context"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	pagerAuth "pager-services/pkg/api/pager_api/auth"
	common "pager-services/pkg/api/pager_api/common"
	"pager-services/pkg/transfers"
	"pager-services/pkg/utils"
	"time"
)

var _ pagerAuth.AuthServiceServer = (*PagerAuth)(nil)

type PagerAuth struct {
}

func (p PagerAuth) Refresh(ctx context.Context, request *pagerAuth.RefreshRequest) (*pagerAuth.RefreshResponse, error) {
	accessToken := request.AccessToken
	validAccessToken, err := utils.ValidateAccessToken(accessToken)

	if err != nil {
		refreshToken := request.RefreshToken
		validateRefresh, err := utils.ValidateRefreshToken(refreshToken)
		if err != nil {
			return nil, err
		}
		_, err = transfers.CheckRefreshToken(ctx, refreshToken)
		if err != nil {
			return nil, err
		}
		newAccessToken, err := utils.RefreshAccessToken(validateRefresh)
		if err != nil {
			return nil, err
		}
		return &pagerAuth.RefreshResponse{AccessToken: newAccessToken}, nil
	}
	return &pagerAuth.RefreshResponse{AccessToken: validAccessToken.Raw}, nil
}

func (p PagerAuth) Registration(ctx context.Context, request *pagerAuth.RegistrationRequest) (*common.Empty, error) {
	if request.GetEmail() == "" {
		return nil, status.Error(codes.InvalidArgument, "email is require")
	}
	if request.GetLogin() == "" {
		return nil, status.Error(codes.InvalidArgument, "login is require")
	}
	if request.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "password is require")
	}
	passHash, err := bcrypt.GenerateFromPassword([]byte(request.GetPassword()), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("%s", err)
	}

	authData := &transfers.AuthRegisterData{
		Email:    request.GetEmail(),
		Login:    request.GetLogin(),
		Password: string(passHash),
	}

	exists, err := transfers.IsUserExistsWithData(ctx, authData.Email, authData.Login)
	if err != nil {
		return nil, err
	}

	if exists {
		return nil, utils.MentorError("user already exists", codes.AlreadyExists, &common.PagerError{
			Code: common.PagerError_ALREADY_EXISTS,
		})
	}

	// Вставляем данные в коллекции
	err = transfers.InsertAuthData(ctx, authData)
	if err != nil {
		return nil, utils.MentorError("failed to register user", codes.Internal, &common.PagerError{
			Code: common.PagerError_INTERNAL,
		})
	}

	return &common.Empty{}, nil
}

func (p PagerAuth) Login(ctx context.Context, request *pagerAuth.LoginRequest) (*pagerAuth.Token, error) {
	if request.GetIdentity() == "" {
		return nil, utils.MentorError("identity require", codes.InvalidArgument, &common.PagerError{
			Code: common.PagerError_INVALID_ARGUMENT,
		})
	}
	if request.GetPassword() == "" {
		return nil, utils.MentorError("password require", codes.InvalidArgument, &common.PagerError{
			Code: common.PagerError_INVALID_ARGUMENT,
		})
	}
	authData := &transfers.AuthLoginData{
		Password: request.GetPassword(),
		Identity: request.GetIdentity(),
	}

	UserId, err := transfers.FindUserIDByIdentifier(ctx, authData.Identity)
	if err != nil {
		return nil, err
	}
	passHash, refreshToken, err := transfers.GetHashedPasswordByIDAndRefreshToken(ctx, UserId)
	if err != nil {
		return nil, utils.MentorError("password not found", codes.NotFound, &common.PagerError{
			Code:    common.PagerError_NOT_FOUND,
			Details: err.Error(),
		})
	}
	if err := bcrypt.CompareHashAndPassword(passHash, []byte(authData.Password)); err != nil {
		return nil, utils.MentorError("unknown error", codes.Unknown, &common.PagerError{
			Code:    common.PagerError_UNKNOWN,
			Details: err.Error(),
		})
	}
	AccessToken, err := utils.NewToken(UserId, authData.Identity, 5*time.Minute)
	if err != nil {
		return nil, err
	}
	return &pagerAuth.Token{
		RefreshToken: refreshToken,
		AccessToken:  AccessToken,
	}, nil
}
