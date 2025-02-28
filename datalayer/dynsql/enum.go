package dynsql

type Enums []*Enum

type Enum struct {
	Name  string `json:"name"`  // 备注说明
	Value any    `json:"value"` // 枚举值
}
