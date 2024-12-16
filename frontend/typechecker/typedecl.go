package typechecker

import (
	"walrus/errgen"
	"walrus/frontend/ast"
	"walrus/utils"
)

func checkTypeDeclaration(node ast.TypeDeclStmt, env *TypeEnvironment) ExprType {

	typeName := node.UDTypeValue

	utils.BLUE.Print("declaring type ")
	utils.PURPLE.Println(node.UDTypeName.Name)

	//if typename is small case
	if !utils.IsCapitalized(node.UDTypeName.Name) {
		errgen.Add(env.filePath, node.UDTypeName.Start.Line, node.UDTypeName.End.Line, node.UDTypeName.Start.Column, node.UDTypeName.End.Column, "Type name should be capitalized").Hint("Make the first letter uppercase").Level(errgen.WARNING)
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
		errgen.Add(env.filePath, node.Start.Line, node.End.Line, node.Start.Column, node.End.Column, err.Error()).Level(errgen.NORMAL_ERROR)
	}

	utils.GREEN.Print("Declared Type ")
	utils.PURPLE.Println(node.UDTypeName.Name)

	return val
}
