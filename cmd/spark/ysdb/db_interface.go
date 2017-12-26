package ysdb

import (
	"time"

	"hia/cmd/spark/types"
)

type YsDb interface {
	Init(driverName, dataSourceName string) error

	UserAdd(user *types.User) error
	UserUpdate(user *types.User) error
	UserQuery(user *types.User, sqls string) (*[]types.User, error)
	UserQueryBetween(user *types.User, start time.Time, end time.Time) (*[]types.User, error)
	UserQueryAfter(user *types.User, time time.Time) (*[]types.User, error)
	UserQueryBefore(user *types.User, time time.Time) (*[]types.User, error)
	UserQuerySimple(user *types.User) (types.User, error)
	UserDelete(user *types.User) error

	VideoAdd(video *types.Video) error
	VideoUpdate(video *types.Video) error
	VideoQuery(video *types.Video, sqls string) (*[]types.Video, error)
	VideoQueryBetween(video *types.Video, start time.Time, end time.Time) (*[]types.Video, error)
	VideoQueryBefore(video *types.Video, end time.Time) (*[]types.Video, error)
	VideoQueryAfter(video *types.Video, start time.Time) (*[]types.Video, error)
	VideoQuerySimple(user *types.Video) (types.Video, error)
	VideoDelete(video *types.Video) error

	VideoTransactionAdd(vt *types.VideoTransaction) error
	VideoTransactionQuery(vt *types.VideoTransaction, sqls string) (*[]types.VideoTransaction, error)
	VideoTransactionQueryBetween(vt *types.VideoTransaction, start time.Time, end time.Time) (*[]types.VideoTransaction, error)
	VideoTransactionQueryAfter(vt *types.VideoTransaction, time time.Time) (*[]types.VideoTransaction, error)
	VideoTransactionQueryBefore(vt *types.VideoTransaction, time time.Time) (*[]types.VideoTransaction, error)
}
