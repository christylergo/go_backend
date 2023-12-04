package app

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
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
	"github.com/redis/go-redis/v9"
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
	val, err := rdb.Get(context.Background(), key).Result()
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
			rdb.Set(context.Background(), gormUser.Name, string(v), time.Millisecond*50)
			rdb.Set(context.Background(), strconv.Itoa(int(gormUser.Phone)), string(v), time.Millisecond*50)
		}

	}
	h := sha256.New()
	_, err = h.Write([]byte(user.PassWord + gormUser.CreatedAt.String()))
	if err != nil {
		return res, err
	}
	err = nil
	if gormUser.PassWord == hex.EncodeToString(h.Sum(nil)) {
		res.ID = uint64(gormUser.ID)
		res.Token, err = generateTokenUsingHs256(&gormUser)
	}
	return res, err
}

func (s *myAuthenServer) GetAuthenRegisterFeedBack(ctx context.Context, user *pb.User) (*pb.UserAuthenResponse, error) {
	gormUser := models.User{}
	gormUser.Name = user.GetName()
	gormUser.Phone = uint(user.GetPhone())
	gormUser.Email = user.Email
	gormUser.PassWord = user.PassWord
	pgDB := store.GetPgConn()
	result := pgDB.Create(&gormUser)
	if result.Error != nil {
		log.Print(result.Error.Error())
	}
	pgDB.Where("phone = ?", user.Phone).Preload("UserInfo.MemberRight").Take(&gormUser)
	rdb := store.GetRedisClient()
	if v, err := json.Marshal(&gormUser); err == nil {
		rdb.Set(context.Background(), gormUser.Name, string(v), time.Millisecond*50)
		rdb.Set(context.Background(), strconv.Itoa(int(gormUser.Phone)), string(v), time.Millisecond*50)
	}
	res := &pb.UserAuthenResponse{}
	res.ID = uint64(gormUser.ID)
	var terr error
	res.Token, terr = generateTokenUsingHs256(&gormUser)
	return res, terr
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
