package dynsql

type Table struct {
	wheres    Wheres
	orders    Orders
	whereMaps map[string]*Where
	orderMaps map[string]*Order
}

func (t Table) RawColumns() (Wheres, Orders) {
	return t.wheres, t.orders
}
