package typechecker

import (
	//Standard packages
	"fmt"
	//Walrus packages
	"walrus/compiler/internal/ast"
	"walrus/compiler/report"
)

func checkFunctionExpr(funcNode ast.FunctionLiteral, env *TypeEnvironment) Tc {
	name := fmt.Sprintf("_FN_%s", RandStringRunes(10))
	return CheckAndDeclareFunction(funcNode, name, env)
}

func CheckAndDeclareFunction(funcNode ast.FunctionLiteral, name string, env *TypeEnvironment) Fn {

	fnEnv := NewTypeENV(env, FUNCTION_SCOPE, name, env.filePath)

	parameters := checkandDeclareParamaters(funcNode.Params, fnEnv)
	//check return type
	returnType := evaluateTypeName(funcNode.ReturnType, fnEnv)

	fn := Fn{
		DataType:      FUNCTION_TYPE,
		Params:        parameters,
		Returns:       returnType,
		FunctionScope: *fnEnv,
	}

	//declare the function
	err := env.declareVar(name, fn, true, false)
	if err != nil {
		report.Add(env.filePath, funcNode.Start.Line, funcNode.End.Line, funcNode.Start.Column, funcNode.End.Column, "error declaring function. "+err.Error()).SetLevel(report.CRITICAL_ERROR)
	}

	checkSatisfaction(funcNode, returnType, fnEnv)

	return fn
}

func checkSatisfaction(funcNode ast.FunctionLiteral, returnType Tc, fnEnv *TypeEnvironment) {
	//check the function body

	var fnSatisfied bool

	//return type is void, fnSatisfied is true by default
	if _, ok := returnType.(Void); ok {
		fnSatisfied = true
	}

	totalBlocks := 0
	unsatisfiedBlocks := make([]Block, 0)

	for _, stmt := range funcNode.Body.Contents {
		val := checkAST(stmt, fnEnv)
		switch v := val.(type) {
		case ReturnType:
			fnSatisfied = true
		case Block:
			totalBlocks++
			if !v.IsSatisfied {
				unsatisfiedBlocks = append(unsatisfiedBlocks, v)
			}
		}
	}

	checkSafeReturns(funcNode, unsatisfiedBlocks, fnSatisfied, totalBlocks, fnEnv)
}


func checkSafeReturns(funcNode ast.FunctionLiteral, unsatisfiedBlocks []Block, fnSatisfied bool, totalBlocks int, env *TypeEnvironment) {

	if len(unsatisfiedBlocks) > 0 {
		for _, block := range unsatisfiedBlocks {
			fnSatisfied = fnSatisfied || block.IsSatisfied
			if !fnSatisfied {
				report.Add(env.filePath, block.ProblemLocation.Start.Line, block.ProblemLocation.End.Line, block.ProblemLocation.Start.Column, block.ProblemLocation.End.Column, "missing return in this block").SetLevel(report.NORMAL_ERROR)
			}
		}
	} else if len(unsatisfiedBlocks) == 0 && totalBlocks > 0 {
		fnSatisfied = true
	}

	//if there is no return statement in the function, and the return type is not void. show an error
	//if the nested blocks are satisfied, then the function is satisfied. No explicit return is needed
	//if the nested blocks are not satisfied, then the function must have a return statement
	if totalBlocks == 0 && !fnSatisfied {
		report.Add(env.filePath, funcNode.Start.Line, funcNode.End.Line, funcNode.Start.Column, funcNode.End.Column, "missing return statement in function").SetLevel(report.NORMAL_ERROR)
	} else if !fnSatisfied {
		report.Add(env.filePath, funcNode.Start.Line, funcNode.End.Line, funcNode.Start.Column, funcNode.End.Column, "missing return statement in function").SetLevel(report.NORMAL_ERROR)
	}
}

func checkandDeclareParamaters(params []ast.FunctionParam, fnEnv *TypeEnvironment) []FnParam {
	var parameters []FnParam

	for _, param := range params {
		checkAndDeclareSingleParameter(param, fnEnv, &parameters)
	}
	return parameters
}

