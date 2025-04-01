package authclient

import (
	"chatY-go/internal/api/grpc/auth/proto"
	"context"
	"google.golang.org/grpc"
)

type AuthClient struct {
	client proto.AuthServiceClient
}

func NewAuthClient(conn *grpc.ClientConn) *AuthClient {
	return &AuthClient{client: proto.NewAuthServiceClient(conn)}
}

func (a *AuthClient) Login(username string, password string) (bool, string, error) {
	res, err := a.client.Login(context.TODO(), &proto.LoginRequest{
		Username: username,
		Password: password,
	})
	if err != nil {
		return false, "", err
	}
	return res.Success, res.Message, nil
}
