package ysdb

import (
	"fmt"
	"hia/cmd/spark/types"
)

type ABCDatabase struct {
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

func NewABCDatabase() (YsDb, error) {
	return &ABCDatabase{}, nil
}
