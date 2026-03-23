package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"slices"
	"sync"
	"syscall"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	inventory_v1 "github.com/ArchibaldKronin/microservices_test/shared/pkg/proto/inventory/v1"
	"google.golang.org/grpc"
)

type Category int

const (
	CategoryUnknown Category = iota
	CategoryEngine
	CategoryFuel
	CategoryPorthole
	CategoryWing
)

func CategoryToProto(c Category) inventory_v1.Category {
	switch c {
	case 1:
		return inventory_v1.Category_CATEGORY_ENGINE
	case 2:
		return inventory_v1.Category_CATEGORY_FUEL
	case 3:
		return inventory_v1.Category_CATEGORY_PORTHOLE
	case 4:
		return inventory_v1.Category_CATEGORY_WING
	default:
		return inventory_v1.Category_CATEGORY_UNSPECIFIED
	}
}

func CategoryToDomain(c inventory_v1.Category) Category {
	switch c {
	case 1:
		return CategoryEngine
	case 2:
		return CategoryFuel
	case 3:
		return CategoryPorthole
	case 4:
		return CategoryWing
	default:
		return CategoryUnknown
	}
}

type Dimensions struct {
	Length float64
	Width  float64
	Height float64
	Weight float64
}

func DimensionsToProto(d Dimensions) inventory_v1.Dimensions {
	return inventory_v1.Dimensions{
		Length: d.Length,
		Width:  d.Width,
		Height: d.Height,
		Weight: d.Weight,
	}
}

type Manufacturer struct {
	Name    string
	Country string
	Website string
}

func ManufactererToProto(m Manufacturer) inventory_v1.Manufacturer {
	return inventory_v1.Manufacturer{
		Name:    m.Name,
		Country: m.Country,
		Website: m.Website,
	}
}

type Value interface {
	isValue()
}

type StringValue struct {
	Value string
}

func (StringValue) isValue() {}

type Int64Value struct {
	Value int64
}

func (Int64Value) isValue() {}

type DoubleValue struct {
	Value float64
}

func (DoubleValue) isValue() {}

type BoolValue struct {
	Value bool
}

func (BoolValue) isValue() {}

func ValueToProto(val Value) *inventory_v1.Value {
	switch v := val.(type) {
	case StringValue:
		return &inventory_v1.Value{
			Value: &inventory_v1.Value_StringValue{
				StringValue: v.Value,
			},
		}
	case Int64Value:
		return &inventory_v1.Value{
			Value: &inventory_v1.Value_Int64Value{
				Int64Value: v.Value,
			},
		}
	case DoubleValue:
		return &inventory_v1.Value{
			Value: &inventory_v1.Value_DoubleValue{
				DoubleValue: v.Value,
			},
		}
	case BoolValue:
		return &inventory_v1.Value{
			Value: &inventory_v1.Value_BoolValue{
				BoolValue: v.Value,
			},
		}
	default:
		return nil
	}
}

func TimeToProto(t time.Time) *timestamppb.Timestamp {
	return timestamppb.New(t)
}

type Part struct {
	Uuid          string
	Name          string
	Description   string
	Price         float64
	StockQuantity int
	Category      Category
	Dimensions    Dimensions
	Manufacturer  Manufacturer
	Tags          []string
	Metadata      map[string]Value
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type PartStorage struct {
	mu    sync.RWMutex
	Parts map[string]*Part
}

func NewPartStorage() *PartStorage {
	return &PartStorage{
		Parts: make(map[string]*Part),
	}
}

const grpcPort = 50051

type PartService struct {
	inventory_v1.UnimplementedPartServiceServer

	Storage *PartStorage
}

func PartToProto(p *Part) *inventory_v1.Part {
	categoryProto := CategoryToProto(p.Category)
	dimensionsProto := DimensionsToProto(p.Dimensions)
	manufacturerProto := ManufactererToProto(p.Manufacturer)

	metadataProto := make(map[string]*inventory_v1.Value)
	for k, v := range p.Metadata {
		metadataProto[k] = ValueToProto(v)
	}

	return &inventory_v1.Part{
		Uuid:          p.Uuid,
		Name:          p.Name,
		Description:   p.Description,
		Price:         p.Price,
		StockQuantity: int64(p.StockQuantity),
		Category:      categoryProto,
		Dimensions:    &dimensionsProto,
		Manufacturer:  &manufacturerProto,
		Tags:          p.Tags,
		Metadata:      metadataProto,
		CreatedAt:     TimeToProto(p.CreatedAt),
		UpdatedAt:     TimeToProto(p.UpdatedAt),
	}
}

func (s *PartService) GetPart(_ context.Context, req *inventory_v1.GetPartRequest) (*inventory_v1.GetPartResponse, error) {
	s.Storage.mu.RLock()
	defer s.Storage.mu.RUnlock()

	part, ok := s.Storage.Parts[req.GetUuid()]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "запчасть с id: %s не найдена", req.GetUuid())
	}

	return &inventory_v1.GetPartResponse{
		Part: PartToProto(part),
	}, nil
}

func (s *PartService) ListPart(_ context.Context, req *inventory_v1.ListPartRequest) (*inventory_v1.ListPartResponse, error) {
	s.Storage.mu.RLock()
	defer s.Storage.mu.RUnlock()

	var parts []*inventory_v1.Part

outer:
	for _, part := range s.Storage.Parts {

		if len(req.Uuids) != 0 {
			if !slices.Contains(req.Uuids, part.Uuid) {
				continue outer
			}
		}

		if len(req.Names) != 0 {
			if !slices.Contains(req.Names, part.Name) {
				continue outer
			}
		}

		if len(req.Categorys) != 0 {
			categoryMatched := false
			for _, cat := range req.Categorys {
				if part.Category == CategoryToDomain(cat) {
					categoryMatched = true
					break
				}
			}
			if !categoryMatched {
				continue outer
			}
		}

		if len(req.Countrys) != 0 {
			if !slices.Contains(req.Countrys, part.Manufacturer.Country) {
				continue outer
			}
		}

		if len(req.Tags) != 0 {
			tagMatched := false

		tagOuter:
			for _, reqTag := range req.Tags {

				for _, innerTag := range part.Tags {

					if innerTag == reqTag {
						tagMatched = true
						break tagOuter
					}
				}
			}
			if !tagMatched {
				continue outer
			}
		}
		parts = append(parts, PartToProto(part))
	}

	return &inventory_v1.ListPartResponse{
		Parts: parts,
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v\n", err)
		return
	}
	defer func() {
		if cerr := lis.Close(); cerr != nil {
			log.Printf("failed to close listenner: %v\n", cerr)
		}
	}()

	s := grpc.NewServer()

	storage := NewPartStorage()
	/////////////////////////////////////
	for _, part := range PartsTest {
		storage.Parts[part.Uuid] = part
	}
	/////////////////////////////////////

	service := &PartService{
		Storage: storage,
	}

	inventory_v1.RegisterPartServiceServer(s, service)

	reflection.Register(s)

	go func() {
		log.Printf("🚀 gRPC Inventory server listening on %d\n", grpcPort)
		err = s.Serve(lis)
		if err != nil {
			log.Printf("failed to serve Inventory: %v\n", err)
			return
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("🛑 Shutting down gRPC Inventory server...")
	s.GracefulStop()
	log.Println("✅ Inventory server stopped")
}
