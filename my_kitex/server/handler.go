package main

import (
	"context"
	math_service "my_kitex/kitex_gen/math_service"
)

// MathImpl implements the last service interface defined in the IDL.
type MathImpl struct{}

// Add implements the MathImpl interface.
func (s *MathImpl) Add(ctx context.Context, req *math_service.AddRequest) (resp *math_service.AddResponse, err error) {
	// TODO: Your code here...
	resp = &math_service.AddResponse{Sum: req.Left + req.Right}
	return
}

// Sub implements the MathImpl interface.
func (s *MathImpl) Sub(ctx context.Context, req *math_service.SubRequest) (resp *math_service.SubResponse, err error) {
	// TODO: Your code here...
	resp = &math_service.SubResponse{Diff: req.Left - req.Right}
	return
}
