package typechecker

import (
	//Standard packages
	"fmt"
	//Walrus packages
	"walrus/compiler/colors"
	"walrus/compiler/internal/ast"
	"walrus/compiler/internal/utils"
	"walrus/compiler/report"
)

func checkAnnonymousStructLiteral(structLit ast.StructLiteral, env *TypeEnvironment) Tc {

	var structType Struct

	structEnv := NewTypeENV(env, STRUCT_SCOPE, "", env.filePath)

	//make the struct name "struct { prop1: type1, prop2: type2, ... }"
	sname := "struct { "

	for i, propval := range structLit.Properties {
		propType := parseNodeValue(propval.Value, env)
		property := StructProperty{
			IsPrivate: false,
			Type:      propType,
		}
		//declare the property on the struct environment
		err := structEnv.declareVar(propval.Prop.Name, property, false, false)
		if err != nil {
			report.Add(env.filePath, propval.Prop.StartPos().Line, propval.Prop.EndPos().Line, propval.Prop.StartPos().Column, propval.Prop.EndPos().Column, fmt.Sprintf("error declaring property '%s': %s", propval.Prop.Name, err.Error())).SetLevel(report.CRITICAL_ERROR)
		}

		sname += fmt.Sprintf("%s: %s", propval.Prop.Name, tcToString(propType))
		if i < len(structLit.Properties)-1 {
			sname += ", "
		}
	}

	sname += " }"

	structType = Struct{
		DataType:    STRUCT_TYPE,
		StructName:  sname,
		StructScope: *structEnv,
	}

	return structType
}

func checkStructLiteral(structLit ast.StructLiteral, env *TypeEnvironment) Tc {

	sName := structLit.Identifier

	colors.BLUE.Printf("Struct Literal: '%s'\n", sName.Name)

	if sName.Name == "" {
		//annonymous struct
		return checkAnnonymousStructLiteral(structLit, env)
	}

	Type, err := getTypeDefinition(sName.Name) // need to get the most deep type
	if err != nil {
		report.Add(env.filePath, sName.StartPos().Line, sName.EndPos().Line, sName.StartPos().Column, sName.EndPos().Column, err.Error()).SetLevel(report.NORMAL_ERROR)
	}

	structType, ok := Type.(Struct)

	if !ok {
		report.Add(env.filePath, sName.StartPos().Line, sName.EndPos().Line, sName.StartPos().Column, sName.EndPos().Column, fmt.Sprintf("'%s' is not a struct", sName.Name)).SetLevel(report.CRITICAL_ERROR)
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
	report.Add(env.filePath, structLit.StartPos().Line, structLit.EndPos().Line, structLit.StartPos().Column, structLit.EndPos().Column, errMsg).SetLevel(report.NORMAL_ERROR)
}

func checkMissingProps(structType Struct, structLit ast.StructLiteral) []string {
	//check for missing properties
	missingProps := []string{}

	for propname := range structType.StructScope.variables {
		//skip methods
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
			report.Add(env.filePath, structProp.Prop.Start.Line, structProp.Prop.End.Line, structProp.Prop.Start.Column, structProp.Prop.End.Column, fmt.Sprintf("property '%s' is not defined on struct '%s'", structProp.Prop.Name, structLit.Identifier.Name)).SetLevel(report.CRITICAL_ERROR)
		}

		//check if the property type matches the defined type
		providedType := parseNodeValue(structProp.Value, env)

		expectedType := structType.StructScope.variables[structProp.Prop.Name].(StructProperty).Type

		err := validateTypeCompatibility(expectedType, providedType)
		if err != nil {
			report.Add(env.filePath, structProp.Prop.StartPos().Column, structProp.Value.EndPos().Line, structProp.Prop.StartPos().Column, structProp.Value.EndPos().Column, err.Error()).SetLevel(report.NORMAL_ERROR)
		}
	}
}

func getObject(expr ast.StructPropertyAccessExpr, env *TypeEnvironment) Tc {
	// if obj is 'this' then we return the struct type
	if ok := expr.Object.(ast.IdentifierExpr); ok.Name == "this" {
		obj, err := env.getStructType()
		if err != nil {
			report.Add(env.filePath, expr.Object.StartPos().Line, expr.Object.EndPos().Line, expr.Object.StartPos().Column, expr.Object.EndPos().Column, err.Error()).SetLevel(report.CRITICAL_ERROR)
		}
		return obj
	} else {
		return parseNodeValue(expr.Object, env)
	}
}

func checkPropertyAccess(expr ast.StructPropertyAccessExpr, env *TypeEnvironment) Tc {

	object := getObject(expr, env)

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
		report.Add(env.filePath, prop.Start.Line, prop.End.Line, prop.Start.Column, prop.End.Column, fmt.Sprintf("interface '%s' does not have a method '%s'", t.InterfaceName, prop.Name)).SetLevel(report.CRITICAL_ERROR)
	}

	propType := ""
	var propValue Tc

	// Check if the property exists on the struct
	if property, ok := structEnv.variables[prop.Name]; ok {
		isPrivate := false
		switch t := property.(type) {
		case StructMethod:
			propType = "method"
			isPrivate = t.IsPrivate
			propValue = t.Fn
		case StructProperty:
			propType = "property"
			isPrivate = t.IsPrivate
			propValue = t.Type
		default:
			report.Add(env.filePath, prop.Start.Line, prop.End.Line, prop.Start.Column, prop.End.Column, fmt.Sprintf("'%s' is not a %s", prop.Name, propType)).SetLevel(report.CRITICAL_ERROR)
		}

		if isPrivate {
			//check the scope we are in
			if !env.isInStructScope() {
				report.Add(env.filePath, prop.Start.Line, prop.End.Line, prop.Start.Column, prop.End.Column, fmt.Sprintf("cannot access private property '%s' from outside of the struct's scope", prop.Name)).SetLevel(report.NORMAL_ERROR)
			}
		}
		return propValue
	}

	report.Add(env.filePath, prop.Start.Line, prop.End.Line, prop.Start.Column, prop.End.Column, fmt.Sprintf("'%s' does not exist on type '%s'", prop.Name, objName)).SetLevel(report.CRITICAL_ERROR)

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
			report.Add(env.filePath, propval.PropType.StartPos().Line, propval.PropType.EndPos().Line, propval.PropType.StartPos().Column, propval.PropType.EndPos().Column, fmt.Sprintf("error declaring property '%s': %s", propval.Prop.Name, err.Error())).SetLevel(report.CRITICAL_ERROR)
		}
	}

	structTypeValue := Struct{
		DataType:    STRUCT_TYPE,
		StructName:  name,
		StructScope: *structEnv,
	}

	return structTypeValue
}
