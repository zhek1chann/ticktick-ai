package logger

type keyType int

const key = keyType(0)

type logCtx struct {
	ShopID  int64
	Message string
}
