package evaluator

import (
	"fmt"
	"interpreter/ast"
	"interpreter/object"
	"reflect"
)

var (
	TRUE  = &object.Boolean{Value: true}
	NULL  = &object.Null{}
	FALSE = &object.Boolean{Value: false}
)

func Eval(node ast.Node, environment *object.Environment) object.Object {

	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node, environment)
	case *ast.ExpressionStatement:
		return Eval(node.Expression, environment)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue, environment)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}
	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)
	case *ast.PrefixExpression:
		// приклад !!5 буде right FALSE
		// оскільки ми підемо в default в evalBangOperatorExpression
		// і після заходу в evalPrefix буде в нас true
		right := Eval(node.Right, environment)
		if isError(right) {
			return right
		}
		return evalPrefix(node.Operator, right)
	case *ast.InfixExpression:
		left := Eval(node.Left, environment)
		if isError(left) {
			return left
		}
		right := Eval(node.Right, environment)
		if isError(right) {
			return right
		}
		return evalInfixExpression(left, right, node.Operator)
	case *ast.IfExpression:
		return evalIfExpression(node, environment)
	case *ast.BlockStatement:
		return evalBlockStatements(node.Statements, environment)
	case *ast.LetStatement:
		val := Eval(node.Value, environment)
		if isError(val) {
			return val
		}
		//тут ми в середовище(наший сторедж) будемо зберігати значення під назвою змінної 'let a = b'  (map (key=a, val=b))
		environment.Set(node.Name.Value, val)
	case *ast.Identifier:
		return evalIdentifier(node, environment)

	case *ast.FunctionLiteral:
		params := node.Parameters
		return &object.Function{
			Parameters: params,
			Body:       &node.Body,
			Env:        environment,
		}
	case *ast.CallExpression:
		function := Eval(node.Function, environment)
		if isError(function) {
			return function
		}
		args := evalExpressions(node.Arguments, environment)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}
	}

	return nil
}

func evalExpressions(expressions []ast.Expression, environment *object.Environment) []object.Object {
	var result []object.Object

	for _, exp := range expressions {
		argRes := Eval(exp, environment)
		if isError(argRes) {
			return []object.Object{argRes}
		}
		result = append(result, argRes)
	}
	return result
}

func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	val, ok := env.Get(node.Value)
	if !ok {
		return newError("identifier not found: %s", node.Value)
	}
	return val
}

func evalIfExpression(ie *ast.IfExpression, environment *object.Environment) object.Object {
	cond := Eval(ie.Condition, environment)
	if isError(cond) {
		return cond
	}
	if isTruthy(cond) {
		return Eval(ie.Consequence, environment)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative, environment)
	} else {
		return NULL
	}
}

func isTruthy(cond object.Object) bool {
	switch cond {
	case NULL:
		return false
	case TRUE:
		return true
	case FALSE:
		return false
	default:
		return true
	}
}

func evalInfixExpression(left object.Object, right object.Object, operator string) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left.(*object.Integer), right.(*object.Integer))
	case operator == "==":
		return nativeBoolToBooleanObject(left == right)
	case operator == "!=":
		return nativeBoolToBooleanObject(left != right)
	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s", left.Type(), operator, right.Type())
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalIntegerInfixExpression(operator string, left, right *object.Integer) object.Object {
	return processIntegerInfixExpression(left, right, operator)
}

func processIntegerInfixExpression(left, right *object.Integer, operator string) object.Object {
	switch operator {
	case "*":
		return &object.Integer{Value: left.Value * right.Value}
	case "-":
		return &object.Integer{Value: left.Value - right.Value}
	case "/":
		return &object.Integer{Value: left.Value / right.Value}
	case "+":
		return &object.Integer{Value: left.Value + right.Value}
	case ">":
		return nativeBoolToBooleanObject(left.Value > right.Value)
	case "==":
		return nativeBoolToBooleanObject(left.Value == right.Value)
	case "!=":
		return nativeBoolToBooleanObject(left.Value != right.Value)
	case "<":
		return nativeBoolToBooleanObject(left.Value < right.Value)
	default:
		return newError("unknown operator: %s %s %s", reflect.TypeOf(left), operator, reflect.TypeOf(right))
	}
}

func evalPrefix(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusOperatorExpression(right)
	default:
		return newError("unknown operator: %s%s", operator, right.Type())
	}
}

func evalMinusOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJ {
		return newError("unknown operator: -%s", right.Type())
	}

	value := right.(*object.Integer).Value

	return &object.Integer{Value: -value}
}

func evalBangOperatorExpression(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		return FALSE
	}
}

func nativeBoolToBooleanObject(b bool) *object.Boolean {
	if b {
		return TRUE
	}
	return FALSE
}

func evalProgram(program *ast.Program, env *object.Environment) object.Object {

	return evalStatements(program.Statements, env)
}

func evalBlockStatements(statements []ast.Statement, environment *object.Environment) object.Object {
	var result object.Object

	for _, stmt := range statements {
		result = Eval(stmt, environment)
		if result.Type() == object.RETURN_VALUE_OBJ || result.Type() == object.ERROR_OBJ {
			return result
		}
	}

	return result
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false
}

func evalStatements(statements []ast.Statement, env *object.Environment) object.Object {
	var result object.Object

	for _, stmt := range statements {
		result = Eval(stmt, env)

		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
	}

	return result
}

func newError(format string, a ...any) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}
