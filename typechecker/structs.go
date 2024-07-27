package typechecker

import (
	"fmt"
	"walrus/ast"
	"walrus/errgen"
)

func checkStructLiteral(structLit ast.StructLiteral, env *TypeEnvironment) ValueTypeInterface {

	sName := structLit.Name
	//check if defined
	scope, err := env.ResolveType(sName.Name)
	if err != nil {
		errgen.MakeError(env.filePath, sName.StartPos().Line, sName.EndPos().Line, sName.StartPos().Column, sName.EndPos().Column, fmt.Sprintf("type '%s' is not defined", sName.Name)).Display()
	}

	udType := scope.types[sName.Name].(UserDefined)

	structType := udType.TypeDef.(Struct)

	// now we match the defined props with the provided props
	// first check names on the provided value
	for propName, propValue := range structLit.Properties {
		//if prop exist
		if elem, ok := structType.Elements[propName]; !ok {
			errgen.MakeError(env.filePath, propValue.StartPos().Line, propValue.EndPos().Line,  propValue.StartPos().Column, propValue.EndPos().Column, fmt.Sprintf("property '%s' does not exists on type '%s'", propName, sName.Name)).Display()
		} else {
			MatchTypes(elem.Type, CheckAST(propValue, env), env.filePath, propValue.StartPos().Line, propValue.EndPos().Line, propValue.StartPos().Column, propValue.EndPos().Column)
		}
	}

	//now check all from defined type
	for propName := range structType.Elements {
		if _, ok := structLit.Properties[propName]; !ok {
			//if does not exists
			errgen.MakeError(env.filePath, structLit.StartPos().Line, structLit.EndPos().Line, structLit.StartPos().Column, structLit.EndPos().Column, fmt.Sprintf("property '%s' is missing from type '%s'", propName, sName.Name)).Display()
		}
	}

	return Struct{
		DataType:   STRUCT_TYPE,
		StructName: sName.Name,
		Elements:   structType.Elements,
	}
}

func checkProperty(expr ast.StructPropertyAccessExpr, env *TypeEnvironment) ValueTypeInterface {
	object := CheckAST(expr.Object, env)
	prop := expr.Property

	lineStart := expr.Object.StartPos().Line
	lineEnd := expr.Object.EndPos().Line
	start := expr.Object.StartPos().Column
	end := expr.Object.EndPos().Column

	val, ok := object.(Struct)

	if !ok {
		errgen.MakeError(env.filePath, lineStart, lineEnd, start, end, "not an object").Display()
	}

	//find if object exists or not
	scope, err := env.ResolveType(val.StructName)
	if err != nil {
		errgen.MakeError(env.filePath, lineStart, lineEnd, start, end, err.Error()).Display()
	}

	//get the type of the prop
	structOnEnv := scope.types[val.StructName]

	elems := structOnEnv.(UserDefined).TypeDef.(Struct).Elements

	valType, ok := elems[prop.Name]

	if !ok {
		errgen.MakeError(env.filePath, prop.Start.Line, prop.End.Line, prop.Start.Column, prop.End.Column, fmt.Sprintf("property '%s' does not exist on type '%s'", prop.Name, val.StructName)).Display()
	}

	if elems[prop.Name].IsPrivate {
		errgen.MakeError(env.filePath, prop.Start.Line, prop.End.Line, prop.Start.Column, prop.End.Column, fmt.Sprintf("cannot access private property '%s'", prop.Name)).Display()
	}

	return valType.Type
}
