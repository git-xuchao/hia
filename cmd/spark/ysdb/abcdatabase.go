package ysdb

import (
	"fmt"
	"hia/cmd/spark/types"
	"time"
)

type ABCDatabase struct {
}

func (this *ABCDatabase) Init(driverName, dataSourceName string) error {
	fmt.Println("init")
	return nil
}

func (this *ABCDatabase) UserAdd(user *types.User) error {
	fmt.Println("UserAdd")
	return nil
}

func (this *ABCDatabase) UserUpdate(user *types.User) error {
	fmt.Println("UserUpdate")
	return nil
}

func (this *ABCDatabase) UserQuery(user *types.User, sqls string) (*[]types.User, error) {
	fmt.Println("UserQuery")
	return nil, nil
}

func (this *ABCDatabase) UserQueryBetween(user *types.User, start time.Time, end time.Time, page int, count int) (*[]types.User, error) {
	fmt.Println("UserQueryBetween")
	return nil, nil
}

func (this *ABCDatabase) UserQueryAfter(user *types.User, time time.Time, page int, count int) (*[]types.User, error) {
	fmt.Println("UserQueryAfter")
	return nil, nil
}

func (this *ABCDatabase) UserQueryBefore(user *types.User, time time.Time, page int, count int) (*[]types.User, error) {
	fmt.Println("UserQueryBefore")
	return nil, nil
}

func (this *ABCDatabase) UserQuerySimple(video *types.User) (types.User, error) {
	fmt.Println("UserQuerySimple")
	var newUser types.User
	return newUser, nil
}

func (this *ABCDatabase) UserDelete(user *types.User) error {
	fmt.Println("UserDelete")
	return nil
}

func (this *ABCDatabase) VideoAdd(video *types.Video) error {
	fmt.Println("VideoAdd")
	return nil
}

func (this *ABCDatabase) VideoUpdate(video *types.Video) error {
	fmt.Println("VideoUpdate")
	return nil
}

func (this *ABCDatabase) VideoQuery(video *types.Video, sqls string) (*[]types.Video, error) {
	fmt.Println("VideoQuery")
	return nil, nil
}

func (this *ABCDatabase) VideoQueryBetween(video *types.Video, start time.Time, end time.Time, page int, count int) (*[]types.Video, error) {
	fmt.Println("VideoQueryBetween")
	return nil, nil
}

func (this *ABCDatabase) VideoQueryBefore(video *types.Video, end time.Time, page int, count int) (*[]types.Video, error) {
	fmt.Println("VideoQueryBefore")
	return nil, nil
}

func (this *ABCDatabase) VideoQueryAfter(video *types.Video, start time.Time, page int, count int) (*[]types.Video, error) {
	fmt.Println("VideoQueryAfter")
	return nil, nil
}

func (this *ABCDatabase) VideoQuerySimple(video *types.Video) (types.Video, error) {
	fmt.Println("VideoQuerySimple")
	var newVideo types.Video
	return newVideo, nil
}

func (this *ABCDatabase) VideoDelete(video *types.Video) error {
	fmt.Println("VideoDelete")
	return nil
}

func (this *ABCDatabase) VideoTransactionAdd(vt *types.VideoTransaction) error {
	fmt.Println("VideoTransactionAdd")
	return nil
}

func (this *ABCDatabase) VideoTransactionQuery(vt *types.VideoTransaction, sqls string) (*[]types.VideoTransaction, error) {
	fmt.Println("VideoTransactionQuery")
	return nil, nil
}

func (this *ABCDatabase) VideoTransactionQueryBetween(vt *types.VideoTransaction, start time.Time, end time.Time, page int, count int) (*[]types.VideoTransaction, error) {
	fmt.Println("VideoTransactionQueryBetwee")
	return nil, nil
}

func (this *ABCDatabase) VideoTransactionQueryAfter(vt *types.VideoTransaction, time time.Time, page int, count int) (*[]types.VideoTransaction, error) {
	fmt.Println("VideoTransactionQueryAfter")
	return nil, nil
}

func (this *ABCDatabase) VideoTransactionQueryBefore(vt *types.VideoTransaction, time time.Time, page int, count int) (*[]types.VideoTransaction, error) {
	fmt.Println("VideoTransactionQueryBefore")
	return nil, nil
}

func NewABCDatabase() (YsDb, error) {
	return &ABCDatabase{}, nil
}
