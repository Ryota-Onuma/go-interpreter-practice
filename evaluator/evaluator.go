package evaluator

import (
	"lox-by-go/ast"
	"lox-by-go/object"
)

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node, env)
	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)
	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right)
	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		right := Eval(node.Right, env)
		if isError(left) {
			return left
		}
		if isError(right) {
			return right
		}
		switch {
		case object.IsNumber(left) && object.IsNumber(right):
			if left.Type() == object.FLOAT && right.Type() == object.INTEGER {
				right = object.NewFloat(float64(right.(*object.Integer).Value))
				return evalFloatInfixExpression(node.Operator, left, right)
			}
			if left.Type() == object.INTEGER && right.Type() == object.FLOAT {
				left = object.NewFloat(float64(left.(*object.Integer).Value))
				return evalFloatInfixExpression(node.Operator, left, right)
			}

			if left.Type() == object.INTEGER && right.Type() == object.INTEGER {
				return evalIntegerInfixExpression(node.Operator, left, right)
			}

			if left.Type() == object.FLOAT && right.Type() == object.FLOAT {
				return evalFloatInfixExpression(node.Operator, left, right)
			}
		case left.Type() == object.STRING && right.Type() == object.STRING:
			return evalStringInfixExpression(node.Operator, left, right)
		default:
			return object.NewError("unknown operator: %s %s %s", left.Type(), node.Operator, right.Type())
		}

	case *ast.BlockStatement:
		return evalBlockStatement(node, env)
	case *ast.IfExpression:
		return evalIfExpression(node, env)
	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue, env)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}
	case *ast.VarStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		env.Set(node.Name.Value, val)
	case *ast.Identifier:
		val, ok := env.Get(node.Value)
		if !ok {
			return object.NewError("undefined identifier %v", node.Value)
		}
		return val
	case *ast.IntegerLiteral:
		return object.NewInteger(node.Value)
	case *ast.FloatLiteral:
		return object.NewFloat(node.Value)
	case *ast.Boolean:
		return object.NewBoolean(node.Value)
	case *ast.StringLiteral:
		return object.NewString(node.Value)
	}
	return nil
}

func evalProgram(program *ast.Program, env *object.Environment) object.Object {
	var result object.Object
	for _, statement := range program.Statements {
		result = Eval(statement, env)
		if result != nil && result.Type() == object.RETURN {
			return result.(*object.ReturnValue).Value // アンラップ
		}
	}
	return result
}

func evalBlockStatement(block *ast.BlockStatement, env *object.Environment) object.Object {
	var result object.Object
	for _, stmt := range block.Statements {
		result = Eval(stmt, env)
		if result != nil && (result.Type() == object.RETURN || result.Type() == object.ERROR) {
			return result
		}
	}
	return result
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	}
	return object.NewError("unknown operator: %s%s", operator, right.Type())
}

// !true => false
// !false => true
// !5 => false
// !0 => true
func evalBangOperatorExpression(right object.Object) object.Object {
	if right.IsTruthy() {
		return object.NewBoolean(false)
	} else {
		return object.NewBoolean(true)
	}

}

func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.INTEGER && right.Type() != object.FLOAT {
		return object.NewError("unknown operator: -%s", right.Type())
	}

	if right.Type() == object.INTEGER {
		value := right.(*object.Integer).Value
		return object.NewInteger(-value)
	}

	if right.Type() == object.FLOAT {
		value := right.(*object.Float).Value
		return object.NewFloat(-value)

	}
	return object.NewError("unknown operator: -%s", right.Type())
}

func evalIntegerInfixExpression(operator string, left, right object.Object) object.Object {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value
	switch operator {
	case "+":
		return object.NewInteger(leftVal + rightVal)
	case "-":
		return object.NewInteger(leftVal - rightVal)
	case "*":
		return object.NewInteger(leftVal * rightVal)
	case "/":
		return object.NewInteger(leftVal / rightVal)
	case "<":
		return object.NewBoolean(leftVal < rightVal)
	case ">":
		return object.NewBoolean(leftVal > rightVal)
	case "==":
		return object.NewBoolean(leftVal == rightVal)
	case "!=":
		return object.NewBoolean(leftVal != rightVal)
	}
	return object.NewError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
}

func evalFloatInfixExpression(operator string, left, right object.Object) object.Object {
	leftVal := left.(*object.Float).Value
	rightVal := right.(*object.Float).Value
	switch operator {
	case "+":
		return object.NewFloat(leftVal + rightVal)
	case "-":
		return object.NewFloat(leftVal - rightVal)
	case "*":
		return object.NewFloat(leftVal * rightVal)
	case "/":
		return object.NewFloat(leftVal / rightVal)
	case "<":
		return object.NewBoolean(leftVal < rightVal)
	case ">":
		return object.NewBoolean(leftVal > rightVal)
	case "==":
		return object.NewBoolean(leftVal == rightVal)
	case "!=":
		return object.NewBoolean(leftVal != rightVal)
	}
	return object.NewError("unknown operator: %s %s %s", left.Type(), operator, right.Type())

}

func evalStringInfixExpression(operator string, left, right object.Object) object.Object {
	leftVal := left.(*object.String).Value
	rightVal := right.(*object.String).Value
	switch operator {
	case "+":
		return object.NewString(leftVal + rightVal)
	}
	return object.NewError("unknown operator: %s %s %s", left.Type(), operator, right.Type())

}

func evalIfExpression(ie *ast.IfExpression, env *object.Environment) object.Object {
	condition := Eval(ie.Condition, env)
	if isError(condition) {
		return condition
	}
	if condition.IsTruthy() {
		return Eval(ie.Consequence, env)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative, env)
	}
	return object.NewNil()
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR
	}
	return false
}
