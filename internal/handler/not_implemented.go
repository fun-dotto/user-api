package handler

import (
	"context"
	"errors"

	api "github.com/fun-dotto/api-template/generated"
)

var errNotImplemented = errors.New("not implemented")

func (h *Handler) FCMTokenV1List(ctx context.Context, request api.FCMTokenV1ListRequestObject) (api.FCMTokenV1ListResponseObject, error) {
	_ = ctx
	_ = request

	// TODO: FCMトークン一覧取得のユースケースを実装する。
	return nil, errNotImplemented
}

func (h *Handler) FCMTokenV1Upsert(ctx context.Context, request api.FCMTokenV1UpsertRequestObject) (api.FCMTokenV1UpsertResponseObject, error) {
	_ = ctx
	_ = request

	// TODO: FCMトークン作成/更新のユースケースを実装する。
	return nil, errNotImplemented
}

func (h *Handler) UsersV1List(ctx context.Context, request api.UsersV1ListRequestObject) (api.UsersV1ListResponseObject, error) {
	_ = ctx
	_ = request

	// TODO: ユーザー一覧取得のユースケースを実装する。
	return nil, errNotImplemented
}
