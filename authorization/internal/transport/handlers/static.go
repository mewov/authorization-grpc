package handlers

import (
	"context"
	"log/slog"

	"github.com/gox7/authorizathion/internal/services"
	"github.com/gox7/authorizathion/models"
	"google.golang.org/grpc/peer"
)

var (
	authService    *services.AuthorizathionService
	sessionService *services.SessionsService
	config         *models.LocalConfig
)

func SetResource(localConfig *models.LocalConfig, auth *services.AuthorizathionService, session *services.SessionsService) {
	authService = auth
	sessionService = session
	config = localConfig
}

func logger(message string, logger *slog.Logger, ctx context.Context) {
	p, ok := peer.FromContext(ctx)
	if ok {
		network := p.Addr.Network()
		address := p.Addr.String()

		logger.Info(message,
			slog.String("address", address),
			slog.String("network", network),
		)
	}
}
