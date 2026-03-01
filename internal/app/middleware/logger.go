package middleware

import (
	"context"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

func Logger(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	str, _ := protojson.Marshal((req).(proto.Message))
	log.Printf("request: method: %s, req: %s\n", info.FullMethod, str)

	if resp, err = handler(ctx, req); err != nil {
		log.Printf("request: method: %s, err: %v\n", info.FullMethod, err)
	}

	if v, ok := resp.(proto.Message); ok {
		str, _ = protojson.Marshal(v)
		log.Printf("response: method: %s, resp: %s\n", info.FullMethod, str)
	}

	return resp, err
}
