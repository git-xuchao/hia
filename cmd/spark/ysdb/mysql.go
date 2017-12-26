package ysdb

import (
	"database/sql"
	"fmt"
	"log"
	"time"

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
            video_id varchar(256) primary key,
            video_name varchar(256),
            url varchar(256) not null,
            user_id bigint not null,
            transaction varchar(64) not null,
            status boolean default true,
            plays int not null default 0,
            buys int not null default 0,
            constraint user_id foreign key(user_id) references user(user_id)
        );
    `
	createVideoTransactionTable := `
        create table if not exists video_transaction (
            buy_time timestamp not null default current_timestamp,
            transaction_id varchar(100) not null,
            video_id varchar(256) not null,
            user_id bigint not null,
            transaction varchar(100) not null,
            constraint  video_id foreign key(video_id) references video(video_id),
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

func (d *DbMysql) Init(driverName, dataSourceName string) error {
	var err error
	/*
	 *dbMysql.db, err = sql.Open("mysql", "root:root@tcp(192.168.31.19)/test")
	 */
	dbMysql.db, err = sql.Open(driverName, dataSourceName)
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

	return err
}

// all user operation
func (d *DbMysql) UserAdd(user *types.User) error {
	if nil == user {
		return nil
	}

	format := " , %s = '%v' "
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

	fmt.Println("sql", s)
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

func (d *DbMysql) UserQueryBetween(user *types.User, start time.Time, end time.Time) (*[]types.User, error) {
	fmt.Println("UserQueryBetween", "timeStart", start, "timeEnd", end)
	sqlStr := fmt.Sprintf(" %d <= register_time <= %d ", start.Unix(), end.Unix())
	return d.UserQuery(user, sqlStr)
}

func (d *DbMysql) UserQueryAfter(user *types.User, time time.Time) (*[]types.User, error) {
	fmt.Println("UserQueryAfter", "time", time)
	sqlStr := fmt.Sprintf(" %d <= register_time", time.Unix())
	return d.UserQuery(user, sqlStr)
}

func (d *DbMysql) UserQueryBefore(user *types.User, time time.Time) (*[]types.User, error) {
	fmt.Println("UserQueryBefore", "time", time)
	sqlStr := fmt.Sprintf(" register_time <= %d ", time.Unix())
	return d.UserQuery(user, sqlStr)
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
	s := fmt.Sprintf("insert video set video_id = \"%s\", url = \"%s\", user_id = %d, transaction = \"%s\"",
		video.VideoID, video.URL, video.UserID, video.Transaction)

	if nil != video.Status {
		s += fmt.Sprintf(", %s = %v ", video.Status)
	}

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
	format2 := " %s %s = %v "
	delimiter := " "

	s := fmt.Sprintf("update video set ")

	if "" != video.VideoName {
		s += fmt.Sprintf(format, delimiter, "video_name", video.VideoName)
		delimiter = ","
	}

	if "" != video.URL {
		s += fmt.Sprintf(format, delimiter, "url", video.URL)
		delimiter = ","
	}

	if 0 != video.UserID {
		s += fmt.Sprintf(format2, delimiter, "user_id", video.UserID)
		delimiter = ","
	}

	if "" != video.Transaction {
		s += fmt.Sprintf(format, delimiter, "transaction", video.Transaction)
		delimiter = ","
	}

	if nil != video.Status {
		s += fmt.Sprintf(format2, delimiter, "status", video.Status)
		delimiter = ","
	}

	if nil != video.Plays {
		s += fmt.Sprintf(format2, delimiter, "plays", video.Plays)
		delimiter = ","
	}
	if nil != video.Buys {
		s += fmt.Sprintf(format2, delimiter, "buys", video.Buys)
		delimiter = ","
	}

	s += fmt.Sprintf(" where video_id = \"%s\";", video.VideoID)

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
            upload_time, url, video_id, video_name, user_id, transaction, status, plays, buys 
        from
            video 
        where 
            1 = 1 
    `

	if "" != video.URL {
		s += fmt.Sprintf(format, "url", video.URL)
	}

	if "" != video.VideoID {
		s += fmt.Sprintf(format, "video_id", video.VideoID)
	}

	// if 0 != video.UploadTime {
	// 	s += fmt.Sprintf(format2, "upload_time", video.UploadTime)
	// }

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

	fmt.Println("sql", s)
	rows, err := d.db.Query(s)
	if err != nil {
		return nil, err
	}

	res := make([]types.Video, 0, 10)
	for rows.Next() {
		var temp types.Video
		var uploadTime mysql.NullTime
		var videoName sql.NullString

		temp.Status = new(bool)
		temp.Plays = new(uint)
		temp.Buys = new(uint)

		err = rows.Scan(&uploadTime, &temp.URL, &temp.VideoID, &videoName, &temp.UserID,
			&temp.Transaction, temp.Status, temp.Plays, temp.Buys)

		if err != nil {
			return nil, err
		}
		fmt.Println(*temp.Status, *temp.Plays, *temp.Buys)
		temp.UploadTime = uploadTime.Time

		if videoName.Valid {
			temp.VideoName = videoName.String
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

func (d *DbMysql) VideoQueryBetween(video *types.Video, start time.Time, end time.Time) (*[]types.Video, error) {
	fmt.Println("VideoQueryBetween", "timeStart", start, "timeEnd", end)
	sqlStr := fmt.Sprintf(" %d <= upload_time <= %d ", start.Unix(), end.Unix())
	return d.VideoQuery(video, sqlStr)
}

func (d *DbMysql) VideoQueryBefore(video *types.Video, end time.Time) (*[]types.Video, error) {
	fmt.Println("VideoQueryBefore", "timeEnd", end)
	sqlStr := fmt.Sprintf(" upload_time <= %d ", end.Unix())
	return d.VideoQuery(video, sqlStr)
}

func (d *DbMysql) VideoQueryAfter(video *types.Video, start time.Time) (*[]types.Video, error) {
	fmt.Println("VideoQueryAfter", "timeStart", start)
	sqlStr := fmt.Sprintf(" upload_time <= %d ", start.Unix())
	return d.VideoQuery(video, sqlStr)
}

func (d *DbMysql) VideoQuerySimple(video *types.Video) (types.Video, error) {
	fmt.Println("VideoQuerySimple")
	searchVideo, err := d.VideoQuery(video, "")

	if searchVideo != nil {
		return (*searchVideo)[0], err
	} else {
		return *video, err
	}

}

func (d *DbMysql) VideoDelete(video *types.Video) error {
	if nil == video {
		return nil
	}

	s := fmt.Sprintf("delete from video where video_id = \"%s\";", video.VideoID)

	fmt.Println("VideoDelete sql", s)

	_, err := d.db.Exec(s)

	return err
}

//all video_transaction operation
func (d *DbMysql) VideoTransactionAdd(vt *types.VideoTransaction) error {
	if nil == vt {
		return nil
	}

	s := fmt.Sprintf("insert video_transaction set transaction_id = \"%s\", video_id = \"%s\", transaction = \"%s\", user_id= %d ;",
		vt.TransactionId, vt.VideoID, vt.Transaction, vt.UserID)

	fmt.Println("VideoTransactionAdd sql", s)
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
            video_id, 
            user_id, 
            transaction 
        from 
            video_transaction 
        where 
            1 = 1
    `
	format := " and %s = \"%v\" "
	format2 := " and %s = %v "

	if "" != vt.VideoID {
		s += fmt.Sprintf(format, "video_id", vt.VideoID)
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

	fmt.Println("VideoTransactionQuery sql", s)
	rows, err := d.db.Query(s)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var temp types.VideoTransaction
		var buyTime mysql.NullTime

		err = rows.Scan(&buyTime, &temp.TransactionId, &temp.VideoID, &temp.UserID, &temp.Transaction)

		if err != nil {
			return nil, err
		}

		temp.BuyTime = buyTime.Time

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

func (d *DbMysql) VideoTransactionQueryBetween(transaction *types.VideoTransaction, start time.Time, end time.Time) (*[]types.VideoTransaction, error) {
	fmt.Println("VideoTransactionQueryBetween", "timeStart", start, "timeEnd", end)
	sqlStr := fmt.Sprintf(" %d <= buy_time <= %d ", start.Unix(), end.Unix())
	return d.VideoTransactionQuery(transaction, sqlStr)
}

func (d *DbMysql) VideoTransactionQueryAfter(transaction *types.VideoTransaction, time time.Time) (*[]types.VideoTransaction, error) {
	fmt.Println("VideoTransactionQueryAfter", "time", time)
	sqlStr := fmt.Sprintf(" %d <= buy_time", time.Unix())
	return d.VideoTransactionQuery(transaction, sqlStr)
}

func (d *DbMysql) VideoTransactionQueryBefore(transaction *types.VideoTransaction, time time.Time) (*[]types.VideoTransaction, error) {
	fmt.Println("VideoTransactionQueryBefore", "time", time)
	sqlStr := fmt.Sprintf(" buy_time <= %d ", time.Unix())
	return d.VideoTransactionQuery(transaction, sqlStr)
}

//////////////////////////////
var (
	dbMysql DbMysql
)

func NewDbMysql() *DbMysql {
	return &dbMysql
}

func Close() {
	dbMysql.db.Close()
}
