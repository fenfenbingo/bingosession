package session

//session加载错误
//error info when loading session
type ErrInfo struct {
	Code    int
	Message string
}


var (
	NoErr = addErr(0, "")
	ErrEmpty = addErr(100, "ErrEmptySessionId")//session id为空
	ErrNoCache = addErr(101, "ErrNoCache(Expire)")//缓存未找到（已过期）
	ErrLoadCache = addErr(102, "ErrLoadCache(Expire)")//读取缓存出现异常（已过期）
	ErrExpire = addErr(103, "ErrExpire")//已过期
	ErrInValidSessId = addErr(800, "ErrInValidSessId")//session id不合法
	ErrCacheFormat = addErr(801, "ErrInValidSessId")//缓存数据格式异常
	ErrInvalidSign = addErr(802, "ErrInvalidSign")//session签名错误
)

func addErr(code int, msg string) *ErrInfo {
	return &ErrInfo{code, msg}
}
