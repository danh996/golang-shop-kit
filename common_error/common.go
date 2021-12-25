package common_error

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

var (
	ErrInvalidKeySize    = grpc.Errorf(codes.InvalidArgument, "Invalid Key size")
	ErrExpiredToken      = grpc.Errorf(codes.DeadlineExceeded, "Expired token")
	ErrCanNotCreateToken = grpc.Errorf(codes.Unknown, "Can not create token")
	ErrInvalidToken      = grpc.Errorf(codes.InvalidArgument, "Invalid token")

	ErrCanNotHashPassword = grpc.Errorf(codes.Unknown, "Can not hash password")

	ErrCanNotExtractInfo     = grpc.Errorf(codes.Unknown, "Can not extract info")
	ErrCanNotMappingMetadata = grpc.Errorf(codes.Unknown, "Can not maping metadata")
)
