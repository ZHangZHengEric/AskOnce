package dto_user

type LoginRes struct {
	UserId          int64  `json:"userId"`
	Account         string `json:"account"`
	Phone           string `json:"phone"`
	AtomechoSession string `json:"atomechoSession"`
}

type LoginInfoRes struct {
	UserId  int64  `json:"userId"`
	Account string `json:"account"`
	Phone   string `json:"phone"`
}
