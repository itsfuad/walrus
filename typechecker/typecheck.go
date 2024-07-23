package typechecker

import (
	"fmt"
	"walrus/ast"
	"walrus/builtins"
	"walrus/errors"
	"walrus/utils"
)

func EvaluateProgram(program ast.ProgramStmt, env *TypeEnvironment) ValueTypeInterface {
	utils.ColorPrint(utils.ORANGE, "### Evaluating program ###")
	for _, item := range program.Contents {
		EvaluateTypesOfNode(item, env)
	}

	utils.ColorPrint(utils.GREEN, "--------- passed ---------")

	return nil
}

func EvaluateTypesOfNode(node ast.Node, env *TypeEnvironment) ValueTypeInterface {
	switch t := node.(type) {
	case ast.ProgramStmt:
		return EvaluateProgram(t, env)
	case ast.VarDeclStmt:
		return checkVariableDeclaration(t, env)
	case ast.VarAssignmentExpr:
		return checkVariableAssignment(t, env)
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
		return evaluateArrayExpr(t, env)
	case ast.ArrayIndexAccess:
		return evaluateArrayAccess(t, env)
	case ast.TypeDeclStmt:
		return checkTypeDeclaration(t, env)
	case ast.StructLiteral:
		return checkStructLiteral(t, env)
	}
	errors.MakeError(env.filePath, node.StartPos().Line, node.StartPos().Column, node.EndPos().Column, fmt.Sprintf("<%T> node is not implemented yet", node)).Display()
	return nil
}

func checkStructLiteral(structLit ast.StructLiteral, env *TypeEnvironment) ValueTypeInterface {

	sName := structLit.Name
	//check if defined
	scope, err := env.ResolveType(sName.Name)
	if err != nil {
		errors.MakeError(env.filePath, sName.StartPos().Line, sName.StartPos().Column, sName.EndPos().Column, fmt.Sprintf("type '%s' is not defined", sName.Name)).Display()
	}

	udType := scope.types[sName.Name].(UserDefined)

	structType := udType.TypeDef.(Struct)

	// now we match the defined props with the provided props
	// first check names on the provided value
	for propName, propValue := range structLit.Properties {
		//if prop exist
		if elem, ok := structType.Elements[propName]; !ok{
			errors.MakeError(env.filePath, propValue.StartPos().Line, propValue.StartPos().Column, propValue.EndPos().Column, fmt.Sprintf("property '%s' does not exists on type '%s'", propName, sName.Name)).Display()
		} else {
			MatchTypes(elem.Type, EvaluateTypesOfNode(propValue, env), env.filePath, propValue.StartPos().Line, propValue.StartPos().Column, propValue.EndPos().Column)
		}
	}

	//now check all from defined type
	for propName := range structType.Elements {
		if _, ok := structLit.Properties[propName]; !ok {
			//if does not exists
			errors.MakeError(env.filePath, structLit.StartPos().Line, structLit.StartPos().Column, structLit.EndPos().Column, fmt.Sprintf("property '%s' is missing from type '%s'", propName, sName.Name)).Display()
		}
	}

	return Struct{
		DataType: STRUCT_TYPE,
		StructName: sName.Name,
		Elements: structType.Elements,
	}
}

func evaluateArrayAccess(array ast.ArrayIndexAccess, env *TypeEnvironment) ValueTypeInterface {
	//Array must be evaluated to an array value
	arrType := EvaluateTypesOfNode(array.Arrayvalue, env)
	if _, ok := arrType.(Array); !ok {
		line := array.Arrayvalue.StartPos().Line
		start := array.Arrayvalue.StartPos().Column
		end := array.Arrayvalue.EndPos().Column
		errors.MakeError(env.filePath, line, start, end, fmt.Sprintf("cannot access index of type %s", arrType.DType())).AddHint("type must be an array", errors.TEXT_HINT).Display()
	}
	//index must be evaluated to int
	indexType := EvaluateTypesOfNode(array.Index, env)
	if _, ok := indexType.(Int); !ok {
		line := array.Index.StartPos().Line
		start := array.Index.StartPos().Column
		end := array.Index.EndPos().Column
		errors.MakeError(env.filePath, line, start, end, fmt.Sprintf("cannot use index of type %s", indexType.DType())).AddHint("index must be valid integer", errors.TEXT_HINT).Display()
	}
	return arrType.(Array).ArrayType
}

func evaluateArrayExpr(array ast.ArrayExpr, env *TypeEnvironment) ValueTypeInterface {
	var expectedType ValueTypeInterface
	for i, value := range array.Values {
		v := EvaluateTypesOfNode(value, env)
		if i == 0 {
			expectedType = v
		}
		//check every type is same or not
		MatchTypes(expectedType, v, env.filePath, array.Start.Line, array.Values[i].StartPos().Column, array.Values[i].EndPos().Column)
	}

	return Array{
		DataType:  ARRAY_TYPE,
		ArrayType: expectedType,
	}
}

