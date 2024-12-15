package analyzer

import (
	"walrus/ast"
	"walrus/errgen"
	"walrus/utils"
)

func checkTypeDeclaration(node ast.TypeDeclStmt, env *TypeEnvironment) TcValue {

	typeName := node.UDType

	utils.BLUE.Print("declaring type ")
	utils.PURPLE.Println(node.UDTypeName)

	var val TcValue

	switch t := typeName.(type) {
	case ast.StructType:
		val = checkStructTypeDecl(node.UDTypeName, t, env)
	case ast.InterfaceType:
		val = checkInterfaceTypeDecl(node.UDTypeName, t, env)
	default:
		val = evaluateTypeName(typeName, env)
	}

	typeVal := UserDefined{
		DataType: USER_DEFINED_TYPE,
		TypeName: node.UDTypeName,
		TypeDef:  val,
	}

	err := declareType(node.UDTypeName, typeVal)
	if err != nil {
		errgen.Add(env.filePath, node.Start.Line, node.End.Line, node.Start.Column, node.End.Column, err.Error()).Level(errgen.NORMAL)
	}

	utils.GREEN.Print("Declared Type ")
	utils.PURPLE.Println(node.UDTypeName)

	return val
}
