package typechecker

import (
	"fmt"
	"walrus/ast"
	"walrus/errgen"
)

func checkStructLiteral(structLit ast.StructLiteral, env *TypeEnvironment) ValueTypeInterface {

	sName := structLit.Identifier
	fmt.Printf("Checking struct literal %s\n", sName.Name)
	//check if defined
	declaredEnv, err := env.ResolveType(sName.Name)
	if err != nil {

		errgen.AddError(env.filePath, sName.StartPos().Line, sName.EndPos().Line, sName.StartPos().Column, sName.EndPos().Column, err.Error())
	}

	structType, ok := declaredEnv.types[sName.Name].(UserDefined).TypeDef.(Struct)
	if !ok {

		errgen.AddError(env.filePath, sName.StartPos().Line, sName.EndPos().Line, sName.StartPos().Column, sName.EndPos().Column, fmt.Sprintf("'%s' is not a struct", sName.Name))
	}

	// now we match the defined props with the provided props
	for propname, propval := range structLit.Properties {
		//check if the property is defined
		if _, ok := structType.StructScope.variables[propname]; !ok {

			errgen.AddError(env.filePath, propval.StartPos().Line, propval.EndPos().Line, propval.StartPos().Column, propval.EndPos().Column, fmt.Sprintf("property '%s' is not defined on struct '%s'", propname, sName.Name))
		}

		//check if the property type matches the defined type
		providedType := GetValueType(propval, env)
		expectedType := structType.StructScope.variables[propname].(StructProperty).Type

		err := MatchTypes(expectedType, providedType)
		if err != nil {

			errgen.AddError(env.filePath, propval.StartPos().Line, propval.EndPos().Line, propval.StartPos().Column, propval.EndPos().Column, err.Error())
		}
	}

	fmt.Printf("Checking if all required properties are provided\n")

	// check if any required property is missing
	for propname := range structType.StructScope.variables {
		// skip methods and 'this' variable
		if propname == "this" {
			continue
		}
		if _, ok := structType.StructScope.variables[propname].(StructMethod); ok {
			continue
		}
		fmt.Printf("Checking prop '%s'\n", propname)
		if _, ok := structLit.Properties[propname]; !ok {

			errgen.AddError(env.filePath, structLit.StartPos().Line, structLit.EndPos().Line, structLit.StartPos().Column, structLit.EndPos().Column, fmt.Sprintf("property '%s' is required on struct '%s'", propname, sName.Name))
		}
		fmt.Printf("Prop '%s' is provided\n", propname)
	}

	structValue := Struct{
		DataType:    STRUCT_TYPE,
		StructName:  sName.Name,
		StructScope: structType.StructScope,
	}

	return UserDefined{
		DataType: USER_DEFINED_TYPE,
		TypeName: sName.Name,
		TypeDef:  structValue,
	}
}

func checkPropertyAccess(expr ast.StructPropertyAccessExpr, env *TypeEnvironment) ValueTypeInterface {

	fmt.Printf("Checking property access\n")
	//show both the object and the property
	fmt.Printf("Object: %v\n", expr.Object)
	fmt.Printf("Property: %v\n", expr.Property)

	object := GetValueType(expr.Object, env)

	fmt.Printf("Object type: %T\n", object)

	prop := expr.Property

	lineStart := expr.Object.StartPos().Line
	lineEnd := expr.Object.EndPos().Line
	start := expr.Object.StartPos().Column
	end := expr.Object.EndPos().Column

	typeName := string(valueTypeInterfaceToString(object))
	// Resolve the struct type from the environment
	fmt.Printf("Resolving type %s\n", typeName)

	declaredEnv, err := env.ResolveType(typeName)
	if err != nil {

		errgen.AddError(env.filePath, lineStart, lineEnd, start, end, err.Error())
	}

	var structEnv TypeEnvironment

	//get the struct's environment
	switch t := declaredEnv.types[typeName].(UserDefined).TypeDef.(type) {
	case Struct:
		structEnv = t.StructScope
	case Interface:
		//prop must be a method
		return t.Methods[prop.Name]
	}

	// Check if the property exists on the struct
	if property, ok := structEnv.variables[prop.Name]; ok {
		isPrivate := false
		switch t := property.(type) {
		case StructMethod:
			isPrivate = t.IsPrivate
		case StructProperty:
			isPrivate = t.IsPrivate
		default:

			errgen.AddError(structEnv.filePath, prop.Start.Line, prop.End.Line, prop.Start.Column, prop.End.Column, fmt.Sprintf("'%s' is not a property or method", prop.Name))
		}

		if isPrivate {
			fmt.Printf("Scope type: %d, name: %s\n", env.scopeType, env.scopeName)
			//check the scope we are in
			if !env.IsInStructScope() {

				errgen.AddError(structEnv.filePath, prop.Start.Line, prop.End.Line, prop.Start.Column, prop.End.Column, fmt.Sprintf("cannot access private property '%s' from outside of the struct's scope", prop.Name))
			}
		}

		return property
	} else {
		// If the property is not found, check if it is a method

		errgen.AddError(structEnv.filePath, prop.Start.Line, prop.End.Line, prop.Start.Column, prop.End.Column, fmt.Sprintf("property '%s' does not exist on type '%s'", prop.Name, typeName))
	}

	// If neither property nor method is found, raise an error
	/*
		errgen.AddError(
			declaredEnv.filePath,
			prop.Start.Line,
			prop.End.Line,
			prop.Start.Column,
			prop.End.Column,
			fmt.Sprintf("property or method '%s' does not exist on type '%s'", prop.Name, typeName)).DisplayWithPanic()
	*/

	errgen.AddError(
		declaredEnv.filePath,
		prop.Start.Line,
		prop.End.Line,
		prop.Start.Column,
		prop.End.Column,
		fmt.Sprintf("property or method '%s' does not exist on type '%s'", prop.Name, typeName))

	return NewVoid()
}

func checkStructTypeDecl(name string, structType ast.StructType, env *TypeEnvironment) Struct {

	structEnv := NewTypeENV(env, STRUCT_SCOPE, name, env.filePath)

	//props := map[string]StructProperty{}
	for propname, propval := range structType.Properties {
		propType := EvaluateTypeName(propval.PropType, env)
		property := StructProperty{
			IsPrivate: propval.IsPrivate,
			Type:      propType,
		}
		//props[propname] = property
		//declare the property on the struct environment
		err := structEnv.DeclareVar(propname, property, false, false)
		if err != nil {

			errgen.AddError(env.filePath, propval.Prop.Start.Line, propval.Prop.End.Line, propval.Prop.Start.Column, propval.Prop.End.Column, err.Error())
		}
	}

	structTypeValue := Struct{
		DataType:    STRUCT_TYPE,
		StructName:  name,
		StructScope: *structEnv,
	}

	//declare 'this' variable to be used in the struct's methods
	err := structEnv.DeclareVar("this", structTypeValue, true, false)
	if err != nil {

		errgen.AddError(env.filePath, structType.Start.Line, structType.End.Line, structType.Start.Column, structType.End.Column, err.Error())
	}

	return structTypeValue
}
