package enums


type JsonResultCode int

const (
	JRCodeSuccess JsonResultCode = iota
	JRCodeFail
	JRCode302 = 302  //跳转至地址
	JRCode401 = 401   //未授权访问
)

const(
	Deleted = iota-1
	Disabled
	Enabled
)

const CURRENT_USER = "currentuser"  // session 当前用户 key