package typechecker

import (
	"fmt"
	"walrus/ast"
	"walrus/builtins"
	"walrus/errors"
)

func EvaluateProgram(program ast.ProgramStmt, env *TypeEnvironment) ValueTypeInterface {
	fmt.Printf("Evaluating program\n")
	for _, item := range program.Contents {
		EvaluateTypesOfNode(item, env)
	}

	return nil
}

func EvaluateTypesOfNode(node ast.Node, env *TypeEnvironment) ValueTypeInterface {
	switch t := node.(type) {
	case ast.ProgramStmt:
		return EvaluateProgram(t, env)
	case ast.VarDeclStmt:
		return checkVariableDeclaration(t, env)
	case ast.VarAssignmentExpr:
		return checkVariableAssignment(t.Assignee, t.Value, env)
	case ast.IdentifierExpr:
		return checkIdentifier(t, env)
	case ast.IntegerLiteralExpr:
		return Int{
			DataType: INT_TYPE,
			Name:     string(INT_TYPE),
		}
	case ast.FloatLiteralExpr:
		return Float{
			DataType: FLOAT_TYPE,
			Name:     string(FLOAT_TYPE),
		}
	case ast.StringLiteralExpr:
		return Str{
			DataType: STRING_TYPE,
			Name:     string(STRING_TYPE),
		}
	case ast.CharLiteralExpr:
		return Chr{
			DataType: CHAR_TYPE,
			Name:     string(CHAR_TYPE),
		}
	case ast.NullLiteralExpr:
		return Null{
			DataType: NULL_TYPE,
			Name:     string(NULL_TYPE),
		}
	case ast.ArrayExpr:

		var expectedType ValueTypeInterface

		for i, value := range t.Values {
			v := EvaluateTypesOfNode(value, env)
			if i == 0 {
				expectedType = v
			}
			//check every type is same or not
			MatchTypes(expectedType, v, env.filePath, t.Start.Line, t.Values[i].StartPos().Column, t.Values[i].EndPos().Column)
		}

		return Array{
			DataType:  ARRAY_TYPE,
			ArrayType: expectedType,
		}

	case ast.ArrayIndexAccess:
		//Array must be evaluated to an array value
		arrType := EvaluateTypesOfNode(t.Arrayvalue, env)
		if _, ok := arrType.(Array); !ok {
			line := t.Arrayvalue.StartPos().Line
			start := t.Arrayvalue.StartPos().Column
			end := t.Arrayvalue.EndPos().Column
			errors.MakeError(env.filePath, line, start, end, fmt.Sprintf("cannot access index of type %s", arrType.DType())).AddHint("type must be an array", errors.TEXT_HINT).Display()
		}
		//index must be evaluated to int
		indexType := EvaluateTypesOfNode(t.Index, env)
		if _, ok := indexType.(Int); !ok {
			line := t.Index.StartPos().Line
			start := t.Index.StartPos().Column
			end := t.Index.EndPos().Column
			errors.MakeError(env.filePath, line, start, end, fmt.Sprintf("cannot use index of type %s", indexType.DType())).AddHint("index must be valid integer", errors.TEXT_HINT).Display()
		}
		return arrType.(Array).ArrayType
	default:
		errors.MakeError(env.filePath, node.StartPos().Line, node.StartPos().Column, node.EndPos().Column, fmt.Sprintf("<%T> node is not implemented yet", t)).Display()
		return nil
	}
}

func EvaluateTypeName(dtype ast.DataType, env *TypeEnvironment) ValueTypeInterface {
	switch t := dtype.(type) {
	case ast.ArrayType:
		arr := Array{
			DataType:  builtins.ARRAY,
			ArrayType: EvaluateTypeName(t.ArrayType, env),
		}
		return arr
	default:
		return makeBuiltinTYPE(VALUE_TYPE(t.Type()))
	}
}

func checkIdentifier(node ast.IdentifierExpr, env *TypeEnvironment) ValueTypeInterface {

	name := node.Name
	//find the scope where the variable was declared
	scope, err := env.ResolveVar(name)
	if err != nil {
		errors.MakeError(env.filePath, node.StartPos().Line, node.StartPos().Column, node.EndPos().Column, err.Error()).Display()
	}
	// if we found value on that scope, return the value. Else make error (though there is no change to reach the error)
	if val, ok := scope.variables[name]; ok {
		return val
	}
	errors.MakeError(env.filePath, node.StartPos().Line, node.StartPos().Column, node.EndPos().Column, "failed to check type").Display()
	return nil
}

