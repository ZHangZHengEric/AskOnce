package dto_user

type LoginRes struct {
	UserId          string `json:"userId"`
	Account         string `json:"account"`
	AtomechoSession string `json:"atomechoSession"`
}

type LoginInfoRes struct {
	UserId  string `json:"userId"`
	Account string `json:"account"`
}
