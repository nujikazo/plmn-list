package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/nujikazo/plmn-list/api/config"
	pb "github.com/nujikazo/plmn-list/api/proto/gen/go"
	"github.com/nujikazo/plmn-list/database"
	"github.com/nujikazo/plmn-list/general"

	"google.golang.org/grpc"
)

// server
type server struct {
	DB *database.Database
}

// ListPlmn
func (s *server) ListPlmn(ctx context.Context, in *pb.ListPlmnRequest) (*pb.ListPlmnsResponses, error) {
	var q = map[string]string{"mcc": in.GetMcc(),
		"mnc":          in.GetMnc(),
		"iso":          in.GetIso(),
		"country":      in.GetCountry(),
		"country_code": in.GetCountryCode(),
		"network":      in.GetNetwork(),
	}

	result, err := s.DB.GetPlmnList(q)
	if err != nil {
		return nil, err
	}

	var list []*pb.Plmn
	for _, v := range result {
		var plmn pb.Plmn

		plmn.Mcc = v.MCC
		plmn.Mnc = v.MNC
		plmn.Iso = v.ISO
		plmn.Country = v.Country
		plmn.CountryCode = v.CountryCode
		plmn.Network = v.Network

		list = append(list, &plmn)
	}
	return &pb.ListPlmnsResponses{
		Plmn: list,
	}, nil
}

func main() {
	generalConf := general.ReadGeneralConf(os.Getenv("GENERAL_CONF"))
	apiConf := config.ReadAPIConf(os.Getenv("API_CONF"))

	db, err := database.New(generalConf)
	if err != nil {
		log.Fatal(err)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", apiConf.ServerAddr, apiConf.ServerPort))
	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer()
	pb.RegisterPlmnServiceServer(s, &server{DB: db})
	if err := s.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
