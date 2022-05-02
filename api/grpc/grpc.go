package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	ac "github.com/nujikazo/plmn-list/api/config"
	pb "github.com/nujikazo/plmn-list/api/proto/gen/go"
	"github.com/nujikazo/plmn-list/config"
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
	var query map[string]string
	query = make(map[string]string)

	if in != nil {
		if mcc := in.GetMcc(); mcc != "" {
			query[general.Mcc] = mcc
		}

		if mnc := in.GetMnc(); mnc != "" {
			query[general.Mnc] = mnc
		}

		if iso := in.GetIso(); iso != "" {
			query[general.Iso] = iso
		}

		if country := in.GetCountry(); country != "" {
			query[general.Country] = country
		}

		if countryCode := in.GetCountryCode(); countryCode != "" {
			query[general.CountryCode] = countryCode
		}

		if network := in.GetNetwork(); network != "" {
			query[general.Network] = network
		}
	}

	if err := s.DB.GetPlmnList(query); err != nil {
		return nil, err
	}

	list := s.toPbResponse()

	return &pb.ListPlmnsResponses{
		Plmn: list,
	}, nil
}

func (s *server) toPbResponse() []*pb.Plmn {
	var list = make([]*pb.Plmn, len(s.DB.Result))
	for _, v := range s.DB.Result {
		var plmn pb.Plmn

		plmn.Mcc = v.MCC
		plmn.Mnc = v.MNC
		plmn.Iso = v.ISO
		plmn.Country = v.Country
		plmn.CountryCode = v.CountryCode
		plmn.Network = v.Network

		list = append(list, &plmn)
	}

	return list
}

func main() {
	generalConf := config.New(os.Getenv("GENERAL_CONF"))
	apiConf := ac.New(os.Getenv("API_CONF"))

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
