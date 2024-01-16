package auth

import (
	context "context"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	pagerAuth "pager-services/pkg/api/pager_api/auth"
	common "pager-services/pkg/api/pager_api/common"
	"pager-services/pkg/mongo_ops"
	"pager-services/pkg/transfers"
)

var _ pagerAuth.AuthServiceServer = (*PagerAuth)(nil)

type PagerAuth struct {
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

	authData := &transfers.AuthData{
		Email:    request.GetEmail(),
		Login:    request.GetLogin(),
		Password: string(passHash),
	}

	// Вставляем данные в коллекции
	err = transfers.InsertAuthData(ctx, mongo_ops.CollectionsPoll.ProfileCollection, mongo_ops.CollectionsPoll.UsersCollection, authData)
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
	//TODO implement me
	panic("implement me")
}