func checkAndDeclareSingleParameter(param ast.FunctionParam, fnEnv *TypeEnvironment, parameters *[]FnParam) {
	if fnEnv.isDeclared(param.Identifier.Name) {
		report.Add(fnEnv.filePath, param.Identifier.Start.Line, param.Identifier.End.Line, param.Identifier.Start.Column, param.Identifier.End.Column, fmt.Sprintf("parameter '%s' is already defined", param.Identifier.Name)).SetLevel(report.NORMAL_ERROR)
	}

	paramType := evaluateTypeName(param.Type, fnEnv)

	err := fnEnv.declareVar(param.Identifier.Name, paramType, false, false)
	if err != nil {
		report.Add(fnEnv.filePath, param.Identifier.Start.Line, param.Identifier.End.Line, param.Identifier.Start.Column, param.Identifier.End.Column, fmt.Sprintf("error defining parameter. %s", err.Error())).SetLevel(report.CRITICAL_ERROR)
	}

	*parameters = append(*parameters, FnParam{
		Name: param.Identifier.Name,
		Type: paramType,
	})
}

func checkFunctionCall(callNode ast.FunctionCallExpr, env *TypeEnvironment) Tc {
	//check if the function is declared
	caller := parseNodeValue(callNode.Caller, env)
	fn, err := userDefinedToFn(caller)

	if err != nil {
		report.Add(env.filePath, callNode.Caller.StartPos().Line, callNode.Caller.EndPos().Line, callNode.Caller.StartPos().Column, callNode.Caller.EndPos().Column, err.Error()).SetLevel(report.CRITICAL_ERROR)
	}

	fnParams := fn.Params
	if len(callNode.Arguments) != len(fnParams) {
		report.Add(env.filePath, callNode.Start.Line, callNode.End.Line, callNode.Start.Column, callNode.End.Column, fmt.Sprintf("function expects %d arguments, got %d", len(fnParams), len(callNode.Arguments))).SetLevel(report.NORMAL_ERROR)
	}

	//check if the arguments match the parameters
	for i := 0; i < len(callNode.Arguments); i++ {
		arg := parseNodeValue(callNode.Arguments[i], env)
		err := validateTypeCompatibility(fnParams[i].Type, arg)
		if err != nil {
			report.Add(env.filePath, callNode.Arguments[i].StartPos().Line, callNode.Arguments[i].EndPos().Line, callNode.Arguments[i].StartPos().Column, callNode.Arguments[i].EndPos().Column, err.Error()).SetLevel(report.NORMAL_ERROR)
		}
	}

	return fn.Returns
}

func userDefinedToFn(ud Tc) (Fn, error) {
	// if UserDefined then chain until Fn or error
	switch t := ud.(type) {
	case Fn:
		return t, nil
	case StructMethod:
		return t.Fn, nil
	case UserDefined:
		return userDefinedToFn(t.TypeDef)
	default:
		return Fn{}, fmt.Errorf("type of '%s' is not callable", tcToString(ud))
	}
}

func checkFunctionDeclStmt(funcNode ast.FunctionDeclStmt, env *TypeEnvironment) Tc {

	// check if function is already declared
	funcName := funcNode.Identifier.Name

	if env.isDeclared(funcName) {
		report.Add(env.filePath, funcNode.Identifier.Start.Line, funcNode.Identifier.End.Line, funcNode.Identifier.Start.Column, funcNode.Identifier.End.Column, fmt.Sprintf("function '%s' is already defined in this scope", funcName)).SetLevel(report.NORMAL_ERROR)
	}

	return CheckAndDeclareFunction(funcNode.FunctionLiteral, funcName, env)
}

func getFunctionReturnValue(env *TypeEnvironment, returnNode ast.Node) Tc {
	funcParent, err := env.resolveFunctionEnv()

	if err != nil {
		report.Add(env.filePath, returnNode.StartPos().Line, returnNode.EndPos().Line, returnNode.StartPos().Column, returnNode.EndPos().Column, err.Error()).SetLevel(report.CRITICAL_ERROR)
	}

	fnName := funcParent.scopeName
	//fn := funcParent.parent.variables[fnName].(Fn)
	switch fn := funcParent.parent.variables[fnName].(type) {
	case Fn:
		return fn.Returns
	case StructMethod:
		return fn.Fn.Returns
	default:
		report.Add(env.filePath, returnNode.StartPos().Line, returnNode.EndPos().Line, returnNode.StartPos().Column, returnNode.EndPos().Column, fmt.Sprintf("'%s' is not a function", fnName)).SetLevel(report.NORMAL_ERROR)
		return NewVoid()
	}
}
