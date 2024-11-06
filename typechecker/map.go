package typechecker

import (
	"walrus/ast"
	"walrus/errgen"
)

func checkMapLiteral(node ast.MapLiteral, env *TypeEnvironment) ValueTypeInterface {

	//get the map definitions

	evaluatedMapType := EvaluateTypeName(node.MapType, env)

	//check if the map type is valid
	if evaluatedMapType.DType() != MAP_TYPE {
		errgen.AddError(env.filePath, node.StartPos().Line, node.EndPos().Line, node.StartPos().Column, node.EndPos().Column, "invalid map type")
	}

	//check the key value pairs
	for _, value := range node.Values {
		keyType := nodeType(value.Key, env)
		valueType := nodeType(value.Value, env)

		err := MatchTypes(evaluatedMapType.(Map).KeyType, keyType)
		if err != nil {
			errgen.AddError(env.filePath, value.Key.StartPos().Line, value.Key.EndPos().Line, value.Key.StartPos().Column, value.Key.EndPos().Column, "incorrect map key. "+err.Error())
		}
		err = MatchTypes(evaluatedMapType.(Map).ValueType, valueType)
		if err != nil {
			errgen.AddError(env.filePath, value.Value.StartPos().Line, value.Value.EndPos().Line, value.Value.StartPos().Column, value.Value.EndPos().Column, "incorrect map value. "+err.Error())
		}
	}

	return evaluatedMapType
}
