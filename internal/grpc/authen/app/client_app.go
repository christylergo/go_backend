package app

import (
	"context"
	"fmt"
	"log"
	"path/filepath"
	"time"

	pb "example.com/go_backend/internal/grpc/authen"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

func AuthenLoginRPC(user *pb.User) (int, string) {
	conn := getConn()
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	client := pb.NewAuthenticationClient(conn)
	res, err := client.GetAuthenLoginFeedBack(ctx, user)
	if err != nil {
		log.Fatalf("client.GetAuthenFeedBack failed: %v", err)
	}
	id := int(res.GetID())
	token := res.GetToken()
	return id, token
}
func AuthenRegisterRPC(user *pb.User) (int, string) {
	conn := getConn()
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	client := pb.NewAuthenticationClient(conn)
	res, err := client.GetAuthenRegisterFeedBack(ctx, user)
	if err != nil {
		log.Fatalf("client.GetAuthenFeedBack failed: %v", err)
	}
	id := int(res.GetID())
	token := res.GetToken()
	return id, token
}

func getConn() *grpc.ClientConn {
	viper.SetConfigFile("./config.yaml")
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err))
	}
	tls := viper.Get("tls").(bool)
	caFile := viper.Get("caFile").(string)
	serverAddr := viper.Get("serverAddr").(string)
	serverHostOverride := viper.Get("serverHostOverride").(string)

	var opts []grpc.DialOption
	if tls {
		if caFile == "" {
			caFile, _ = filepath.Abs("x509/ca_cert.pem")
		}
		creds, err := credentials.NewClientTLSFromFile(caFile, serverHostOverride)
		if err != nil {
			log.Fatalf("Failed to create TLS credentials: %v", err)
		}
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	conn, err := grpc.Dial(serverAddr, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	return conn
}
