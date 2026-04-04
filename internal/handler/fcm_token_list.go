package handler

import (
	"context"

	api "github.com/fun-dotto/api-template/generated"
)

func (h *Handler) FCMTokenV1List(ctx context.Context, request api.FCMTokenV1ListRequestObject) (api.FCMTokenV1ListResponseObject, error) {
	_ = ctx
	_ = request

	// TODO: FCMトークン一覧取得のユースケースを実装する。
	return nil, errNotImplemented
}
