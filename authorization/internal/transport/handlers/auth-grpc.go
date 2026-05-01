package handlers

import (
	"context"
	"log/slog"
	"time"

	rpc "github.com/gox7/AuthorizationRPC/proto/gen"
	"github.com/gox7/authorizathion/pkg/tokens"
	"github.com/gox7/authorizathion/pkg/validator"
)

type (
	AuthEngine struct {
		logger *slog.Logger
		rpc.UnimplementedAuthorizationServer
	}
)

func NewAuthEngine(logger *slog.Logger) *AuthEngine {
	logger.Info("auth.engine: new gRPC handler")
	return &AuthEngine{logger: logger}
}

func (auth *AuthEngine) Register(ctx context.Context, request *rpc.RequestRegister) (*rpc.ResponseToken, error) {
	logger("grpc.register", auth.logger, ctx)
	if err := validator.ValidateLogin(request.Login); err != nil {
		return nil, err
	}
	if err := validator.ValidatePassword(request.Password); err != nil {
		return nil, err
	}
	if err := validator.Email(request.Email); err != nil {
		return nil, err
	}

	if request.Client == "" {
		request.Client = "none"
	}
	if request.Role == "" {
		request.Role = "user"
	}

	lastId, err := authService.CreateUser(
		request.Login, request.Email, request.Password,
		request.Client, request.Role,
	)
	if err != nil {
		return nil, err
	}

	refresh := tokens.GenerateRefresh()
	expires := time.Now().Add(time.Hour * 24 * 7)
	if _, err := sessionService.CreateSession(lastId, refresh, request.Client, expires.Unix()); err != nil {
		return nil, err
	}

	access, err := tokens.GenerateAccess(config, lastId, request.Login, request.Email, request.Role)
	if err != nil {
		return nil, err
	}

	return &rpc.ResponseToken{
		Status:  "success",
		Message: "successfuly create account",
		Tokens: &rpc.Tokens{
			Refresh: refresh,
			Access:  access,
		},
	}, nil
}

func (auth *AuthEngine) Login(ctx context.Context, request *rpc.RequestLogin) (*rpc.ResponseToken, error) {
	logger("grpc.login", auth.logger, ctx)
	if err := validator.ValidateLogin(request.Login); err != nil {
		return nil, err
	}
	if err := validator.ValidatePassword(request.Password); err != nil {
		return nil, err
	}

	user, err := authService.SearchUser(request.Login, request.Password)
	if err != nil {
		return nil, err
	}

	refresh := tokens.GenerateRefresh()
	expires := time.Now().Add(time.Hour * 24 * 7)
	if _, err := sessionService.CreateSession(user.ID, refresh, user.Client, expires.Unix()); err != nil {
		return nil, err
	}

	access, err := tokens.GenerateAccess(config, user.ID, user.Login, user.Email, user.Role)
	if err != nil {
		return nil, err
	}

	return &rpc.ResponseToken{
		Status:  "success",
		Message: "successfuly login account",
		Tokens: &rpc.Tokens{
			Refresh: refresh,
			Access:  access,
		},
	}, nil
}
