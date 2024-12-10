package analyzer

import (
	"fmt"
	"walrus/ast"
	"walrus/errgen"
)

func checkStructLiteral(structLit ast.StructLiteral, env *TypeEnvironment) TcValue {

	sName := structLit.Identifier

	Type, err := getTypeDefinition(sName.Name) // need to get the most deep type
	if err != nil {
		errgen.AddError(env.filePath, sName.StartPos().Line, sName.EndPos().Line, sName.StartPos().Column, sName.EndPos().Column, err.Error()).ErrorLevel(errgen.NORMAL)
	}

	structType, ok := Type.(Struct)

	if !ok {
		errgen.AddError(env.filePath, sName.StartPos().Line, sName.EndPos().Line, sName.StartPos().Column, sName.EndPos().Column, fmt.Sprintf("'%s' is not a struct", sName.Name)).ErrorLevel(errgen.CRITICAL)
	}

	// now we match the defined props with the provided props
	for propname, prop := range structLit.Properties {
		//check if the property is defined
		if _, ok := structType.StructScope.variables[propname]; !ok {
			errgen.AddError(env.filePath, prop.Prop.Start.Line, prop.Prop.End.Line, prop.Prop.Start.Column, prop.Prop.End.Column, fmt.Sprintf("property '%s' is not defined on struct '%s'", propname, sName.Name)).ErrorLevel(errgen.CRITICAL)
		}

		//check if the property type matches the defined type
		providedType := CheckAST(prop.Value, env)

		expectedType := structType.StructScope.variables[propname].(StructProperty).Type

		err := matchTypes(expectedType, providedType)
		if err != nil {
			errgen.AddError(env.filePath, prop.Prop.StartPos().Column, prop.Value.EndPos().Line, prop.Prop.StartPos().Column, prop.Value.EndPos().Column, err.Error()).ErrorLevel(errgen.NORMAL)
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
			errgen.AddError(env.filePath, structLit.StartPos().Line, structLit.EndPos().Line, structLit.StartPos().Column, structLit.EndPos().Column, fmt.Sprintf("property '%s' is missing on @%s", propname, sName.Name)).ErrorLevel(errgen.NORMAL)
		}
	}

	structValue := Struct{
		DataType:    STRUCT_TYPE,
		StructName:  tcValueToString(Type),
		StructScope: structType.StructScope,
	}

	return UserDefined{
		DataType: USER_DEFINED_TYPE,
		TypeName: sName.Name,
		TypeDef:  structValue,
	}
}

func checkPropertyAccess(expr ast.StructPropertyAccessExpr, env *TypeEnvironment) TcValue {

	fmt.Printf("Property Access: %s\n", expr.Property.Name)

	object := CheckAST(expr.Object, env)

	prop := expr.Property

	fmt.Printf("Obj type: %T\n", object)

	objName := tcValueToString(object)

	var structEnv TypeEnvironment

	//get the struct's environment
	switch t := object.(type) {
	case Struct:
		structEnv = t.StructScope
	case Interface:
		//prop must be a method
		if _, ok := t.Methods[prop.Name]; !ok {
			errgen.AddError(env.filePath, prop.Start.Line, prop.End.Line, prop.Start.Column, prop.End.Column, fmt.Sprintf("interface '%s' does not have a method '%s'", t.InterfaceName, prop.Name)).ErrorLevel(errgen.CRITICAL)
			return NewVoid() // unreachable but needed to avoid accidental nil pointer dereference
		}
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
			errgen.AddError(env.filePath, prop.Start.Line, prop.End.Line, prop.Start.Column, prop.End.Column, fmt.Sprintf("'%s' is not a %s", prop.Name, propType)).ErrorLevel(errgen.CRITICAL)
		}

		if isPrivate {
			//check the scope we are in
			if !env.isInStructScope() {
				errgen.AddError(env.filePath, prop.Start.Line, prop.End.Line, prop.Start.Column, prop.End.Column, fmt.Sprintf("cannot access private property '%s' from outside of the struct's scope", prop.Name)).ErrorLevel(errgen.NORMAL)
			}
		}
		return property
	}

	errgen.AddError(env.filePath, prop.Start.Line, prop.End.Line, prop.Start.Column, prop.End.Column, fmt.Sprintf("'%s' does not exist on type '%s'", prop.Name, objName)).ErrorLevel(errgen.CRITICAL)

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
		//declare the property on the struct environment
		err := structEnv.declareVar(propname, property, false, false)
		if err != nil {
			errgen.AddError(env.filePath, propval.PropType.StartPos().Line, propval.PropType.EndPos().Line, propval.PropType.StartPos().Column, propval.PropType.EndPos().Column, fmt.Sprintf("error declaring property '%s': %s", propname, err.Error())).ErrorLevel(errgen.CRITICAL)
		}
	}

	structTypeValue := Struct{
		DataType:    STRUCT_TYPE,
		StructName:  name,
		StructScope: *structEnv,
	}

	//declare 'this' variable to be used in the struct's methods
	err := structEnv.declareVar("this", structTypeValue, true, false)
	if err != nil {
		errgen.AddError(env.filePath, structType.StartPos().Line, structType.EndPos().Line, structType.StartPos().Column, structType.EndPos().Column, err.Error()).ErrorLevel(errgen.CRITICAL)
	}

	return structTypeValue
}
