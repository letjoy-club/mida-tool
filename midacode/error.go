package midacode

import (
	"errors"
	"fmt"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

var ErrorStrMap = map[string]string{}

var (
	LogLevelFatal = "fatal"
	LogLevelError = "error"
	LogLevelWarn  = "warn"
	LogLevelInfo  = "info"
)

type Error2 struct {
	cn      string
	message string
	level   string

	extra interface{}
}

func (e Error2) ToExtensions() map[string]interface{} {
	m := map[string]interface{}{}
	m["cn"] = e.cn
	if e.extra != nil {
		m["extra"] = e.extra
	}
	return m
}

func (e Error2) WithExtra(extra interface{}) Error2 {
	e.extra = extra
	return e
}

func (e Error2) CN() string {
	return e.cn
}

func (e Error2) Error() string {
	return e.message
}

func (e Error2) LogLevel() string {
	return e.level
}

func NewError(code string, name string, level string) error {
	ErrorStrMap[code] = name
	return Error2{
		cn:      name,
		message: code,
		level:   level,
	}
}

var (
	ErrClientTokenExpired  = NewError("CLIENT_TOKEN_EXPIRED", "用户登录已过期", LogLevelWarn)
	ErrClientTokenInvalid  = NewError("CLIENT_TOKEN_INVALID", "用户 token 不正确", LogLevelWarn)
	ErrNotPermitted        = NewError("NOT_PERMITTED", "没有权限", LogLevelWarn)
	ErrStateMayHaveChanged = NewError("STATE_MAY_CHANGED", "当前状态已发生改变，请重新进入页面", LogLevelError)
	ErrInternalError       = NewError("INTERNAL_ERROR", "内部错误", LogLevelError)
	ErrUnknownError        = NewError("UNKNOWN_ERROR", "发生了未知错误，你可以重试或者联系客服", LogLevelError)
	ErrResourceBusy        = NewError("RESOURCE_BUSY", "资源繁忙，请稍后再试", LogLevelWarn)

	ErrRecordExists = NewError("RECORD_EXISTS", "数据已存在", LogLevelWarn)
	ErrItemNotFound = NewError("ITEM_NOT_FOUND", "没有找到所需的数据", LogLevelWarn)
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
