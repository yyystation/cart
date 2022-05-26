package handler

import (
	"context"

	"github.com/yyystation/common"
	"go.micro.service.cart/domain/model"
	"go.micro.service.cart/domain/service"
	cart "go.micro.service.cart/proto"
)

type Cart struct {
	CartDataService service.ICartDataService
}

func (h *Cart) AddCart(ctx context.Context, request *cart.CartInfo, response *cart.ResponseAdd) (err error) {
	cart := &model.Cart{}
	common.SwapTo(request, cart)
	response.CartId, err = h.CartDataService.AddCart(cart)
	return err
}

func (h *Cart) CleanCart(ctx context.Context, request *cart.Clean, reponse *cart.Response) error {
	if err := h.CartDataService.CleanCart(request.UserId); err != nil {
		return err
	}
	reponse.Msg = "清空购物车成功"
	return nil
}

func (h *Cart) Incr(ctx context.Context, request *cart.Item, reponse *cart.Response) error {
	if err := h.CartDataService.IncrNum(request.Id, request.ChangeNum); err != nil {
		return err
	}
	reponse.Msg = "购物车添加成功"
	return nil
}

//购物车减少商品数量
func (h *Cart) Decr(ctx context.Context, request *cart.Item, reponse *cart.Response) error {
	if err := h.CartDataService.DecrNum(request.Id, request.ChangeNum); err != nil {
		return err
	}
	reponse.Msg = "购物车减少成功"
	return nil
}

func (h *Cart) DeleteItemByID(ctx context.Context, request *cart.CartID, reponse *cart.Response) error {
	if err := h.CartDataService.DeleteCart(request.Id); err != nil {
		return err
	}
	reponse.Msg = "购物车删除成功"
	return nil
}

//查询用户所有的购物车信息
func (h *Cart) GetAll(ctx context.Context, request *cart.CartFindAll, reponse *cart.CartAll) error {
	cartAll, err := h.CartDataService.FindAllCart(request.UserId)
	if err != nil {
		return err
	}
	for _, v := range cartAll {
		cart := &cart.CartInfo{}
		if err := common.SwapTo(v, cart); err != nil {
			return err
		}
		reponse.Cartinfo = append(reponse.Cartinfo, cart)
	}
	return nil
}
