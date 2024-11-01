package typechecker

import (
	"fmt"
	"walrus/ast"
	"walrus/errgen"
)

func checkFunctionExpr(funcNode ast.FunctionLiteral, env *TypeEnvironment) ValueTypeInterface {
	name := fmt.Sprintf("_FN_%s", RandStringRunes(10))
	return CheckAndDeclareFunction(funcNode, name, env)
}

func CheckAndDeclareFunction(funcNode ast.FunctionLiteral, name string, env *TypeEnvironment) Fn {

	fnEnv := NewTypeENV(env, FUNCTION_SCOPE, name, env.filePath)

	parameters := checkandDeclareParamaters(funcNode.Params, fnEnv)

	//check return type
	returnType := EvaluateTypeName(funcNode.ReturnType, fnEnv)

	fn := Fn{
		DataType:      FUNCTION_TYPE,
		Params:        parameters,
		Returns:       returnType,
		FunctionScope: *fnEnv,
	}

	//declare the function
	err := env.DeclareVar(name, fn, true, false)
	if err != nil {
		//errgen.AddError(env.filePath, funcNode.Start.Line, funcNode.End.Line, funcNode.Start.Column, funcNode.End.Column, err.Error()).DisplayWithPanic()
		errgen.AddError(env.filePath, funcNode.Start.Line, funcNode.End.Line, funcNode.Start.Column, funcNode.End.Column, err.Error())
	}

	//check the function body
	for _, stmt := range funcNode.Body.Contents {
		CheckAST(stmt, fnEnv)
	}

	return fn
}

func checkandDeclareParamaters(params []ast.FunctionParam, fnEnv *TypeEnvironment) []FnParam {

	var parameters []FnParam

	for _, param := range params {

		if fnEnv.isDeclared(param.Identifier.Name) {
			//errgen.AddError(fnEnv.filePath, param.Identifier.Start.Line, param.Identifier.End.Line, param.Identifier.Start.Column, param.Identifier.End.Column, fmt.Sprintf("Parameter %s is already declared", param.Identifier.Name)).DisplayWithPanic()
			errgen.AddError(fnEnv.filePath, param.Identifier.Start.Line, param.Identifier.End.Line, param.Identifier.Start.Column, param.Identifier.End.Column, fmt.Sprintf("Parameter %s is already declared", param.Identifier.Name))
		}

		paramType := EvaluateTypeName(param.Type, fnEnv)

		if param.IsOptional {
			//default value type
			defaultValue := GetValueType(param.DefaultValue, fnEnv)
			err := MatchTypes(paramType, defaultValue, fnEnv.filePath, param.DefaultValue.StartPos().Line, param.DefaultValue.EndPos().Line, param.DefaultValue.StartPos().Column, param.DefaultValue.EndPos().Column)
			if err != nil {
				//errgen.AddError(fnEnv.filePath, param.DefaultValue.StartPos().Line, param.DefaultValue.EndPos().Line, param.DefaultValue.StartPos().Column, param.DefaultValue.EndPos().Column, err.Error()).DisplayWithPanic()
				errgen.AddError(fnEnv.filePath, param.DefaultValue.StartPos().Line, param.DefaultValue.EndPos().Line, param.DefaultValue.StartPos().Column, param.DefaultValue.EndPos().Column, err.Error())
			}
		}

		err := fnEnv.DeclareVar(param.Identifier.Name, paramType, false, param.IsOptional)
		if err != nil {
			//errgen.AddError(fnEnv.filePath, param.Identifier.Start.Line, param.Identifier.End.Line, param.Identifier.Start.Column, param.Identifier.End.Column, err.Error()).DisplayWithPanic()
			errgen.AddError(fnEnv.filePath, param.Identifier.Start.Line, param.Identifier.End.Line, param.Identifier.Start.Column, param.Identifier.End.Column, err.Error())
		}

		fmt.Printf("Declared parameter %s of type %s\n", param.Identifier.Name, paramType.DType())

		parameters = append(parameters, FnParam{
			Name:       param.Identifier.Name,
			IsOptional: param.IsOptional,
			Type:       paramType,
		})
	}
	return parameters
}