func EvaluateTypeName(dtype ast.DataType, env *TypeEnvironment) (ValueTypeInterface, error) {
	switch t := dtype.(type) {
	case ast.ArrayType:
		val, err := EvaluateTypeName(t.ArrayType, env)
		if err != nil {
			return nil, err
		}
		arr := Array{
			DataType:  builtins.ARRAY,
			ArrayType: val,
		}
		return arr, nil
	default:
		fmt.Println("data type used:", dtype.Type())
		return makeBuiltinTYPE(VALUE_TYPE(t.Type()), env)
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
func makeBuiltinTYPE(typ VALUE_TYPE, env *TypeEnvironment) (ValueTypeInterface, error) {
	switch typ {
	case INT_TYPE:
		return Int{
			DataType: typ,
			Name:     string(typ),
		}, nil
	case FLOAT_TYPE:
		return Float{
			DataType: typ,
			Name:     string(typ),
		}, nil
	case CHAR_TYPE:
		return Chr{
			DataType: typ,
			Name:     string(typ),
		}, nil
	case STRING_TYPE:
		return Str{
			DataType: typ,
			Name:     string(typ),
		}, nil
	case BOOLEAN_TYPE:
		return Bool{
			DataType: typ,
			Name:     string(typ),
		}, nil
	case NULL_TYPE:
		return Null{
			DataType: typ,
			Name:     string(typ),
		}, nil
	case VOID_TYPE:
		return Void{
			DataType: typ,
			Name:     string(typ),
		}, nil
	default:
		fmt.Println("use type:", typ)
		//search for the type
		e, err := env.ResolveType(string(typ))
		if err != nil {
			return nil, err
		}
		return UserDefined{
			DataType: typ,
			TypeDef: e.types[string(typ)],
		}, nil
	}
}

func handleExplicitType(node ast.VarDeclStmt, env *TypeEnvironment) ValueTypeInterface {
	//Explicit type is defined
	var expectedTypeInterface ValueTypeInterface
	switch t := node.ExplicitType.(type) {
	case ast.ArrayType:
		val, err := EvaluateTypeName(t, env)
		if err != nil {
			errors.MakeError(env.filePath, t.Start.Line, t.Start.Column, t.End.Column, err.Error()).Display()
		}
		expectedTypeInterface = val
	default:
		val, err := makeBuiltinTYPE(VALUE_TYPE(node.ExplicitType.Type()), env)
		if err != nil {
			errors.MakeError(env.filePath, node.ExplicitType.StartPos().Line, node.ExplicitType.StartPos().Column, node.ExplicitType.EndPos().Column, err.Error()).Display()
		}
		expectedTypeInterface = val
	}
	return expectedTypeInterface
}

func checkVariableDeclaration(node ast.VarDeclStmt, env *TypeEnvironment) ValueTypeInterface {

	varToDecl := node.Variable

	var expectedTypeInterface ValueTypeInterface

	if node.ExplicitType != nil {
		expectedTypeInterface = handleExplicitType(node, env)
	} else {
		typ := EvaluateTypesOfNode(node.Value, env)
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

func checkTypeDeclaration(node ast.TypeDeclStmt, env *TypeEnvironment) ValueTypeInterface {

	typeName := node.UDType
	fmt.Println("UDT:", typeName.Type())

	var val ValueTypeInterface

	switch t := typeName.(type) {
	case ast.StructType:
		props := map[string]StructProperty{}
		for propname, propval := range t.Properties {
			propType, err := EvaluateTypeName(propval.PropType, env)
			if err != nil {
				errors.MakeError(env.filePath, t.Start.Line, t.Start.Column, t.End.Column, err.Error()).Display()
			}
			p := StructProperty{
				IsPrivate: propval.IsPrivate,
				Type: propType,
			}
			props[propname] = p
		}
		val = Struct{
			DataType: VALUE_TYPE(t.TypeName),
			StructName: node.UDTypeName,
			Elements: props,
		}
	default:
		typ , err := EvaluateTypeName(node.UDType, env)
		if err != nil {
			errors.MakeError(env.filePath, node.UDType.StartPos().Line, node.UDType.StartPos().Column, node.UDType.EndPos().Column, err.Error()).Display()
		}
		val = typ
	}

	typeVal := UserDefined{
		DataType: "user-defined",
		TypeDef: val,
	}

	env.DeclareType(node.UDTypeName, typeVal)
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
	case UserDefined:
		return getTypename(t.TypeDef)
	default:
		return t.DType()
	}
}

func MatchTypes(expected ValueTypeInterface, provided ValueTypeInterface, filePath string, line, start, end int) {

	fmt.Printf("%v\n", expected)

	expectedType := getTypename(expected)
	gotType := getTypename(provided)

	if expectedType != gotType {
		errors.MakeError(filePath, line, start, end, fmt.Sprintf("expected '%s', got '%s'", expectedType, gotType)).Display()
	}
}

func IsLValue(node ast.Node) bool {
	switch t := node.(type) {
	case ast.IdentifierExpr:
		return true
	case ast.ArrayIndexAccess:
		return IsLValue(t.Arrayvalue)
	default:
		return false
	}
}

func checkVariableAssignment(node ast.VarAssignmentExpr, env *TypeEnvironment) ValueTypeInterface {

	Assignee := node.Assignee
	valueToAssign := node.Value

	//varToAssign := node.Identifier
	expected := EvaluateTypesOfNode(Assignee, env)
	provided := EvaluateTypesOfNode(valueToAssign, env)
	
	MatchTypes(expected, provided, env.filePath, valueToAssign.StartPos().Line, valueToAssign.StartPos().Column, valueToAssign.EndPos().Column)

	var varName string

	if !IsLValue(Assignee) {
		errors.MakeError(env.filePath, Assignee.StartPos().Line, Assignee.StartPos().Column, Assignee.EndPos().Column, "invalid assignment expression. the assignee must be a lvalue").AddHint("lvalue is something that has a memory address\nFor example: variables.\nso you cannot assign values something which does not exist in memory as an independent identifier.", errors.TEXT_HINT).Display()
	}

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
