package handler

import (
	"context"
	"fmt"

	api "github.com/fun-dotto/api-template/generated"
)

func (h *Handler) FCMTokenV1Upsert(_ context.Context, _ api.FCMTokenV1UpsertRequestObject) (api.FCMTokenV1UpsertResponseObject, error) {
	// TODO: FCMトークン作成/更新のユースケースを実装する。
	return nil, fmt.Errorf("not implemented")
}
