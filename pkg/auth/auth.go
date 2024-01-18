package auth

import (
	context "context"
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

func (p PagerAuth) SearchUsersByIdentifier(ctx context.Context, request *pagerAuth.SearchUsersRequest) (*pagerAuth.SearchUsersResponse, error) {
	userIds, err := transfers.FindUserIDsByIdentifier(ctx, "test", request.GetIdentifier())
	if err != nil {
		return nil, err
	}

	// Возвращаем список ID в ответе
	response := &pagerAuth.SearchUsersResponse{
		UserIds: userIds,
	}

	return response, nil
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
	exists, err := transfers.IsUserExistsWithData(ctx, "test", authData.Email, authData.Login)
	if err != nil {
		return nil, err
	}

	if exists {
		return nil, status.Error(codes.AlreadyExists, "user already exists")
	}

	// Вставляем данные в коллекции
	err = transfers.InsertAuthData(ctx, authData)
	if err != nil {
		return nil, fmt.Errorf("failed to register user: %s", err)
	}

	return &common.Empty{}, nil
}

func (p PagerAuth) Logout(ctx context.Context, token *pagerAuth.Token) (*common.Empty, error) {
	//TODO implement me
	panic("implement me")
}

func (p PagerAuth) Login(ctx context.Context, request *pagerAuth.LoginRequest) (*pagerAuth.Token, error) {
	if request.GetIdentity() == "" {
		return nil, status.Error(codes.InvalidArgument, "identity is require")
	}
	if request.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "password is require")
	}
	authData := &transfers.AuthLoginData{
		Password: request.GetPassword(),
		Identity: request.GetIdentity(),
	}

	UserId, err := transfers.FindUserIDByIdentifier(ctx, "test", authData.Identity)
	if err != nil {
		return nil, status.Error(codes.NotFound, "user id not found")
	}
	passHash, err := transfers.GetHashedPasswordByID(ctx, UserId)
	if err != nil {
		return nil, status.Error(codes.NotFound, "password not found")
	}

	if err := bcrypt.CompareHashAndPassword(passHash, []byte(authData.Password)); err != nil {
		return nil, err
	}

	token, err := utils.NewToken(UserId, authData.Identity, 5*time.Minute)
	if err != nil {
		return nil, err
	}
	return &pagerAuth.Token{
		Token: token,
	}, nil
}
