package mongodsl_test

import (
	"encoding/json"
	"fmt"

	"security-gitlab-devel.eastmoney.com/security/siem-backend/mongodsl"
)

// mustJSON 将过滤器序列化为紧凑 JSON，便于在示例中做稳定的输出断言。
func mustJSON(v any) string {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return string(b)
}

// ExampleParser_Parse 演示最基本的等值查询。
func ExampleParser_Parse() {
	psr := mongodsl.NewParser()

	filter, err := psr.Parse(`status == "active"`)
	if err != nil {
		panic(err)
	}

	fmt.Println(mustJSON(filter))
	// Output: {"status":"active"}
}

// ExampleParser_Parse_comparison 演示各类比较运算符。
func ExampleParser_Parse_comparison() {
	psr := mongodsl.NewParser()

	filter, _ := psr.Parse(`age >= 18`)
	fmt.Println(mustJSON(filter))

	filter, _ = psr.Parse(`score < 60`)
	fmt.Println(mustJSON(filter))

	filter, _ = psr.Parse(`level != 0`)
	fmt.Println(mustJSON(filter))

	// Output:
	// {"age":{"$gte":18}}
	// {"score":{"$lt":60}}
	// {"level":{"$ne":0}}
}

// ExampleParser_Parse_logical 演示逻辑与、逻辑或的组合。
func ExampleParser_Parse_logical() {
	psr := mongodsl.NewParser()

	filter, _ := psr.Parse(`age >= 18 && city == "北京"`)
	fmt.Println(mustJSON(filter))

	filter, _ = psr.Parse(`role == "admin" || role == "ops"`)
	fmt.Println(mustJSON(filter))

	// Output:
	// {"$and":[{"age":{"$gte":18}},{"city":"北京"}]}
	// {"$or":[{"role":"admin"},{"role":"ops"}]}
}

// ExampleParser_Parse_in 演示集合成员判断及其取反形式。
func ExampleParser_Parse_in() {
	psr := mongodsl.NewParser()

	filter, _ := psr.Parse(`role in ["admin", "ops"]`)
	fmt.Println(mustJSON(filter))

	filter, _ = psr.Parse(`!(role in ["guest"])`)
	fmt.Println(mustJSON(filter))

	// Output:
	// {"role":{"$in":["admin","ops"]}}
	// {"role":{"$nin":["guest"]}}
}

// ExampleParser_Parse_matches 演示正则匹配。
func ExampleParser_Parse_matches() {
	psr := mongodsl.NewParser()

	filter, _ := psr.Parse(`name matches "^admin"`)
	fmt.Println(mustJSON(filter))

	// Output: {"name":{"$regex":"^admin"}}
}

// ExampleParser_Parse_exists 演示字段存在性判断及其取反形式。
func ExampleParser_Parse_exists() {
	psr := mongodsl.NewParser()

	filter, _ := psr.Parse(`exists("email")`)
	fmt.Println(mustJSON(filter))

	filter, _ = psr.Parse(`!exists("email")`)
	fmt.Println(mustJSON(filter))

	// Output:
	// {"email":{"$exists":true}}
	// {"email":{"$exists":false}}
}

// ExampleParser_Parse_nested 演示点号嵌套字段路径。
func ExampleParser_Parse_nested() {
	psr := mongodsl.NewParser()

	filter, _ := psr.Parse(`user.age < 30`)
	fmt.Println(mustJSON(filter))

	// Output: {"user.age":{"$lt":30}}
}

// ExampleParser_Parse_error 演示错误处理：使用了不支持的运算符或类型不匹配。
func ExampleParser_Parse_error() {
	psr := mongodsl.NewParser()

	// in 右侧必须是数组，否则返回错误。
	_, err := psr.Parse(`role in "admin"`)
	fmt.Println(err)

	// Output: 字段 'role' 的 in 操作符右侧必须是数组
}
