package handler

import (
	"context"
	"fmt"

	api "github.com/fun-dotto/api-template/generated"
)

func (h *Handler) FCMTokenV1Upsert(ctx context.Context, request api.FCMTokenV1UpsertRequestObject) (api.FCMTokenV1UpsertResponseObject, error) {
	_ = ctx
	_ = request

	// TODO: FCMトークン作成/更新のユースケースを実装する。
	return nil, fmt.Errorf("not implemented")
}
