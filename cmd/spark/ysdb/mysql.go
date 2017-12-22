package ysdb

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/go-sql-driver/mysql"

	"hia/cmd/spark/types"
)

type DbMysql struct {
	db *sql.DB
}

// create tables
func (d *DbMysql) createAllTables() error {
	createUserTable := `
        create table if not exists user (
            register_time timestamp not null default current_timestamp,
            user_id bigint primary key,
            password varchar(128) not null,
            user_type ENUM('common','author') not null,
            name varchar(128) not null unique,
            user_id_card varchar(64),
            email varchar(64),
            account varchar(128),
            secret_key varchar(128),
            ethr_addr varchar(128),
            abi varchar(4096),
            last_update_time timestamp not null default current_timestamp
                                       on update current_timestamp
        );
    `

	createVideoTable := `
        create table if not exists video (
            upload_time timestamp not null default current_timestamp,
            url varchar(256) primary key,
            video_name varchar(128),
            user_id bigint not null,
            transaction varchar(64) not null,
            status boolean default true,
            plays int not null,
            buys int not null,
            constraint user_id foreign key(user_id) references user(user_id)
        );
    `
	createVideoTransactionTable := `
        create table if not exists video_transaction (
            buy_time timestamp not null default current_timestamp,
            transaction_id varchar(100) not null,
            url varchar(100) not null,
            user_id bigint not null,
            transaction varchar(100) not null,
            constraint url foreign key(url) references video(url),
            constraint buy_uid foreign key(user_id) references user(user_id)
        );
    `
	_, err := d.db.Exec(createUserTable)
	if nil != err {
		return err
	}

	_, err = d.db.Exec(createVideoTable)
	if nil != err {
		return err
	}

	_, err = d.db.Exec(createVideoTransactionTable)
	if nil != err {
		return err
	}

	return nil
}

// all user operation
func (d *DbMysql) UserAdd(user *types.User) error {
	if nil == user {
		return nil
	}

	format := " , %s = '%v'"
	s := fmt.Sprintf("insert user set user_id = %d, password = '%s', user_type = '%s'",
		user.UserID, user.Password, user.UserType)

	if "" != user.UserName {
		s += fmt.Sprintf(format, "name", user.UserName)
	}

	if "" != user.UserIdCard {
		s += fmt.Sprintf(format, "user_id_card", user.UserIdCard)
	}

	if "" != user.Email {
		s += fmt.Sprintf(format, "email", user.Email)
	}

	if "" != user.EthAccount {
		s += fmt.Sprintf(format, "account", user.EthAccount)
	}

	if "" != user.EthKey {
		s += fmt.Sprintf(format, "secret_key", user.EthKey)
	}

	if "" != user.EthKeyFileName {
		///////////
	}

	if "" != user.EthContractAddr {
		s += fmt.Sprintf(format, "ethr_addr", user.EthContractAddr)
	}

	if "" != user.EthAbi {
		s += fmt.Sprintf(format, "abi", user.EthAbi)
	}

	s += ";"

	fmt.Println(s)

	_, err := d.db.Exec(s)
	return err
}

func (d *DbMysql) UserUpdate(user *types.User) error {

	format := "%s %s = \"%v\""
	delimiter := " "

	s := "update user set "

	if "" != user.Password {
		s += fmt.Sprintf(format, delimiter, "password", user.Password)
		delimiter = ","
	}

	if "" != user.UserType {
		s += fmt.Sprintf(format, delimiter, "type", user.UserType)
		delimiter = ","
	}

	if "" != user.UserName {
		s += fmt.Sprintf(format, delimiter, "name", user.UserName)
		delimiter = ","
	}

	if "" != user.UserIdCard {
		s += fmt.Sprintf(format, delimiter, "user_id_card", user.UserIdCard)
		delimiter = ","
	}

	if "" != user.Email {
		s += fmt.Sprintf(format, delimiter, "email", user.Email)
		delimiter = ","
	}

	if "" != user.EthAccount {
		s += fmt.Sprintf(format, delimiter, "account", user.EthAccount)
		delimiter = ","
	}

	if "" != user.EthKey {
		s += fmt.Sprintf(format, delimiter, "secret_key", user.EthKey)
		delimiter = ","
	}

	if "" != user.EthContractAddr {
		s += fmt.Sprintf(format, delimiter, "ethr_addr", user.EthContractAddr)
		delimiter = ","
	}

	if "" != user.EthAbi {
		s += fmt.Sprintf(format, delimiter, "abi", user.EthAbi)
		delimiter = ","
	}

	s += fmt.Sprintf(" where user_id = %d;", user.UserID)

	if " " == delimiter {
		return nil
	}

	_, err := d.db.Exec(s)
	return err
}

