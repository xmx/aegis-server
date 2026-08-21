package mongodsl

import (
	"fmt"

	"github.com/expr-lang/expr/ast"
	"github.com/expr-lang/expr/parser"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// Parser 将类表达式语法的 DSL 字符串解析为 MongoDB 查询过滤器（bson.M）。
//
// Parser 基于 github.com/expr-lang/expr 的词法与语法分析器，遍历其生成的
// AST 并将其翻译为等价的 MongoDB 查询条件。它支持比较、逻辑组合、集合成员、
// 正则匹配以及一组内置函数（如 now、date、oid、exists）。
//
// 支持的运算符与对应的 MongoDB 语义：
//   - ==       字段等于                 {field: value}
//   - !=       字段不等于               {field: {$ne: value}}
//   - >        大于                     {field: {$gt: value}}
//   - >=       大于等于                 {field: {$gte: value}}
//   - <        小于                     {field: {$lt: value}}
//   - <=       小于等于                 {field: {$lte: value}}
//   - in       属于数组                 {field: {$in: [...]}}
//   - !(x in)  不属于数组               {field: {$nin: [...]}}
//   - matches  正则匹配                 {field: {$regex: "..."}}
//   - &&       逻辑与                   {$and: [...]}
//   - ||       逻辑或                   {$or: [...]}
//   - !exists  字段不存在               {field: {$exists: false}}
//
// 内置函数：
//   - now()              返回当前时间，可链式调用 .add("1h") 进行偏移
//   - date("2024-01-02") 将字符串解析为时间，支持 RFC3339、DateTime、DateOnly 等格式
//   - oid("...")         将 24 位十六进制字符串转换为 bson.ObjectID
//   - exists("field")    生成字段存在性判断 {field: {$exists: true}}
//
// 字段名支持点号嵌套（如 user.name），会被翻译为 MongoDB 的嵌套字段路径。
//
// 示例：
//
//	psr := mongodsl.NewParser()
//
//	// 简单等值：{"status": "active"}
//	filter, _ := psr.Parse(`status == "active"`)
//
//	// 逻辑组合：{"$and": [{"age": {"$gte": 18}}, {"city": "北京"}]}
//	filter, _ = psr.Parse(`age >= 18 && city == "北京"`)
//
//	// 集合成员：{"role": {"$in": ["admin", "ops"]}}
//	filter, _ = psr.Parse(`role in ["admin", "ops"]`)
//
//	// 取反成员：{"role": {"$nin": ["guest"]}}
//	filter, _ = psr.Parse(`!(role in ["guest"])`)
//
//	// 正则匹配：{"name": {"$regex": "^admin"}}
//	filter, _ = psr.Parse(`name matches "^admin"`)
//
//	// 字段存在：{"email": {"$exists": true}}
//	filter, _ = psr.Parse(`exists("email")`)
//
//	// 时间偏移：{"created_at": {"$gt": <now-1h>}}
//	filter, _ = psr.Parse(`created_at > now().add("-1h")`)
//
//	// 日期解析：{"created_at": {"$gte": <2024-01-01>}}
//	filter, _ = psr.Parse(`created_at >= date("2024-01-01")`)
//
//	// ObjectID：{"_id": <ObjectID>}
//	filter, _ = psr.Parse(`_id == oid("507f1f77bcf86cd799439011")`)
//
//	// 嵌套字段：{"user.age": {"$lt": 30}}
//	filter, _ = psr.Parse(`user.age < 30`)
type Parser struct {
	funcs functionMap // 支持的函数
}

// NewParser 创建并返回一个装载了内置 BSON 表达式函数的 Parser。
//
// 返回的 Parser 可复用于多次 Parse 调用。
func NewParser() *Parser {
	return &Parser{
		funcs: builtinBSONExpressionFunctions(),
	}
}

// Parse 将 DSL 表达式字符串解析为 MongoDB 查询过滤器 bson.M。
//
// 解析流程分为两步：先由 expr 词法/语法分析器生成 AST，随后遍历 AST 将其
// 翻译为 bson.M。若输入存在语法错误、使用了不支持的运算符或函数、或最终结果
// 不是一个合法的查询过滤器，则返回相应的错误（通常为 *DSLError）。
func (psr *Parser) Parse(input string) (bson.M, error) {
	tree, err := parser.Parse(input)
	if err != nil {
		return nil, err
	}

	val, err := psr.walkAST(tree.Node)
	if err != nil {
		return nil, err
	}

	if ret, ok := val.(bson.M); ok {
		return ret, nil
	}

	return nil, newDSLErrorf("无效的表达式")
}

// walkAST 递归遍历单个 AST 节点，并根据节点类型分派到对应的处理方法。
//
// 返回值可能是 bson.M（表示一个查询条件）、原生标量值（字符串、整数、浮点、
// 布尔、nil）、[]any（数组字面量）或实现了 object 接口的封装值（如时间对象）。
// 调用方通常需要通过 unwrapValue 将 object 解包为可放入 BSON 的原生类型。
func (psr *Parser) walkAST(node ast.Node) (any, error) {
	switch n := node.(type) {
	case *ast.BinaryNode:
		return psr.walkBinaryNode(n)
	case *ast.BuiltinNode:
		return psr.invokeFunction(n.Name, n.Arguments)
	case *ast.CallNode:
		return psr.walkCallNode(n)
	case *ast.ArrayNode:
		return psr.walkArrayNode(n)
	case *ast.UnaryNode:
		return psr.walkUnaryNode(n)
	case *ast.StringNode:
		return n.Value, nil
	case *ast.IntegerNode:
		return n.Value, nil
	case *ast.FloatNode:
		return n.Value, nil
	case *ast.BoolNode:
		return n.Value, nil
	case *ast.NilNode:
		return nil, nil
	default:
		return nil, newDSLErrorf("不支持的语法: %s", node.String())
	}
}

// walkBinaryNode 处理二元运算节点。
//
// 对于逻辑运算符（&& 与 ||），左右子节点各自被解析为子过滤器，并组合为
// $and / $or。对于其余比较类运算符（== != > >= < <= in matches），左侧被
// 解析为字段名，右侧被解析为值，再翻译为对应的 MongoDB 操作符表达式。
func (psr *Parser) walkBinaryNode(node *ast.BinaryNode) (any, error) {
	op := node.Operator
	switch op {
	case "&&", "||":
		left, err := psr.walkAST(node.Left)
		if err != nil {
			return nil, err
		}
		left = psr.unwrapValue(left) // 解包左侧

		right, err := psr.walkAST(node.Right)
		if err != nil {
			return nil, err
		}
		right = psr.unwrapValue(right) // 解包右侧

		if op == "&&" {
			return bson.M{"$and": []any{left, right}}, nil
		}

		return bson.M{"$or": []any{left, right}}, nil
	}

	key, err := psr.walkLeftKey(node.Left)
	if err != nil {
		return nil, err
	}

	val, err := psr.walkAST(node.Right)
	if err != nil {
		return nil, err
	}
	val = psr.unwrapValue(val) // 确保放入 BSON 中的是原生类型

	if op == "==" {
		return bson.M{key: val}, nil
	} else if op == "!=" {
		return bson.M{key: bson.M{"$ne": val}}, nil
	} else if op == ">" {
		return bson.M{key: bson.M{"$gt": val}}, nil
	} else if op == ">=" {
		return bson.M{key: bson.M{"$gte": val}}, nil
	} else if op == "<" {
		return bson.M{key: bson.M{"$lt": val}}, nil
	} else if op == "<=" {
		return bson.M{key: bson.M{"$lte": val}}, nil
	} else if op == "in" {
		if _, ok := val.([]any); !ok {
			return nil, newDSLErrorf("字段 '%s' 的 in 操作符右侧必须是数组", key)
		}

		return bson.M{key: bson.M{"$in": val}}, nil
	} else if op == "matches" {
		reg, ok := val.(string)
		if !ok {
			return nil, newDSLErrorf("字段 '%s' 的 matches 操作符右侧必须是正则表达式字符串", key)
		}
		regex := bson.Regex{Pattern: reg, Options: "i"}

		return bson.M{key: bson.M{"$regex": regex}}, nil
	}

	return nil, newDSLErrorf("不支持的运算符: %s", op)
}

// walkCallNode 处理函数调用节点。
//
// 支持两种调用形式：普通函数调用（如 now()、date("...")）以及基于接收者的
// 方法调用（如 now().add("1h")）。前者直接在内置函数表中查找并调用；后者先
// 解析出接收者对象，再在该对象上调用相应方法。
func (psr *Parser) walkCallNode(node *ast.CallNode) (any, error) {
	if ident, ok := node.Callee.(*ast.IdentifierNode); ok {
		return psr.invokeFunction(ident.Value, node.Arguments)
	}

	member, ok := node.Callee.(*ast.MemberNode)
	if !ok {
		return nil, newDSLErrorf("不支持的函数调用格式")
	}
	strnode, yes := member.Property.(*ast.StringNode)
	if !yes {
		return nil, newDSLErrorf("不支持的函数调用格式")
	}

	receiver, err := psr.walkAST(member.Node)
	if err != nil {
		return nil, err
	}

	var args []any
	method := strnode.Value
	for idx, argNode := range node.Arguments {
		arg, err1 := psr.walkAST(argNode)
		if err1 != nil {
			return nil, newDSLErrorf("方法 %s 的第 %d 个参数解析失败: %v", method, idx, err1)
		}
		args = append(args, arg)
	}

	return psr.invokeMethod(receiver, method, args)
}

// walkArrayNode 处理数组字面量节点（如 ["a", "b"]），逐一解析元素并解包为
// 原生类型，返回 []any。常用于 in 操作符右侧的取值集合。
func (psr *Parser) walkArrayNode(node *ast.ArrayNode) ([]any, error) {
	var arr []any
	for _, item := range node.Nodes {
		val, err := psr.walkAST(item)
		if err != nil {
			return nil, err
		}
		arr = append(arr, psr.unwrapValue(val))
	}

	return arr, nil
}

// walkUnaryNode 处理一元运算节点，仅支持取反运算符 '!'。
//
// 取反只能作用于 in 或 exists 表达式：!(x in [...]) 被翻译为 $nin，
// !exists("field") 被翻译为 {$exists: false}。用于其他表达式时返回错误。
func (psr *Parser) walkUnaryNode(node *ast.UnaryNode) (bson.M, error) {
	if node.Operator != "!" {
		return nil, newDSLErrorf("不支持的一元操作符: '%s'", node.Operator)
	}

	val, err := psr.walkAST(node.Node)
	if err != nil {
		return nil, err
	}
	val = psr.unwrapValue(val)

	filter, yes := val.(bson.M)
	if !yes {
		return nil, newDSLErrorf("运算符 '!' 只能用于 [in exists] 表达式取反操作")
	}

	for k, sub := range filter {
		smap, yes1 := sub.(bson.M)
		if !yes1 || len(smap) != 1 {
			return nil, fmt.Errorf("运算符 '!' 只能用于 [in exists] 表达式取反操作")
		}

		if vals, exists := smap["$in"]; exists {
			return bson.M{k: bson.M{"$nin": vals}}, nil
		}
		if _, exists := smap["$exists"]; exists {
			return bson.M{k: bson.M{"$exists": false}}, nil
		}

		return nil, newDSLErrorf("运算符 '!' 只能用于 [in exists] 表达式取反操作")
	}

	return nil, newDSLErrorf("运算符 '!' 只能用于 [in exists] 表达式取反操作")
}

// walkLeftKey 从比较表达式的左侧节点中提取字段名。
//
// 支持普通标识符（如 status）以及点号嵌套的成员访问（如 user.name），后者
// 会被递归拼接为以 '.' 分隔的 MongoDB 嵌套字段路径。若左侧不是合法字段名则
// 返回错误。
func (psr *Parser) walkLeftKey(node ast.Node) (string, error) {
	switch n := node.(type) {
	case *ast.IdentifierNode:
		return n.Value, nil
	case *ast.MemberNode:
		base, err := psr.walkLeftKey(n.Node)
		if err != nil {
			return "", err
		}
		if prop, ok := n.Property.(*ast.StringNode); ok {
			return base + "." + prop.Value, nil
		}
	}
	return "", newDSLErrorf("表达式左侧必须是合法的字段名")
}

// invokeFunction 解析函数实参并在内置函数表中查找并调用指定函数。
//
// 每个实参先经 walkAST 求值再解包为原生类型，随后按位置传入目标函数。
func (psr *Parser) invokeFunction(name string, nodes []ast.Node) (any, error) {
	var args []any
	for idx, node := range nodes {
		arg, err := psr.walkAST(node)
		if err != nil {
			return nil, newDSLErrorf("函数 '%s' 的第 %d 个参数解析失败: %v", name, idx, err)
		}
		args = append(args, psr.unwrapValue(arg))
	}

	return psr.funcs.invoke(name, args)
}

func (psr *Parser) invokeMethod(receiver any, method string, args []any) (any, error) {
	obj, ok := receiver.(object)
	if !ok {
		return nil, newDSLErrorf("此类型不支持调用 %s 方法", method)
	}

	return obj.invoke(method, args)
}

func (psr *Parser) unwrapValue(v any) any {
	if uv, ok := v.(object); ok {
		return uv.unwrap()
	}

	return v
}
