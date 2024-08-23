package request

type TermResize struct {
	Rows int `json:"rows" query:"rows" form:"rows"`
	Cols int `json:"cols" query:"cols" form:"cols"`
}

func (t TermResize) IsZero() bool {
	return t.Rows <= 0 && t.Cols <= 0
}

type TermSSH struct {
	TermResize
	Bastion  string `json:"bastion"  query:"bastion"  form:"bastion"  validate:"required"` // 堡垒机地址
	Username string `json:"username" query:"username" form:"username" validate:"required"` // 堡垒机用户名
	Password string `json:"password" query:"password" form:"password"`                     // 堡垒机密码
}
