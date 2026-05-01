package handlers

import (
	"context"
	"errors"
	"log/slog"
	"time"

	rpc "github.com/gox7/AuthorizationRPC/proto/gen"
	"github.com/gox7/authorizathion/pkg/tokens"
)

type (
	SessionEngine struct {
		logger *slog.Logger
		rpc.UnimplementedSessionServer
	}
)

func NewSessionEngine(logger *slog.Logger) *SessionEngine {
	logger.Info("session.engine: new gRPC handler")
	return &SessionEngine{logger: logger}
}

func (session *SessionEngine) Refresh(ctx context.Context, request *rpc.RequestRefresh) (*rpc.ResponseToken, error) {
	logger("grpc.refresh", session.logger, ctx)
	if request.Refresh == "" {
		return nil, errors.New("refresh token is null")
	}

	refresh, err := sessionService.SearchSession(request.Refresh)
	if err != nil {
		return nil, err
	}
	if err := sessionService.RemoveSession(refresh.Token); err != nil {
		return nil, err
	}

	user, err := authService.SearchUserByID(refresh.UserID)
	if err != nil {
		return nil, err
	}

	token := tokens.GenerateRefresh()
	expires := time.Now().Add(time.Hour * 24 * 7)
	if _, err := sessionService.CreateSession(user.ID, token, user.Client, expires.Unix()); err != nil {
		return nil, err
	}

	access, err := tokens.GenerateAccess(config, user.ID, user.Login, user.Email, user.Role)
	if err != nil {
		return nil, err
	}

	return &rpc.ResponseToken{
		Status:  "success",
		Message: "successfuly refresh token",
		Tokens: &rpc.Tokens{
			Refresh: token,
			Access:  access,
		},
	}, nil
}

func (session *SessionEngine) Logout(ctx context.Context, request *rpc.RequestLogout) (*rpc.ResponseLogout, error) {
	logger("grpc.logout", session.logger, ctx)
	if request.Refresh == "" {
		return nil, errors.New("refresh token is null")
	}

	if err := sessionService.RemoveSession(request.Refresh); err != nil {
		return nil, err
	}

	return &rpc.ResponseLogout{
		Status:  "success",
		Message: "refresh token is remove",
	}, nil
}

func (session *SessionEngine) Info(ctx context.Context, request *rpc.RequstInfo) (*rpc.ResponseInfo, error) {
	logger("grpc.info", session.logger, ctx)
	claims, err := tokens.CheckAccess(config, request.Access)
	if err != nil {
		return nil, err
	}

	return &rpc.ResponseInfo{
		Status:  "success",
		Message: "successfuly information from access",
		User: &rpc.User{
			UserId: claims.UserId,
			Login:  claims.Login,
			Email:  claims.Email,
			Role:   claims.Role,
		},
	}, nil
}
