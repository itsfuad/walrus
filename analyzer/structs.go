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
	checkPropsType(structType, structLit, env)
	// check if any required property is missing
	missingProps := checkMissingProps(structType, structLit)
	// if there are missing properties, we compose the error
	composeErrors(structLit, missingProps, env)

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

func composeErrors(structLit ast.StructLiteral, missingProps []string, env *TypeEnvironment) {
	if len(missingProps) == 0 {
		return
	}
	errMsg := fmt.Sprintf("incomplete struct literal for struct '%s'\n", structLit.Identifier.Name)
	errMsg += errgen.TreeFormatString(missingProps...)
	errgen.AddError(env.filePath, structLit.StartPos().Line, structLit.EndPos().Line, structLit.StartPos().Column, structLit.EndPos().Column, errMsg).ErrorLevel(errgen.NORMAL)
}

func checkMissingProps(structType Struct, structLit ast.StructLiteral) []string {
	//check for missing properties
	missingProps := []string{}

	for propname := range structType.StructScope.variables {
		// skip methods and 'this' variable
		if propname == "this" {
			continue
		}
		if _, ok := structType.StructScope.variables[propname].(StructMethod); ok {
			continue
		}
		//find the missing properties
		toHaveProps := structLit.Properties //slice of provided properties
		found := false
		for _, prop := range toHaveProps {
			if prop.Prop.Name == propname {
				found = true
				break
			}
		}
		if !found {
			missingProps = append(missingProps, fmt.Sprintf("missing property '%s'", propname))
		}
	}

	return missingProps
}

func checkPropsType(structType Struct, structLit ast.StructLiteral, env *TypeEnvironment) {
	for _, structProp := range structLit.Properties {
		//check if the property is defined
		if _, ok := structType.StructScope.variables[structProp.Prop.Name]; !ok {
			errgen.AddError(env.filePath, structProp.Prop.Start.Line, structProp.Prop.End.Line, structProp.Prop.Start.Column, structProp.Prop.End.Column, fmt.Sprintf("property '%s' is not defined on struct '%s'", structProp.Prop.Name, structLit.Identifier.Name)).ErrorLevel(errgen.CRITICAL)
		}

		//check if the property type matches the defined type
		providedType := parseNodeValue(structProp.Value, env)

		expectedType := structType.StructScope.variables[structProp.Prop.Name].(StructProperty).Type

		err := matchTypes(expectedType, providedType)
		if err != nil {
			errgen.AddError(env.filePath, structProp.Prop.StartPos().Column, structProp.Value.EndPos().Line, structProp.Prop.StartPos().Column, structProp.Value.EndPos().Column, err.Error()).ErrorLevel(errgen.NORMAL)
		}
	}
}

func checkPropertyAccess(expr ast.StructPropertyAccessExpr, env *TypeEnvironment) TcValue {

	fmt.Printf("Property Access: %s\n", expr.Property.Name)

	object := parseNodeValue(expr.Object, env)

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
		for _, method := range t.Methods {
			if method.Name == prop.Name {
				return method.Method
			}
		}
		errgen.AddError(env.filePath, prop.Start.Line, prop.End.Line, prop.Start.Column, prop.End.Column, fmt.Sprintf("interface '%s' does not have a method '%s'", t.InterfaceName, prop.Name)).ErrorLevel(errgen.CRITICAL)
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

	for _, propval := range structType.Properties {
		propType := evaluateTypeName(propval.PropType, env)
		property := StructProperty{
			IsPrivate: propval.IsPrivate,
			Type:      propType,
		}
		//declare the property on the struct environment
		err := structEnv.declareVar(propval.Prop.Name, property, false, false)
		if err != nil {
			errgen.AddError(env.filePath, propval.PropType.StartPos().Line, propval.PropType.EndPos().Line, propval.PropType.StartPos().Column, propval.PropType.EndPos().Column, fmt.Sprintf("error declaring property '%s': %s", propval.Prop.Name, err.Error())).ErrorLevel(errgen.CRITICAL)
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
