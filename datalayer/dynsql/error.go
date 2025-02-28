package dynsql

import "fmt"

type OpError struct {
	Operator Operator
}

func (e OpError) Error() string {
	return fmt.Sprintf("dynsql operation error: %s", e.Operator)
}
