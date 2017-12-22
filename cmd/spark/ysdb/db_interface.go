package ysdb

import (
	"hia/cmd/spark/types"
)

type YsDb interface {
	UserAdd(user *types.User) error
	UserUpdate(user *types.User) error
	UserQuery(user *types.User, sqls string) (*[]types.User, error)
	UserQuerySimple(user *types.User) (types.User, error)
	UserDelete(user *types.User) error

	VideoAdd(video *types.Video) error
	VideoUpdate(video *types.Video) error
	VideoQuery(video *types.Video, sqls string) (*[]types.Video, error)
	VideoQuerySimple(user *types.Video) (types.Video, error)
	VideoDelete(video *types.Video) error

	VideoTransactionAdd(vt *types.VideoTransaction) error
	VideoTransactionQuery(vt *types.VideoTransaction, sqls string) (*[]types.VideoTransaction, error)
}
