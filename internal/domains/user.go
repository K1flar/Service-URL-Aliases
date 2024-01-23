package domains

type User struct {
	ID       uint32 `json:"id"`
	Login    string `json:"login"`
	Password string `json:"password"`
	Email    string `json:"emain"`
}
