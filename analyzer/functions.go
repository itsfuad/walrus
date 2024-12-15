package analyzer

import (
	"fmt"
	"walrus/ast"
	"walrus/errgen"
)

func checkFunctionExpr(funcNode ast.FunctionLiteral, env *TypeEnvironment) TcValue {
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
		errgen.Add(env.filePath, funcNode.Start.Line, funcNode.End.Line, funcNode.Start.Column, funcNode.End.Column, "error declaring function. "+err.Error()).Level(errgen.CRITICAL)
	}
	//check the function body
	for _, stmt := range funcNode.Body.Contents {
		CheckAST(stmt, fnEnv)
	}

	return fn
}

func checkandDeclareParamaters(params []ast.FunctionParam, fnEnv *TypeEnvironment) []FnParam {
	var parameters []FnParam

	for i, param := range params {
		checkAndDeclareSingleParameter(param, i, params, fnEnv, &parameters)
	}
	return parameters
}

func checkAndDeclareSingleParameter(param ast.FunctionParam, i int, params []ast.FunctionParam, fnEnv *TypeEnvironment, parameters *[]FnParam) {
	if fnEnv.isDeclared(param.Identifier.Name) {
		errgen.Add(fnEnv.filePath, param.Identifier.Start.Line, param.Identifier.End.Line, param.Identifier.Start.Column, param.Identifier.End.Column, fmt.Sprintf("parameter '%s' is already defined", param.Identifier.Name)).Level(errgen.NORMAL)
	}

	paramType := evaluateTypeName(param.Type, fnEnv)

	if param.IsOptional {
		checkOptionalParameter(param, i, params, fnEnv, paramType)
	}

	err := fnEnv.declareVar(param.Identifier.Name, paramType, false, param.IsOptional)
	if err != nil {
		errgen.Add(fnEnv.filePath, param.Identifier.Start.Line, param.Identifier.End.Line, param.Identifier.Start.Column, param.Identifier.End.Column, fmt.Sprintf("error defining parameter. %s", err.Error())).Level(errgen.CRITICAL)
	}

	*parameters = append(*parameters, FnParam{
		Name:       param.Identifier.Name,
		IsOptional: param.IsOptional,
		Type:       paramType,
	})
}

func checkOptionalParameter(param ast.FunctionParam, i int, params []ast.FunctionParam, fnEnv *TypeEnvironment, paramType TcValue) {
	for j := i + 1; j < len(params); j++ {
		if !params[j].IsOptional {
			errgen.Add(fnEnv.filePath, params[j].Identifier.Start.Line, params[j].Identifier.End.Line, params[j].Identifier.Start.Column, params[j].Identifier.End.Column, fmt.Sprintf("parameter '%s' cannot be non-optional after an optional parameter", params[j].Identifier.Name)).Level(errgen.CRITICAL)
		}
	}

	defaultValue := parseNodeValue(param.DefaultValue, fnEnv)

	err := matchTypes(paramType, defaultValue)
	if err != nil {
		errgen.Add(fnEnv.filePath, param.DefaultValue.StartPos().Line, param.DefaultValue.EndPos().Line, param.DefaultValue.StartPos().Column, param.DefaultValue.EndPos().Column, fmt.Sprintf("error defining parameter. %s", err.Error())).Level(errgen.CRITICAL)
	}
}

func checkFunctionCall(callNode ast.FunctionCallExpr, env *TypeEnvironment) TcValue {
	//check if the function is declared
	caller := parseNodeValue(callNode.Caller, env)
	fn, err := userDefinedToFn(caller)

	if err != nil {
		errgen.Add(env.filePath, callNode.Caller.StartPos().Line, callNode.Caller.EndPos().Line, callNode.Caller.StartPos().Column, callNode.Caller.EndPos().Column, err.Error()).Level(errgen.CRITICAL)
	}

	fnParams := fn.Params
	if len(callNode.Arguments) != len(fnParams) {
		// exclude the optional parameters from the count
		optionalParams := 0
		for _, param := range fnParams {
			if param.IsOptional {
				optionalParams++
			}
		}
		if len(callNode.Arguments) < len(fnParams)-optionalParams {
			errgen.Add(env.filePath, callNode.Start.Line, callNode.End.Line, callNode.Start.Column, callNode.End.Column, fmt.Sprintf("function expects at least %d arguments, got %d", len(fnParams)-optionalParams, len(callNode.Arguments))).Level(errgen.NORMAL)
		}
		if len(callNode.Arguments) > len(fnParams) {
			errgen.Add(env.filePath, callNode.Start.Line, callNode.End.Line, callNode.Start.Column, callNode.End.Column, fmt.Sprintf("function expects at most %d arguments, got %d", len(fnParams), len(callNode.Arguments))).Level(errgen.NORMAL)
		}
	}

	//check if the arguments match the parameters
	for i := 0; i < len(callNode.Arguments); i++ {
		arg := parseNodeValue(callNode.Arguments[i], env)
		err := matchTypes(fnParams[i].Type, arg)
		if err != nil {
			errgen.Add(env.filePath, callNode.Arguments[i].StartPos().Line, callNode.Arguments[i].EndPos().Line, callNode.Arguments[i].StartPos().Column, callNode.Arguments[i].EndPos().Column, err.Error()).Level(errgen.NORMAL)
		}
	}

	return fn.Returns
}

func userDefinedToFn(ud TcValue) (Fn, error) {
	// if UserDefined then chain until Fn or error
	switch t := ud.(type) {
	case Fn:
		return t, nil
	case StructMethod:
		return t.Fn, nil
	case UserDefined:
		return userDefinedToFn(t.TypeDef)
	default:
		return Fn{}, fmt.Errorf("type of '%s' is not callable", ud.DType())
	}
}

func checkFunctionDeclStmt(funcNode ast.FunctionDeclStmt, env *TypeEnvironment) TcValue {

	// check if function is already declared
	funcName := funcNode.Identifier.Name

	if env.isDeclared(funcName) {
		errgen.Add(env.filePath, funcNode.Identifier.Start.Line, funcNode.Identifier.End.Line, funcNode.Identifier.Start.Column, funcNode.Identifier.End.Column, fmt.Sprintf("function '%s' is already defined in this scope", funcName)).Level(errgen.NORMAL)
	}

	return CheckAndDeclareFunction(funcNode.FunctionLiteral, funcName, env)
}

func getFunctionReturnValue(env *TypeEnvironment, returnNode ast.Node) TcValue {
	funcParent, err := env.resolveFunctionEnv()

	if err != nil {
		errgen.Add(env.filePath, returnNode.StartPos().Line, returnNode.EndPos().Line, returnNode.StartPos().Column, returnNode.EndPos().Column, err.Error()).Level(errgen.CRITICAL)
	}

	fnName := funcParent.scopeName
	//fn := funcParent.parent.variables[fnName].(Fn)
	switch fn := funcParent.parent.variables[fnName].(type) {
	case Fn:
		return fn.Returns
	case StructMethod:
		return fn.Fn.Returns
	default:

		errgen.Add(env.filePath, returnNode.StartPos().Line, returnNode.EndPos().Line, returnNode.StartPos().Column, returnNode.EndPos().Column, fmt.Sprintf("'%s' is not a function", fnName)).Level(errgen.NORMAL)
		return NewVoid()
	}
}