func checkFunctionCall(callNode ast.FunctionCallExpr, env *TypeEnvironment) ValueTypeInterface {
	//check if the function is declared
	caller := GetValueType(callNode.Caller, env)
	fn, err := userDefinedToFn(caller)

	if err != nil {
		//errgen.AddError(env.filePath, callNode.Caller.StartPos().Line, callNode.Caller.EndPos().Line, callNode.Caller.StartPos().Column, callNode.Caller.EndPos().Column, err.Error()).DisplayWithPanic()
		errgen.AddError(env.filePath, callNode.Caller.StartPos().Line, callNode.Caller.EndPos().Line, callNode.Caller.StartPos().Column, callNode.Caller.EndPos().Column, err.Error())
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
			//errgen.AddError(env.filePath, callNode.Start.Line, callNode.End.Line, callNode.Start.Column, callNode.End.Column, fmt.Sprintf("Function expects at least %d arguments, got %d", len(fnParams)-optionalParams, len(callNode.Arguments))).DisplayWithPanic()
			errgen.AddError(env.filePath, callNode.Start.Line, callNode.End.Line, callNode.Start.Column, callNode.End.Column, fmt.Sprintf("Function expects at least %d arguments, got %d", len(fnParams)-optionalParams, len(callNode.Arguments)))
		}
		if len(callNode.Arguments) > len(fnParams) {
			//errgen.AddError(env.filePath, callNode.Start.Line, callNode.End.Line, callNode.Start.Column, callNode.End.Column, fmt.Sprintf("Function expects at most %d arguments, got %d", len(fnParams), len(callNode.Arguments))).DisplayWithPanic()
			errgen.AddError(env.filePath, callNode.Start.Line, callNode.End.Line, callNode.Start.Column, callNode.End.Column, fmt.Sprintf("Function expects at most %d arguments, got %d", len(fnParams), len(callNode.Arguments)))
		}
	}

	//check if the arguments match the parameters
	for i := 0; i < len(callNode.Arguments); i++ {
		arg := GetValueType(callNode.Arguments[i], env)
		err := MatchTypes(fnParams[i].Type, arg, env.filePath, callNode.Arguments[i].StartPos().Line, callNode.Arguments[i].EndPos().Line, callNode.Arguments[i].StartPos().Column, callNode.Arguments[i].EndPos().Column)
		if err != nil {
			//errgen.AddError(env.filePath, callNode.Arguments[i].StartPos().Line, callNode.Arguments[i].EndPos().Line, callNode.Arguments[i].StartPos().Column, callNode.Arguments[i].EndPos().Column, err.Error()).DisplayWithPanic()
			errgen.AddError(env.filePath, callNode.Arguments[i].StartPos().Line, callNode.Arguments[i].EndPos().Line, callNode.Arguments[i].StartPos().Column, callNode.Arguments[i].EndPos().Column, err.Error())
		}
	}

	return fn.Returns
}

func userDefinedToFn(ud ValueTypeInterface) (Fn, error) {
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

func checkFunctionDeclStmt(funcNode ast.FunctionDeclStmt, env *TypeEnvironment) ValueTypeInterface {

	// check if function is already declared
	funcName := funcNode.Identifier.Name

	if env.isDeclared(funcName) {
		//errgen.AddError(env.filePath, funcNode.Identifier.Start.Line, funcNode.Identifier.End.Line, funcNode.Identifier.Start.Column, funcNode.Identifier.End.Column, fmt.Sprintf("Function %s is already declared", funcName)).DisplayWithPanic()
		errgen.AddError(env.filePath, funcNode.Identifier.Start.Line, funcNode.Identifier.End.Line, funcNode.Identifier.Start.Column, funcNode.Identifier.End.Column, fmt.Sprintf("Function %s is already declared", funcName))
	}

	return CheckAndDeclareFunction(funcNode.FunctionLiteral, funcName, env)
}

func getFunctionReturnValue(env *TypeEnvironment, returnNode ast.Node) ValueTypeInterface {
	funcParent, err := env.ResolveFunctionEnv()
	if err != nil {
		//errgen.AddError(env.filePath, returnNode.StartPos().Line, returnNode.EndPos().Line, returnNode.StartPos().Column, returnNode.EndPos().Column, err.Error()).DisplayWithPanic()
		errgen.AddError(env.filePath, returnNode.StartPos().Line, returnNode.EndPos().Line, returnNode.StartPos().Column, returnNode.EndPos().Column, err.Error())
	}

	fnName := funcParent.scopeName
	//fn := funcParent.parent.variables[fnName].(Fn)
	switch fn := funcParent.parent.variables[fnName].(type) {
	case Fn:
		return fn.Returns
	case StructMethod:
		return fn.Fn.Returns
	default:
		//errgen.AddError(env.filePath, returnNode.StartPos().Line, returnNode.EndPos().Line, returnNode.StartPos().Column, returnNode.EndPos().Column, fmt.Sprintf("'%s' is not a function", fnName)).DisplayWithPanic()
		errgen.AddError(env.filePath, returnNode.StartPos().Line, returnNode.EndPos().Line, returnNode.StartPos().Column, returnNode.EndPos().Column, fmt.Sprintf("'%s' is not a function", fnName))
		return NewVoid()
	}
}
