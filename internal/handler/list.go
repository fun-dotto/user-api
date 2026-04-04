package handler

import (
	"context"

	api "github.com/fun-dotto/api-template/generated"
)

func (h *Handler) UsersV1List(ctx context.Context, request api.UsersV1ListRequestObject) (api.UsersV1ListResponseObject, error) {
	_ = ctx
	_ = request

	// TODO: ユーザー一覧取得のユースケースを実装する。
	return nil, errNotImplemented
}
