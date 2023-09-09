package grpc

import (
	"context"
	"errors"

	"github.com/quangdangfit/gocommon/logger"

	"goshop/internal/user/dto"
	"goshop/internal/user/service"
	"goshop/pkg/utils"
	pb "goshop/proto/gen/go/user"
)

type UserHandler struct {
	pb.UnimplementedUserServiceServer

	service service.IUserService
}

func NewUserHandler(service service.IUserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

func (h *UserHandler) Login(ctx context.Context, req *pb.LoginReq) (*pb.LoginRes, error) {
	user, accessToken, refreshToken, err := h.service.Login(ctx, &dto.LoginReq{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		logger.Error("Failed to register ", err)
		return nil, err
	}

	var res pb.LoginRes
	utils.Copy(&res.User, &user)
	res.AccessToken = accessToken
	res.RefreshToken = refreshToken
	return &res, nil
}

func (h *UserHandler) Register(ctx context.Context, req *pb.RegisterReq) (*pb.RegisterRes, error) {
	user, err := h.service.Register(ctx, &dto.RegisterReq{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		logger.Error("Failed to register ", err)
		return nil, err
	}

	var res pb.RegisterRes
	utils.Copy(&res.User, &user)
	return &res, nil
}

func (h *UserHandler) GetMe(ctx context.Context, _ *pb.GetMeReq) (*pb.GetMeRes, error) {
	userID, _ := ctx.Value("userId").(string)
	if userID == "" {
		return nil, errors.New("unauthorized")
	}

	user, err := h.service.GetUserByID(ctx, userID)
	if err != nil {
		logger.Error("Failed to register ", err)
		return nil, err
	}

	var res pb.GetMeRes
	utils.Copy(&res.User, &user)
	return &res, nil
}

func (h *UserHandler) RefreshToken(ctx context.Context, req *pb.RefreshTokenReq) (*pb.RefreshTokenRes, error) {
	userID, _ := ctx.Value("userId").(string)
	if userID == "" {
		return nil, errors.New("unauthorized")
	}

	accessToken, err := h.service.RefreshToken(ctx, userID)
	if err != nil {
		logger.Error("Failed to register ", err)
		return nil, err
	}

	res := pb.RefreshTokenRes{
		AccessToken: accessToken,
	}
	return &res, nil
}

func (h *UserHandler) ChangePassword(ctx context.Context, req *pb.ChangePasswordReq) (*pb.ChangePasswordRes, error) {
	userID, _ := ctx.Value("userId").(string)
	if userID == "" {
		return nil, errors.New("unauthorized")
	}

	err := h.service.ChangePassword(ctx, userID, &dto.ChangePasswordReq{
		Password:    req.Password,
		NewPassword: req.NewPassword,
	})
	if err != nil {
		logger.Error("Failed to register ", err)
		return nil, err
	}

	return &pb.ChangePasswordRes{}, nil
}
