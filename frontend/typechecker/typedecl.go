package typechecker

import (
	"walrus/frontend/ast"
	"walrus/report"
	"walrus/utils"
)

func checkTypeDeclaration(node ast.TypeDeclStmt, env *TypeEnvironment) ExprType {

	typeName := node.UDTypeValue

	utils.BLUE.Print("declaring type ")
	utils.PURPLE.Println(node.UDTypeName.Name)

	//if typename is small case
	if !utils.IsCapitalized(node.UDTypeName.Name) {
		report.Add(env.filePath, node.UDTypeName.Start.Line, node.UDTypeName.Start.Line, node.UDTypeName.Start.Column, node.UDTypeName.Start.Column, "User defined type name should be capitalized").Hint("Make the first letter uppercase").Level(report.INFO)
	}

	var val ExprType

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
		report.Add(env.filePath, node.Start.Line, node.End.Line, node.Start.Column, node.End.Column, err.Error()).Level(report.NORMAL_ERROR)
	}

	utils.GREEN.Print("Declared Type ")
	utils.PURPLE.Println(node.UDTypeName.Name)

	return val
}
