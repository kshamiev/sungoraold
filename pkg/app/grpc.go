package app

import (
	"context"
	"fmt"
	"net"

	"google.golang.org/grpc"
	grpcMetadata "google.golang.org/grpc/metadata"

	"github.com/kshamiev/sungora/pkg/app/response"
	"github.com/kshamiev/sungora/pkg/logger"
)

// клиент GRPC
type GRPCClient struct {
	Conn *grpc.ClientConn
}

// NewGRPCClient создание и старт клиента GRPC
func NewGRPCClient(cfg *ConfigGRPC) (comp *GRPCClient, err error) {
	comp = &GRPCClient{}
	comp.Conn, err = grpc.Dial(fmt.Sprintf("%s:%d", cfg.Host, cfg.Port), grpc.WithInsecure())
	return comp, err
}

// Wait завершение работы клиента GRPC
func (comp *GRPCClient) Wait() {
	_ = comp.Conn.Close()
}

// сервер GRPC
type GRPCServer struct {
	Addr      string        // адрес сервера grpc
	chControl chan struct{} // управление ожиданием завершения работы сервера
	lis       net.Listener
}

// NewGRPCServer создание и старт сервера GRPC
func NewGRPCServer(cfg *ConfigGRPC, serve *grpc.Server, lg logger.Logger) (comp *GRPCServer, err error) {
	comp = &GRPCServer{
		Addr:      fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		chControl: make(chan struct{}),
	}

	if comp.lis, err = net.Listen("tcp", comp.Addr); err != nil {
		return
	}

	go func() {
		_ = serve.Serve(comp.lis)
		close(comp.chControl)
	}()

	lg.Info("start grpc server: ", comp.Addr)

	return comp, nil
}

// Wait завершение работы сервера GRPC
func (comp *GRPCServer) Wait(lg logger.Logger) {
	if comp.lis == nil {
		return
	}

	if err := comp.lis.Close(); err != nil {
		return
	}

	<-comp.chControl
	lg.Info("stop grpc server: ", comp.Addr)
}

type GRPCKit struct{}

func (kit *GRPCKit) CtxOut(ctx context.Context, m map[string]string) context.Context {
	if m == nil {
		m = make(map[string]string)
	}
	m[string(response.CtxUUID)] = ctx.Value(response.CtxUUID).(string)
	m[string(response.CtxAPI)] = ctx.Value(response.CtxAPI).(string)

	return grpcMetadata.NewOutgoingContext(ctx, grpcMetadata.New(m))
}

func (kit *GRPCKit) CtxIn(ctx context.Context) grpcMetadata.MD {
	if md, ok := grpcMetadata.FromIncomingContext(ctx); ok {
		return md
	} else {
		return nil
	}
}

func (kit *GRPCKit) Log(md grpcMetadata.MD, lg logger.Logger) logger.Logger {
	if md.Get(response.LogUUID) != nil && md.Get(response.LogAPI) != nil {
		lg = lg.WithField(response.LogUUID, md.Get(response.LogUUID)[0])
		lg = lg.WithField(response.LogAPI, md.Get(response.LogAPI)[0])
	}
	return lg
}
