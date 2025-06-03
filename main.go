package main

import (
	"context"
	"github.com/igntnk/stocky-2pc-controller/clients"
	"github.com/igntnk/stocky-2pc-controller/config"
	"github.com/igntnk/stocky-2pc-controller/controllers"
	"github.com/igntnk/stocky-2pc-controller/grpc"
	"github.com/igntnk/stocky-2pc-controller/services"
	"github.com/igntnk/stocky-2pc-controller/web"
	"github.com/rs/zerolog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	var l = zerolog.New(os.Stdout).With().Timestamp().Logger()

	mainCtx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	cfg, err := config.Get()
	if err != nil {
		l.Fatal().Err(err).Send()
		return
	}

	iimsConn, err := grpc.NewGrpcClientConn(
		mainCtx,
		cfg.IIMS.Address,
		cfg.IIMS.Timeout,
		cfg.IIMS.Tries,
		cfg.IIMS.Insecure,
	)
	if err != nil {
		l.Fatal().Err(err).Send()
		return
	}
	iimsClient := clients.NewIIMSClient(iimsConn)

	scsConn, err := grpc.NewGrpcClientConn(
		mainCtx,
		cfg.SCS.Address,
		cfg.SCS.Timeout,
		cfg.SCS.Tries,
		cfg.SCS.Insecure,
	)
	if err != nil {
		l.Fatal().Err(err).Send()
		return
	}
	scsClient := clients.NewSCSClient(scsConn)

	smsConn, err := grpc.NewGrpcClientConn(
		mainCtx,
		cfg.SMS.Address,
		cfg.SMS.Timeout,
		cfg.SMS.Tries,
		cfg.SMS.Insecure,
	)
	if err != nil {
		l.Fatal().Err(err).Send()
		return
	}
	smsClient := clients.NewSMSClient(smsConn)

	omsConn, err := grpc.NewGrpcClientConn(
		mainCtx,
		cfg.OMS.Address,
		cfg.OMS.Timeout,
		cfg.OMS.Tries,
		cfg.OMS.Insecure,
	)
	if err != nil {
		l.Fatal().Err(err).Send()
		return
	}
	clients.NewOMSClient(omsConn)
	omsClient := clients.NewOMSClient(omsConn)

	orderService := services.NewOrderService(omsClient, smsClient, scsClient, iimsClient)

	orderController := controllers.NewOrderController(orderService)

	httpServer, err := web.New(l, cfg.RestServer.Port, orderController)
	if err != nil {
		l.Fatal().Err(err).Send()
		return
	}

	serverErrorChan := make(chan error, 1)
	go func() {
		serverErrorChan <- httpServer.ListenAndServe()
	}()
	l.Info().Msgf("Server started on port: %s", cfg.RestServer.Port)

	select {
	case err := <-serverErrorChan:
		l.Error().Err(err).Msgf("http server is shuting down")
	case <-mainCtx.Done():
		l.Info().Msg("shutting down http server, press Ctrl+C again to force")
	}
}
