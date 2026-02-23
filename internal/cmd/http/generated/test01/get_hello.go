package test01

import (
	"context"

	"github.com/pro1ms/framework/internal/cmd/http/data/test01"
)

type GetHello struct {
}

func NewGetHello() *GetHello {
	return &GetHello{}
}

func (h *GetHello) GetHello(ctx context.Context, request test01.GetHelloRequestObject) (test01.GetHelloResponseObject, error) {
	return nil, nil
}
