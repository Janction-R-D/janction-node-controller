package utils

import (
	"github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
)

const (
	duplicateEntryErrCode = 1062
)

// IsDuplicateEntryErr 唯一索引引起的插入失败
func IsDuplicateEntryErr(err error) bool {
	var sqlError *mysql.MySQLError
	ok := errors.As(err, &sqlError)
	if !ok {
		return false
	}
	return sqlError.Number == duplicateEntryErrCode
}
