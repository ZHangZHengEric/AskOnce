package dto_user

type LoginAccountReq struct {
	Account  string `json:"account" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterAccountReq struct {
	Account  string `json:"account" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginPhoneReq struct {
	Phone   string `json:"phone" binding:"required"`
	SmsCode string `json:"smsCode" binding:"required"`
}

type LoginSendSmsReq struct {
	Phone string `json:"phone"`
}
