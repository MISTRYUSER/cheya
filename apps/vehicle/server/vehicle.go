package server

import (
	"context"
	"errors"
	"fmt"
	"strings"

	vehiclev1 "github.com/xuewentao/cheya/api/vehicle/v1"
	"github.com/xuewentao/cheya/apps/vehicle/ent"
	"github.com/xuewentao/cheya/apps/vehicle/ent/vehicle"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// vehicleServer 是对 Service 接口的具体实现
type VehicleServer struct {
	vehiclev1.UnimplementedVehicleServiceServer
	client ent.Client //hold database client
}

// NewVehicleServer 是构造函数
// 接收 ent.client
func NewVehicleServer(client ent.Client) *VehicleServer {
	return &VehicleServer{
		client: client,
	}
}

// GetVehicle 实现.proto 中定义的 rpc GetVehicle
func (s *VehicleServer) GetVehicle(ctx context.Context, req *vehiclev1.GetVehicleRequest) (*vehiclev1.GetVehicleResponse, error) {
	//校验参数
	if req.VehicleId == "" {
		return nil, errors.New("vehicle_id is requied")
	}
	//数据库查询
	//SELECT * FROM vehicles WHERE vin = ? LIMIT 1
	v, err := s.client.Vehicle.Query().
		Where(vehicle.Vin(req.VehicleId)).
		Only(ctx)
	//错误处理
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, status.Errorf(codes.NotFound, "vehicle not found:%s ", req.VehicleId)
		}
		return nil, status.Errorf(codes.Internal, "database error %v", err)
	}
	return &vehiclev1.GetVehicleResponse{
		Vehicle: &vehiclev1.Vehicle{
			Id:           fmt.Sprintf("%d", v.ID),
			Vin:          v.Vin,
			LicensePlate: v.LicensePlate,
			Status:       mapStatusToProto(v.Status),
		},
	}, nil
}
func mapStatusToProto(s string) vehiclev1.VehicleStatus {
	switch strings.ToLower(s) {
	case "online":
		return vehiclev1.VehicleStatus_VEHICLE_STATUS_ONLINE
	case "offline":
		return vehiclev1.VehicleStatus_VEHICLE_STATUS_OFFLINE
	default:
		return vehiclev1.VehicleStatus_VEHICLE_STATUS_UNSPECIFIED
	}
}

func (s *VehicleServer) CreateVehicle(ctx context.Context, req *vehiclev1.CreateVehicleRequest) (*vehiclev1.CreateVehicleReponse, error) {
	//1.简单校验
	if req.Vin == "" || req.LicensePlate == "" {
		return nil, status.Errorf(codes.InvalidArgument, "vin or license are not require")
	}
	//2.使用 Ent 插入 db
	// SQL: INSERT INTO vehicles (vin, license_plate, status, ...) VALUES (...)
	v, err := s.client.Vehicle.Create().
		SetVin(req.Vin).
		SetLicensePlate(req.LicensePlate).
		SetStatus("offline").
		Save(ctx)
	//3.错误处理
	if err != nil {
		if ent.IsConstraintError(err) {
			return nil, status.Errorf(codes.AlreadyExists, "vehicle with vin %s is already exists", req.Vin)
		}
		return nil, status.Errorf(codes.Internal, "failed to create vehicle: %v", err)
	}
	//4.返回结果
	return &vehiclev1.CreateVehicleReponse{
		VehicleId: fmt.Sprintf("%d", v.ID),
	}, nil

}

func (s *VehicleServer) ListVehicles(ctx context.Context, req *vehiclev1.ListVehiclesRequest) (*vehiclev1.ListVehiclesResponse, error) {
	//1.分页处理
	page := req.Page
	if page < 1 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize < 1 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize
	//2.查询总数
	total, err := s.client.Vehicle.Query().Count(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to count: %v", err)
	}
	//3.查询列表
	vehicles, err := s.client.Vehicle.Query().
		Limit(int(pageSize)).
		Offset(int(offset)).
		Order(ent.Desc("created_at")).
		All(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list: %v", err)
	}

	//转换数据
	pbVehicles := make([]*vehiclev1.Vehicle, len(vehicles))
	for i, v := range vehicles {
		pbVehicles[i] = &vehiclev1.Vehicle{
			Id:           fmt.Sprintf("%d", v.ID),
			Vin:          v.Vin,
			LicensePlate: v.LicensePlate,
			Status:       mapStatusToProto(v.Status),
		}
	}
	return &vehiclev1.ListVehiclesResponse{
		Vehicles:   pbVehicles,
		TotalCount: int32(total),
	}, nil
}
