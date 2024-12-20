package typechecker

import (
	//Standard packages
	"fmt"
	//Walrus packages
	"walrus/frontend/ast"
	"walrus/report"
	"walrus/utils"
)

func checkStructLiteral(structLit ast.StructLiteral, env *TypeEnvironment) Tc {

	sName := structLit.Identifier

	Type, err := getTypeDefinition(sName.Name) // need to get the most deep type
	if err != nil {
		report.Add(env.filePath, sName.StartPos().Line, sName.EndPos().Line, sName.StartPos().Column, sName.EndPos().Column, err.Error()).Level(report.NORMAL_ERROR)
	}

	structType, ok := Type.(Struct)

	if !ok {
		report.Add(env.filePath, sName.StartPos().Line, sName.EndPos().Line, sName.StartPos().Column, sName.EndPos().Column, fmt.Sprintf("'%s' is not a struct", sName.Name)).Level(report.CRITICAL_ERROR)
	}

	// now we match the defined props with the provided props
	checkPropsType(structType, structLit, env)
	// check if any required property is missing
	missingProps := checkMissingProps(structType, structLit)
	// if there are missing properties, we compose the error
	composeErrors(structLit, missingProps, env)

	structValue := Struct{
		DataType:    STRUCT_TYPE,
		StructName:  tcToString(Type),
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
	errMsg += report.TreeFormatString(missingProps...)
	report.Add(env.filePath, structLit.StartPos().Line, structLit.EndPos().Line, structLit.StartPos().Column, structLit.EndPos().Column, errMsg).Level(report.NORMAL_ERROR)
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

		if utils.None(toHaveProps, func(p ast.StructProp) bool {
			return p.Prop.Name == propname
		}) {
			missingProps = append(missingProps, fmt.Sprintf("missing property '%s'", propname))
		}
	}

	return missingProps
}

func checkPropsType(structType Struct, structLit ast.StructLiteral, env *TypeEnvironment) {
	for _, structProp := range structLit.Properties {
		//check if the property is defined
		if _, ok := structType.StructScope.variables[structProp.Prop.Name]; !ok {
			report.Add(env.filePath, structProp.Prop.Start.Line, structProp.Prop.End.Line, structProp.Prop.Start.Column, structProp.Prop.End.Column, fmt.Sprintf("property '%s' is not defined on struct '%s'", structProp.Prop.Name, structLit.Identifier.Name)).Level(report.CRITICAL_ERROR)
		}

		//check if the property type matches the defined type
		providedType := parseNodeValue(structProp.Value, env)

		expectedType := structType.StructScope.variables[structProp.Prop.Name].(StructProperty).Type

		err := validateTypeCompatibility(expectedType, providedType)
		if err != nil {
			report.Add(env.filePath, structProp.Prop.StartPos().Column, structProp.Value.EndPos().Line, structProp.Prop.StartPos().Column, structProp.Value.EndPos().Column, err.Error()).Level(report.NORMAL_ERROR)
		}
	}
}

func checkPropertyAccess(expr ast.StructPropertyAccessExpr, env *TypeEnvironment) Tc {

	fmt.Printf("Property Access: %s\n", expr.Property.Name)

	object := parseNodeValue(expr.Object, env)

	prop := expr.Property

	objName := tcToString(object)

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
		report.Add(env.filePath, prop.Start.Line, prop.End.Line, prop.Start.Column, prop.End.Column, fmt.Sprintf("interface '%s' does not have a method '%s'", t.InterfaceName, prop.Name)).Level(report.CRITICAL_ERROR)
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
			report.Add(env.filePath, prop.Start.Line, prop.End.Line, prop.Start.Column, prop.End.Column, fmt.Sprintf("'%s' is not a %s", prop.Name, propType)).Level(report.CRITICAL_ERROR)
		}

		if isPrivate {
			//check the scope we are in
			if !env.isInStructScope() {
				report.Add(env.filePath, prop.Start.Line, prop.End.Line, prop.Start.Column, prop.End.Column, fmt.Sprintf("cannot access private property '%s' from outside of the struct's scope", prop.Name)).Level(report.NORMAL_ERROR)
			}
		}
		return property
	}

	report.Add(env.filePath, prop.Start.Line, prop.End.Line, prop.Start.Column, prop.End.Column, fmt.Sprintf("'%s' does not exist on type '%s'", prop.Name, objName)).Level(report.CRITICAL_ERROR)

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
			report.Add(env.filePath, propval.PropType.StartPos().Line, propval.PropType.EndPos().Line, propval.PropType.StartPos().Column, propval.PropType.EndPos().Column, fmt.Sprintf("error declaring property '%s': %s", propval.Prop.Name, err.Error())).Level(report.CRITICAL_ERROR)
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
		report.Add(env.filePath, structType.StartPos().Line, structType.EndPos().Line, structType.StartPos().Column, structType.EndPos().Column, err.Error()).Level(report.CRITICAL_ERROR)
	}

	return structTypeValue
}
