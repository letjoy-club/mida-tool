package mitacode

import (
	"errors"
	"fmt"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

var ErrorStrMap = map[string]string{}

type Error2 struct {
	cn      string
	message string
	log     bool
}

func (e Error2) CN() string {
	return e.cn
}

func (e Error2) Error() string {
	return e.message
}

func (e Error2) ShouldBeLogged() bool {
	return e.log
}

func NewError(code string, name string, log bool) error {
	ErrorStrMap[code] = name
	return Error2{
		cn:      name,
		message: code,
		log:     log,
	}
}

var (
	ErrClientTokenExpired  = NewError("CLIENT_TOKEN_EXPIRED", "用户登录已过期", true)
	ErrClientTokenInvalid  = NewError("CLIENT_TOKEN_INVALD", "用户 token 不正确", true)
	ErrNotPermitted        = NewError("NOT_PERMITTED", "没有权限", true)
	ErrStateMayHaveChanged = NewError("STATE_MAY_CHANGED", "当前状态已发生改变，请重新进入页面", true)

	ErrRecordExists = NewError("RECORD_EXISTS", "数据已存在", false)
	ErrItemNotFound = NewError("ITEM_NOT_FOUND", "没有找到所需的数据", false)
)

var (
	dbError = "DB_ERROR"
)

func DBError(err error) error {
	return fmt.Errorf("%s: %s", dbError, err)
}

func ItemMayExist(err error) error {
	if err == nil {
		return nil
	}
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
		return ErrRecordExists
	}
	return err
}

func ItemMayNotFound(err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrItemNotFound
	}
	return err
}

func ItemIsNotFound(err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrItemNotFound
	}
	return err
}

func ItemCustomNotFound(err error, customized error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return customized
	}
	return err
}
