package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"context"

	runtime "github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/nujikazo/plmn-list/api/config"
	gw "github.com/nujikazo/plmn-list/api/proto/gen/go"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	apiConf := config.New(os.Getenv("API_CONF"))

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	endpoint := fmt.Sprintf("%s:%d", apiConf.ServerAddr, apiConf.ServerPort)
	err := gw.RegisterPlmnServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	if err != nil {
		log.Fatal(err)
	}

	if err := http.ListenAndServe(fmt.Sprintf(":%d", apiConf.GatewayPort), mux); err != nil {
		log.Fatal(err)
	}
}