// generate interfaces from the type enum
func makeBuiltinTYPE(typ VALUE_TYPE) ValueTypeInterface {
	switch typ {
	case INT_TYPE:
		return Int{
			DataType: typ,
			Name:     string(typ),
		}
	case FLOAT_TYPE:
		return Float{
			DataType: typ,
			Name:     string(typ),
		}
	case CHAR_TYPE:
		return Chr{
			DataType: typ,
			Name:     string(typ),
		}
	case STRING_TYPE:
		return Str{
			DataType: typ,
			Name:     string(typ),
		}
	case BOOLEAN_TYPE:
		return Bool{
			DataType: typ,
			Name:     string(typ),
		}
	case NULL_TYPE:
		return Null{
			DataType: typ,
			Name:     string(typ),
		}
	case VOID_TYPE:
		return Void{
			DataType: typ,
			Name:     string(typ),
		}
	}
	return nil
}

func checkVariableDeclaration(node ast.VarDeclStmt, env *TypeEnvironment) ValueTypeInterface {

	varToDecl := node.Variable

	var expectedTypeInterface ValueTypeInterface

	if node.ExplicitType != nil {
		//Explicit type is defined
		switch t := node.ExplicitType.(type) {
		case ast.ArrayType:
			expectedTypeInterface = EvaluateTypeName(t, env)
		default:
			expectedTypeInterface = makeBuiltinTYPE(VALUE_TYPE(node.ExplicitType.Type()))
		}
	} else {
		typ := EvaluateTypesOfNode(node.Value, env)
		switch t := typ.(type) {
		case Struct:
			err := env.DeclareStruct(varToDecl.Name, t)
			if err != nil {
				errors.MakeError(env.filePath, node.Variable.StartPos().Line, node.Variable.StartPos().Column, node.Variable.EndPos().Column, err.Error()).Display()
			}
		}
		expectedTypeInterface = typ
	}

	if node.IsAssigned && node.ExplicitType != nil {
		providedValue := EvaluateTypesOfNode(node.Value, env)
		MatchTypes(expectedTypeInterface, providedValue, env.filePath, node.Value.StartPos().Line, node.Value.StartPos().Column, node.Value.EndPos().Column)
	}

	err := env.DeclareVar(varToDecl.Name, expectedTypeInterface, node.IsConst)
	if err != nil {
		errors.MakeError(env.filePath, node.Variable.StartPos().Line, node.Variable.StartPos().Column, node.Variable.EndPos().Column, err.Error()).Display()
	}
	return nil
}

func getTypename(typeName ValueTypeInterface) VALUE_TYPE {
	switch t := typeName.(type) {
	case Array:
		return "[]" + getTypename(t.ArrayType)
	case Struct:
		return VALUE_TYPE(t.StructName)
	case Fn:
		return VALUE_TYPE(t.DataType)
	default:
		return t.DType()
	}
}

func MatchTypes(expected ValueTypeInterface, provided ValueTypeInterface, filePath string, line, start, end int) {
	expectedType := getTypename(expected)
	gotType := getTypename(provided)
	if expectedType != gotType {
		errors.MakeError(filePath, line, start, end, fmt.Sprintf("expected '%s', got '%s'", expectedType, gotType)).Display()
	}
}

func checkVariableAssignment(Assignee ast.Node, valueToAssign ast.Node, env *TypeEnvironment) ValueTypeInterface {

	//varToAssign := node.Identifier
	expected := EvaluateTypesOfNode(Assignee, env)
	provided := EvaluateTypesOfNode(valueToAssign, env)
	
	MatchTypes(expected, provided, env.filePath, valueToAssign.StartPos().Line, valueToAssign.StartPos().Column, valueToAssign.EndPos().Column)

	var varName string

	switch assignee := Assignee.(type) {
	case ast.IdentifierExpr:
		varName = assignee.Name
	case ast.ArrayIndexAccess:
		return nil
	default:
		panic("cannot assign to this type")
	}

	//get the stored type
	scope, err := env.ResolveVar(varName)
	if err != nil {
		errors.MakeError(env.filePath, valueToAssign.StartPos().Line, valueToAssign.StartPos().Column, valueToAssign.EndPos().Column, err.Error()).Display()
	}

	//if constant
	if scope.constants[varName] {
		errors.MakeError(env.filePath, valueToAssign.StartPos().Line, valueToAssign.StartPos().Column, valueToAssign.EndPos().Column, fmt.Sprintf("'%s' is constant", varName)).AddHint("cannot assign value to constant variables", errors.TEXT_HINT).Display()
	}
	scope.variables[varName] = provided
	return nil
}
