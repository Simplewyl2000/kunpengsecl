/***
Description: Using grpc to implement the service API
***/

package clientapi

import (
	"context"
	"errors"
	"log"
	"net"
	"time"

	"gitee.com/openeuler/kunpengsecl/attestation/tas/akissuer"
	"google.golang.org/grpc"
)

type (
	service struct {
		UnimplementedTasServer
	}
	tasConn struct {
		ctx    context.Context
		cancel context.CancelFunc
		conn   *grpc.ClientConn
		c      TasClient
	}
)

func (s *service) GetAKCert(ctx context.Context, in *GetAKCertRequest) (*GetAKCertReply, error) {
	akcert, err := akissuer.GenerateAKCert(in.Akcert, in.Dvcert)
	if err != nil {
		return nil, err
	}
	return &GetAKCertReply{Akcert: akcert}, nil
}

func (s *service) RegisterClient(ctx context.Context, in *RegisterClientRequest) (*RegisterClientReply, error) {
	id, err := akissuer.RegisterClient(in.GetClientinfo(), in.GetDvcert())
	if err != nil {
		return nil, err
	}
	return &RegisterClientReply{Clientid: id}, nil
}

func StartServer(addr string) {
	log.Print("Start tee ak server...")
	listen, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Server: fail to listen at %s, %v", addr, err)
	}
	akServer := grpc.NewServer()
	RegisterTasServer(akServer, &service{})
	if err := akServer.Serve(listen); err != nil {
		log.Fatalf("Server: fail to serve, %v", err)
	}
}

func makesock(addr string) (*tasConn, error) {
	tas := &tasConn{}
	tas.ctx, tas.cancel = context.WithTimeout(context.Background(), 60*time.Second)
	conn, err := grpc.DialContext(tas.ctx, addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, errors.New("Client: fail to connect " + addr)
	}
	tas.conn = conn
	tas.c = NewTasClient(conn)
	log.Printf("Client: connect to %s", addr)
	return tas, nil
}

func DoGetAKCert(addr string, in *GetAKCertRequest) (*GetAKCertReply, error) {
	tas, err := makesock(addr)
	if err != nil {
		return nil, err
	}
	defer tas.conn.Close()
	defer tas.cancel()

	rpy, err := tas.c.GetAKCert(tas.ctx, in)
	if err != nil {
		return nil, err
	}
	return rpy, nil
}

func DoRegisterClient(addr string, in *RegisterClientRequest) (*RegisterClientReply, error) {
	tas, err := makesock(addr)
	if err != nil {
		return nil, err
	}
	defer tas.conn.Close()
	defer tas.cancel()

	rpy, err := tas.c.RegisterClient(tas.ctx, in)
	if err != nil {
		return nil, err
	}
	return rpy, nil
}