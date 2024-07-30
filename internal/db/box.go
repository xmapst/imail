package db

import (
	"fmt"
	_ "github.com/midoks/imail/internal/tools"
	// "gorm.io/gorm"
	"strings"
	// "errors"
	// "time"
)

type Box struct {
	Id         int64 `gorm:"primaryKey"`
	Uid        int64 `gorm:"comment:UID"`
	Mid        int64 `gorm:"comment:MID"`
	Type       int   `gorm:"comment:TYPE|0:接收邮件;1:发送邮件"`
	Size       int   `gorm:"comment:邮件内容大小[byte]"`
	UpdateTime int64 `gorm:"autoCreateTime;comment:更新时间"`
	CreateTime int64 `gorm:"autoCreateTime;comment:创建时间"`
}

func (Box) TableName() string {
	return TablePrefix("user_box")
}

func BoxUserList(uid int64) (int64, int64) {

	var resultBox Box
	sql := fmt.Sprintf("SELECT count(uid) as count, sum(size) as size FROM `%s` WHERE uid=?", TablePrefix("user_box"))
	db.Raw(sql, uid).Find(&resultBox)

	// fmt.Println(uid, num, resultBox)
	return 0, 0
}

func BoxUserTotal(uid int64) (int64, int64) {

	var resultBox Box
	sql := fmt.Sprintf("SELECT count(uid) as count, sum(size) as size FROM `%s` WHERE uid=?", TablePrefix("user_box"))
	db.Raw(sql, uid).Find(&resultBox)

	// fmt.Println(num, resultBox)
	return 0, 0
}

// Get statistics under classification
func BoxUserMessageCountByClassName(uid int64, className string) (int64, int64) {
	type Result struct {
		Count int64
		Size  int64
	}
	var result Result

	sql := fmt.Sprintf("SELECT count(uid) as count, sum(size) as size FROM `%s` WHERE uid=?", TablePrefix("mail"))

	if strings.EqualFold(className, "Sent Messages") {
		sql = fmt.Sprintf("%s and type='%d' and is_delete='0' and is_draft='0' ", sql, 0)
	}

	if strings.EqualFold(className, "INBOX") {
		sql = fmt.Sprintf("%s and type='%d' and is_delete='0' and is_junk='0' and is_flags='0' and is_draft='0'", sql, 1)
	}

	if strings.EqualFold(className, "Deleted Messages") {
		sql = fmt.Sprintf("%s and is_delete='1' and is_draft='0'", sql)
	}

	if strings.EqualFold(className, "Drafts") {
		sql = fmt.Sprintf("%s and is_delete='0' and is_draft='1'", sql)
	}

	if strings.EqualFold(className, "Junk") {
		sql = fmt.Sprintf("%s and is_junk='1' and is_draft='0'", sql)
	}

	// fmt.Println("BoxUserMessageCountByClassName:", sql, className)
	num := db.Raw(sql, uid).Scan(&result)
	if num.Error != nil {
		return 0, 0
	}

	return result.Count, result.Size
}

// // Paging List of Imap Protocol
func BoxListByImap(uid int64, className string, start int64, end int64) ([]Mail, error) {
	var result []Mail

	sql := ""
	if end > 0 {
		sql = fmt.Sprintf("SELECT * FROM `%s` WHERE uid=? and id>='%d' and id<='%d'", TablePrefix("mail"), start, end)
	} else {
		sql = fmt.Sprintf("SELECT * FROM `%s` WHERE uid=? and id>='%d'", TablePrefix("mail"), start)
	}

	if strings.EqualFold(className, "Sent Messages") {
		sql = fmt.Sprintf("%s and type='%d' and is_delete='0' ", sql, 0)
	}

	if strings.EqualFold(className, "INBOX") {
		sql = fmt.Sprintf("%s and type='%d' and is_delete='0' and is_junk='0'", sql, 1)
	}

	if strings.EqualFold(className, "Deleted Messages") {
		sql = fmt.Sprintf("%s and is_delete='1' ", sql)
	}

	if strings.EqualFold(className, "Drafts") {
		sql = fmt.Sprintf("%s and is_draft='1' ", sql)
	}

	if strings.EqualFold(className, "Junk") {
		sql = fmt.Sprintf("%s and is_junk='1'", sql)
	}

	sql = fmt.Sprintf("%s limit 500", sql)

	// fmt.Println("BoxListByImap:", sql, className)
	db.Raw(sql, uid).Find(&result)
	return result, err
}

func BoxListByMid(uid int64, className string, mid int64) ([]Mail, error) {
	var result []Mail
	sql := fmt.Sprintf("SELECT * FROM `%s` WHERE uid=? and id='%d' limit 500", TablePrefix("mail"), mid)

	db.Raw(sql, uid).Find(&result)
	return result, err
}