func (d *DbMysql) UserQuery(user *types.User, sqls string) (*[]types.User, error) {
	format := " and %s = '%v' "

	s := `
        select 
            register_time, user_id, password, user_type, name, user_id_card,
            email, account, secret_key, ethr_addr, abi, last_update_time
        from
            user
        where 
            1 = 1 
    `
	// if 0 != user.RegisterTime {
	// 	s += fmt.Sprintf(" and %s = %v ", "register_time", user.RegisterTime)
	// }

	if 0 != user.UserID {
		s += fmt.Sprintf(" and %s = %v ", "user_id", user.UserID)
	}

	if "" != user.Password {
		s += fmt.Sprintf(format, "password", user.Password)
	}

	if "" != user.UserType {
		s += fmt.Sprintf(format, "user_type", user.UserType)
	}

	if "" != user.UserName {
		s += fmt.Sprintf(format, "name", user.UserName)
	}

	if "" != user.UserIdCard {
		s += fmt.Sprintf(format, "user_id_card", user.UserIdCard)
	}

	if "" != user.Email {
		s += fmt.Sprintf(format, "email", user.Email)
	}

	if "" != user.EthAccount {
		s += fmt.Sprintf(format, "account", user.EthAccount)
	}

	if "" != user.EthKey {
		s += fmt.Sprintf(format, "secret_key", user.EthKey)
	}

	if "" != user.EthContractAddr {
		s += fmt.Sprintf(format, "ethr_addr", user.EthContractAddr)
	}

	if "" != user.EthAbi {
		s += fmt.Sprintf(format, "abi", user.EthAbi)
	}

	if "" != sqls {
		s += " and " + sqls
	}

	s += ";"

	rows, err := d.db.Query(s)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	res := make([]types.User, 0, 10)

	for rows.Next() {
		var temp types.User
		var registerTime, lastUpdateTime mysql.NullTime
		var userIdCard, email, ethAccount, ethKey, ethContractAddr, ethAbi sql.NullString
		err = rows.Scan(&registerTime, &temp.UserID, &temp.Password, &temp.UserType, &temp.UserName,
			&userIdCard, &email, &ethAccount, &ethKey, &ethContractAddr, &ethAbi, &lastUpdateTime)
		if err != nil {
			return nil, err
		}

		temp.RegisterTime = registerTime.Time
		temp.LastUpdateTime = lastUpdateTime.Time

		if userIdCard.Valid {
			temp.UserIdCard = userIdCard.String
		}

		if email.Valid {
			temp.Email = email.String
		}

		if ethAccount.Valid {
			temp.EthAccount = ethAccount.String
		}

		if ethKey.Valid {
			temp.EthKey = ethKey.String
		}

		if ethContractAddr.Valid {
			temp.EthContractAddr = ethContractAddr.String
		}

		if ethAbi.Valid {
			temp.EthAbi = ethAbi.String
		}

		res = append(res, temp)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &res, nil
}

func (d *DbMysql) UserQuerySimple(user *types.User) (types.User, error) {
	searchUser, err := d.UserQuery(user, "")

	return (*searchUser)[0], err
}

func (d *DbMysql) UserDelete(user *types.User) error {
	if nil == user {
		return nil
	}

	s := fmt.Sprintf("delete from user where user_id = %v;", user.UserID)
	_, err := d.db.Exec(s)
	return err
}

//all video operation
func (d *DbMysql) VideoAdd(video *types.Video) error {
	if nil == video {
		return nil
	}

	format := ", %s = \"%v\" "
	s := fmt.Sprintf("insert video set url = \"%s\", user_id = %d, transaction = \"%s\", status = %v, plays = %d, buys = %d",
		video.URL, video.UserID, video.Transaction, video.Status, video.Plays, video.Buys)

	if "" != video.VideoName {
		s += fmt.Sprintf(format, "video_name", video.VideoName)
	}

	s += ";"
	_, err := d.db.Exec(s)

	return err
}

func (d *DbMysql) VideoUpdate(video *types.Video) error {
	if nil == video {
		return nil
	}

	format := " %s %s = \"%v\" "
	delimiter := " "

	s := fmt.Sprintf("update video set ")

	if "" != video.VideoName {
		s += fmt.Sprintf(format, delimiter, "video_name", video.VideoName)
		delimiter = ","
	}

	if 0 != video.UserID {
		s += fmt.Sprintf(" %s %s = %v ", delimiter, "user_id", video.UserID)
		delimiter = ","
	}

	if "" != video.Transaction {
		s += fmt.Sprintf(format, delimiter, "transaction", video.Transaction)
		delimiter = ","
	}

	if nil != video.Status {
		s += fmt.Sprintf(" %s %s = %v ", delimiter, "status", video.Status)
		delimiter = ","
	}

	if nil != video.Plays {
		s += fmt.Sprintf(" %s %s = %v ", delimiter, "plays", video.Plays)
		delimiter = ","
	}
	if nil != video.Buys {
		s += fmt.Sprintf(" %s %s = %v ", delimiter, "buys", video.Buys)
		delimiter = ","
	}

	s += fmt.Sprintf(" where url = \"%s\";", video.URL)

	if " " == delimiter {
		return nil
	}

	_, err := d.db.Exec(s)
	return err
}

func (d *DbMysql) VideoQuery(video *types.Video, sqls string) (*[]types.Video, error) {
	if nil == video {
		return nil, nil
	}

	format := " and %s = \"%v\" "
	format2 := " and %s = %v "
	s := `
        select 
            upload_time, url, video_name, user_id, transaction, status, plays, buys 
        from
            video 
        where 
            1 = 1 
    `

	if "" != video.URL {
		s += fmt.Sprintf(format, "url", video.URL)
	}

	if 0 != video.UploadTime {
		s += fmt.Sprintf(format2, "upload_time", video.UploadTime)
	}

	if "" != video.VideoName {
		s += fmt.Sprintf(format, "video_name", video.VideoName)
	}

	if 0 != video.UserID {
		s += fmt.Sprintf(format2, "user_id", video.UserID)
	}

	if "" != video.Transaction {
		s += fmt.Sprintf(format, "transaction", video.Transaction)
	}

	if nil != video.Status {
		s += fmt.Sprintf(format2, "status", video.Status)
	}

	if nil != video.Plays {
		s += fmt.Sprintf(format2, "plays", video.Plays)
	}

	if nil != video.Buys {
		s += fmt.Sprintf(format2, "buys", video.Buys)
	}

	if "" != sqls {
		s += " and " + sqls
	}

	s += ";"

	rows, err := d.db.Query(s)
	res := make([]types.Video, 0, 10)
	for rows.Next() {
		var temp types.Video
		err = rows.Scan(&temp.UploadTime, &temp.URL, &temp.VideoName, &temp.UserID,
			&temp.Transaction, &temp.Status, &temp.Plays, &temp.Buys)

		if err != nil {
			return nil, err
		}

		res = append(res, temp)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	if len(res) == 0 {
		return nil, nil
	}

	return &res, nil
}

func (d *DbMysql) VideoQuerySimple(video *types.Video) (types.Video, error) {
	fmt.Println("VideoQuerySimple")
	searchVideo, err := d.VideoQuery(video, "")

	return (*searchVideo)[0], err
}

func (d *DbMysql) VideoDelete(video *types.Video) error {
	if nil == video {
		return nil
	}

	s := fmt.Sprintf("delete from video where url = \"%s\";", video.URL)

	_, err := d.db.Exec(s)

	return err
}

//all video_transaction operation
func (d *DbMysql) VideoTransactionAdd(vt *types.VideoTransaction) error {
	if nil == vt {
		return nil
	}

	s := fmt.Sprintf("insert video_transaction set transaction_id = \"%s\", url = \"%s\", transaction = \"%s\", user_id= %d ;",
		vt.TransactionId, vt.URL, vt.Transaction, vt.UserID)

	_, err := d.db.Exec(s)

	return err
}

func (d *DbMysql) VideoTransactionQuery(vt *types.VideoTransaction, sqls string) (*[]types.VideoTransaction, error) {
	if nil == vt {
		return nil, nil
	}

	s := `
        select 
            buy_time, 
            transaction_id, 
            url, 
            user_id, 
            transaction 
        from 
            video_transaction 
        where 
            1 = 1
    `
	format := " and %s = \"%v\" "
	format2 := " and %s = %v "

	if "" != vt.URL {
		s += fmt.Sprintf(format, "url", vt.URL)
	}

	if 0 != vt.BuyTime {
		s += fmt.Sprintf(format2, "buy_time", vt.BuyTime)
	}

	if "" != vt.TransactionId {
		s += fmt.Sprintf(format, "transaction_id", vt.TransactionId)
	}

	if 0 != vt.UserID {
		s += fmt.Sprintf(format2, "user_id", vt.UserID)
	}

	if "" != vt.Transaction {
		s += fmt.Sprintf(format, "transaction", vt.Transaction)
	}

	if "" != sqls {
		s += " and " + sqls
	}

	s += ";"

	res := make([]types.VideoTransaction, 0, 10)

	rows, err := d.db.Query(s)
	if rows.Next() {
		var temp types.VideoTransaction

		err = rows.Scan(&temp.BuyTime, &temp.TransactionId, &temp.URL, &temp.UserID, &temp.Transaction)

		if err != nil {
			return nil, nil
		}

		res = append(res, temp)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	if len(res) == 0 {
		return nil, nil
	}
	return &res, nil
}

//////////////////////////////
var (
	dbMysql DbMysql
)

func init() {
	var err error
	dbMysql.db, err = sql.Open("mysql", "root:root@tcp(192.168.31.19)/test")
	if err != nil {
		log.Fatalln(err)
	}

	err = dbMysql.db.Ping()
	if err != nil {
		log.Fatalln(err)
	}

	err = dbMysql.createAllTables()
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("init db successful!")
}

func NewDbMysql() *DbMysql {
	return &dbMysql
}

func Close() {
	dbMysql.db.Close()
}
