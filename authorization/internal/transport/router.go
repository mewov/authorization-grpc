package transport

import (
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"

	proto "github.com/gox7/AuthorizationRPC/proto/gen"
	"github.com/gox7/authorizathion/internal/services"
	"github.com/gox7/authorizathion/internal/transport/handlers"
	"github.com/gox7/authorizathion/models"
	"google.golang.org/grpc"
)

var (
	config *models.LocalConfig
	s1     *services.AuthorizathionService
	s2     *services.SessionsService
)

func SetResource(localConfig *models.LocalConfig, auth *services.AuthorizathionService, session *services.SessionsService) {
	config = localConfig
	s2 = session
	s1 = auth
}

func Register(logger *slog.Logger, gRPC *grpc.Server) {
	handlers.SetResource(config, s1, s2)
	auth := handlers.NewAuthEngine(logger)
	session := handlers.NewSessionEngine(logger)

	proto.RegisterAuthorizationServer(gRPC, auth)
	proto.RegisterSessionServer(gRPC, session)
}

func Listen(logger *slog.Logger, gRPC *grpc.Server) {
	listener, err := net.Listen("tcp", config.SERVER_ADDRESS)
	if err != nil {
		logger.Error("tcp new listener: "+err.Error(), slog.String("addr", config.SERVER_ADDRESS))
		fmt.Println("[-] transport.tcp:", err.Error())
		return
	}

	logger.Info("grpc listen...", slog.String("addr", config.SERVER_ADDRESS))
	fmt.Println("[+] transport.listen:", config.SERVER_ADDRESS)
	if err := gRPC.Serve(listener); err != nil {
		logger.Error("grpc listen: "+err.Error(), slog.String("addr", config.SERVER_ADDRESS))
		fmt.Println("[-] transport.listen:", err.Error())
	}
}

func GracefulShutdown(logger *slog.Logger, gRPC *grpc.Server) {
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGTERM, syscall.SIGINT)

	<-shutdown
	logger.Info("grpc graceful shutdown")
	gRPC.GracefulStop()
}
