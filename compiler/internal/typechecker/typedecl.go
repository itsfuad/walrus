package typechecker

import (
	"walrus/compiler/colors"
	"walrus/compiler/internal/ast"
	"walrus/compiler/internal/utils"
	"walrus/compiler/report"
)

func checkTypeDeclaration(node ast.TypeDeclStmt, env *TypeEnvironment) Tc {

	typeName := node.UDTypeValue

	colors.BLUE.Print("declaring type ")
	colors.PURPLE.Println(node.UDTypeName.Name)

	//if typename is small case
	if !utils.IsCapitalized(node.UDTypeName.Name) {
		report.Add(env.filePath, node.UDTypeName.Start.Line, node.UDTypeName.Start.Line, node.UDTypeName.Start.Column, node.UDTypeName.Start.Column, "User defined type name should be capitalized").Hint("Make the first letter uppercase").SetLevel(report.INFO)
	}

	var val Tc

	switch t := typeName.(type) {
	case ast.StructType:
		val = checkStructTypeDecl(node.UDTypeName.Name, t, env)
	case ast.InterfaceType:
		val = checkInterfaceTypeDecl(node.UDTypeName.Name, t, env)
	default:
		val = evaluateTypeName(typeName, env)
	}

	typeVal := UserDefined{
		DataType: USER_DEFINED_TYPE,
		TypeName: node.UDTypeName.Name,
		TypeDef:  val,
	}

	err := declareType(node.UDTypeName.Name, typeVal)
	if err != nil {
		report.Add(env.filePath, node.Start.Line, node.End.Line, node.Start.Column, node.End.Column, err.Error()).SetLevel(report.NORMAL_ERROR)
	}

	colors.GREEN.Print("Declared Type ")
	colors.PURPLE.Println(node.UDTypeName.Name)

	return val
}
