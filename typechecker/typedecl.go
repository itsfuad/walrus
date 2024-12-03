package typechecker

import (
	"fmt"
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

	err := DeclareType(node.UDTypeName, typeVal)
	if err != nil {
		fmt.Printf("cannot declare type %s: %s\n", node.UDTypeName, err.Error())
		errgen.AddError(env.filePath, node.Start.Line, node.End.Line, node.Start.Column, node.End.Column, err.Error())
	}

	utils.GREEN.Print("Declared Type ")
	utils.PURPLE.Println(node.UDTypeName)

	return val
}
