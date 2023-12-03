package app

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"path/filepath"
	"strconv"
	"time"

	pb "example.com/go_backend/internal/grpc/authen"
	"example.com/go_backend/internal/models"
	"example.com/go_backend/internal/store"
	"github.com/go-redis/redis"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type myAuthenServer struct {
	pb.UnimplementedAuthenticationServer
}

func (s *myAuthenServer) GetAuthenLoginFeedBack(ctx context.Context, user *pb.User) (*pb.UserAuthenResponse, error) {
	res := &pb.UserAuthenResponse{}
	gormUser := models.User{}
	rdb := store.GetRedisClient()
	var key string
	if user.Phone != 0 {
		key = strconv.Itoa(int(user.GetPhone()))
	} else {
		key = user.GetName()
	}
	val, err := rdb.Get(key).Result()
	if err != nil {
		if err == redis.Nil {
			fmt.Println("key does not exists")
		}
		panic(err)
	}
	if len(val) != 0 {
		err := json.Unmarshal([]byte(val), &gormUser)
		if err != nil {
			val = ""
			fmt.Println("value is invalid in redis")
		}
	}
	if len(val) == 0 {
		pgDB := store.GetPgConn()
		if user.Phone != 0 {
			pgDB.Where("phone = ?", user.Phone).Preload("UserInfo.MemberRight").Take(&gormUser)
		} else {
			pgDB.Where("name = ?", user.Name).Preload("UserInfo.MemberRight").Take(&gormUser)
		}
		if v, err := json.Marshal(&gormUser); err == nil {
			rdb.Set(gormUser.Name, string(v), time.Millisecond*50)
			rdb.Set(strconv.Itoa(int(gormUser.Phone)), string(v), time.Millisecond*50)
		}

	}
	return res, nil
}

func AuthenServerRPC() {
	viper.SetConfigFile("./config.yaml")
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err))
	}
	tls := viper.Get("tls").(bool)
	certFile := viper.Get("certFile").(string)
	keyFile := viper.Get("keyFile").(string)
	serverAddr := viper.Get("serverAddr").(string)
	lis, err := net.Listen("tcp4", serverAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	if tls {
		certFile, _ = filepath.Abs(certFile)
		keyFile, _ = filepath.Abs(keyFile)
		creds, err := credentials.NewServerTLSFromFile(certFile, keyFile)
		if err != nil {
			log.Fatalf("Failed to generate credentials: %v", err)
		}
		opts = []grpc.ServerOption{grpc.Creds(creds)}
	}
	grpcServer := grpc.NewServer(opts...)
	myServer := &myAuthenServer{}
	pb.RegisterAuthenticationServer(grpcServer, myServer)
	grpcServer.Serve(lis)
}
