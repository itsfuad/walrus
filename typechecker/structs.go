package typechecker

import (
	"fmt"
	"walrus/ast"
	"walrus/errgen"
)

func checkStructLiteral(structLit ast.StructLiteral, env *TypeEnvironment) ValueTypeInterface {

	sName := structLit.Identifier

	Type, err := getTypeDefinition(sName.Name)
	if err != nil {
		errgen.AddError(env.filePath, sName.StartPos().Line, sName.EndPos().Line, sName.StartPos().Column, sName.EndPos().Column, fmt.Sprintf("'%s' is not a struct", sName.Name))
	}

	structType := Type.(Struct)

	// now we match the defined props with the provided props
	for propname, propval := range structLit.Properties {
		//check if the property is defined
		if _, ok := structType.StructScope.variables[propname]; !ok {

			errgen.AddError(env.filePath, propval.StartPos().Line, propval.EndPos().Line, propval.StartPos().Column, propval.EndPos().Column, fmt.Sprintf("property '%s' is not defined on struct '%s'", propname, sName.Name))
		}

		//check if the property type matches the defined type
		providedType := nodeType(propval, env)
		expectedType := structType.StructScope.variables[propname].(StructProperty).Type

		err := matchTypes(expectedType, providedType)
		if err != nil {

			errgen.AddError(env.filePath, propval.StartPos().Line, propval.EndPos().Line, propval.StartPos().Column, propval.EndPos().Column, err.Error())
		}
	}

	// check if any required property is missing
	for propname := range structType.StructScope.variables {
		// skip methods and 'this' variable
		if propname == "this" {
			continue
		}
		if _, ok := structType.StructScope.variables[propname].(StructMethod); ok {
			continue
		}
		if _, ok := structLit.Properties[propname]; !ok {

			errgen.AddError(env.filePath, structLit.StartPos().Line, structLit.EndPos().Line, structLit.StartPos().Column, structLit.EndPos().Column, fmt.Sprintf("property '%s' is required on struct '%s'", propname, sName.Name))
		}
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

	object := nodeType(expr.Object, env)

	prop := expr.Property

	lineStart := expr.Object.StartPos().Line
	lineEnd := expr.Object.EndPos().Line
	start := expr.Object.StartPos().Column
	end := expr.Object.EndPos().Column

	typeName := string(valueTypeInterfaceToString(object))

	Type, err := getTypeDefinition(typeName)
	if err != nil {
		errgen.AddError(env.filePath, lineStart, lineEnd, start, end, err.Error())
	}

	var structEnv TypeEnvironment

	//get the struct's environment
	switch t := Type.(type) {
	case Struct:
		structEnv = t.StructScope
	case Interface:
		//prop must be a method
		return t.Methods[prop.Name]
	}

	propType := ""

	// Check if the property exists on the struct
	if property, ok := structEnv.variables[prop.Name]; ok {
		isPrivate := false
		switch t := property.(type) {
		case StructMethod:
			propType = "method"
			isPrivate = t.IsPrivate
		case StructProperty:
			propType = "property"
			isPrivate = t.IsPrivate
		default:
			errgen.AddError(env.filePath, prop.Start.Line, prop.End.Line, prop.Start.Column, prop.End.Column, fmt.Sprintf("'%s' is not a %s", prop.Name, propType)).DisplayWithPanic()
		}

		if isPrivate {
			//check the scope we are in
			if !env.IsInStructScope() {
				errgen.AddError(env.filePath, prop.Start.Line, prop.End.Line, prop.Start.Column, prop.End.Column, fmt.Sprintf("cannot access private property '%s' from outside of the struct's scope", prop.Name))
			}
		}
		return property
	}

	errgen.AddError(env.filePath, prop.Start.Line, prop.End.Line, prop.Start.Column, prop.End.Column, fmt.Sprintf("'%s' does not exist on type '%s'", prop.Name, typeName)).DisplayWithPanic()

	return NewVoid()
}

func checkStructTypeDecl(name string, structType ast.StructType, env *TypeEnvironment) Struct {

	structEnv := NewTypeENV(env, STRUCT_SCOPE, name, env.filePath)

	for propname, propval := range structType.Properties {
		propType := evaluateTypeName(propval.PropType, env)
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
