package typechecker

import (
	"fmt"
	"walrus/ast"
	"walrus/errgen"
)

func checkStructLiteral(structLit ast.StructLiteral, env *TypeEnvironment) ValueTypeInterface {

	sName := structLit.Identifier
	//check if defined
	scope, err := env.ResolveType(sName.Name)
	if err != nil {
		errgen.MakeError(env.filePath, sName.StartPos().Line, sName.EndPos().Line, sName.StartPos().Column, sName.EndPos().Column, err.Error()).Display()
	}

	udType := scope.types[sName.Name].(UserDefined)

	structType, ok := udType.TypeDef.(Struct)
	if !ok {
		errgen.MakeError(env.filePath, sName.StartPos().Line, sName.EndPos().Line, sName.StartPos().Column, sName.EndPos().Column, fmt.Sprintf("'%s' is not a struct", sName.Name)).Display()
	}

	// now we match the defined props with the provided props
	// first check names on the provided value
	for propName, propValue := range structLit.Properties {
		//if prop exist
		if elem, ok := structType.Elements[propName]; !ok {
			errgen.MakeError(env.filePath, propValue.StartPos().Line, propValue.EndPos().Line, propValue.StartPos().Column, propValue.EndPos().Column, fmt.Sprintf("property '%s' does not exist on type '%s'", propName, sName.Name)).Display()
		} else {
			MatchTypes(elem.Type, CheckAST(propValue, env), env.filePath, propValue.StartPos().Line, propValue.EndPos().Line, propValue.StartPos().Column, propValue.EndPos().Column)
		}
	}

	//check if all required props are provided
	hint := ""
	//now check all from defined type
	for propName := range structType.Elements {
		if _, ok := structLit.Properties[propName]; !ok {
			errgen.MakeError(env.filePath, structLit.StartPos().Line, structLit.EndPos().Line, structLit.StartPos().Column, structLit.EndPos().Column, fmt.Sprintf("property '%s' is missing from type '%s'", propName, sName.Name)).AddHint(hint, errgen.TEXT_HINT).Display()
		}
	}

	return Struct{
		DataType:   STRUCT_TYPE,
		StructName: sName.Name,
		Elements:   structType.Elements,
	}
}

func checkPropertyAccess(expr ast.StructPropertyAccessExpr, env *TypeEnvironment) ValueTypeInterface {
	object := CheckAST(expr.Object, env)
	prop := expr.Property

	lineStart := expr.Object.StartPos().Line
	lineEnd := expr.Object.EndPos().Line
	start := expr.Object.StartPos().Column
	end := expr.Object.EndPos().Column

	// Ensure the object is a struct
	val, ok := object.(Struct)
	if !ok {
		errgen.MakeError(env.filePath, lineStart, lineEnd, start, end, "not an object").Display()
		return nil
	}

	// Resolve the struct type from the environment
	scope, err := env.ResolveType(val.StructName)
	if err != nil {
		errgen.MakeError(env.filePath, lineStart, lineEnd, start, end, err.Error()).Display()
		return nil
	}

	// Get the struct definition
	structDef := scope.types[val.StructName].(UserDefined).TypeDef.(Struct)

	// First check if the property exists in the struct's elements (fields)
	if valType, ok := structDef.Elements[prop.Name]; ok {
		if valType.IsPrivate {
			errgen.MakeError(
				env.filePath, 
				prop.Start.Line, 
				prop.End.Line, 
				prop.Start.Column, 
				prop.End.Column,
				fmt.Sprintf("cannot access private property '%s'", prop.Name,
			)).Display()
		}
		return valType.Type
	}

	// If not found in elements, check if it exists as a method
	if method, ok := structDef.Methods[prop.Name]; ok {
		if method.IsPrivate {
			errgen.MakeError(
				env.filePath, 
				prop.Start.Line, 
				prop.End.Line, 
				prop.Start.Column, 
				prop.End.Column,
				fmt.Sprintf("cannot access private method '%s'", prop.Name,
			)).Display()
		}
		return method
	}

	// If neither property nor method is found, raise an error
	errgen.MakeError(
		env.filePath, 
		prop.Start.Line, 
		prop.End.Line, 
		prop.Start.Column, 
		prop.End.Column,
		fmt.Sprintf("property or method '%s' does not exist on type '%s'", prop.Name, val.StructName,
	)).Display()

	return nil
}
