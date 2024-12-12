package analyzer

import (
	"walrus/ast"
	"walrus/errgen"
)

func checkMapLiteral(node ast.MapLiteral, env *TypeEnvironment) TcValue {

	//get the map definitions
	evaluatedMapType := evaluateTypeName(node.MapType, env)

	//check the key value pairs
	for _, value := range node.Values {
		keyType := parseNodeValue(value.Key, env)
		valueType := parseNodeValue(value.Value, env)

		err := matchTypes(evaluatedMapType.(Map).KeyType, keyType)
		if err != nil {
			errgen.AddError(env.filePath, value.Key.StartPos().Line, value.Key.EndPos().Line, value.Key.StartPos().Column, value.Key.EndPos().Column, "incorrect map key. "+err.Error()).ErrorLevel(errgen.NORMAL)
		}
		err = matchTypes(evaluatedMapType.(Map).ValueType, valueType)
		if err != nil {
			errgen.AddError(env.filePath, value.Value.StartPos().Line, value.Value.EndPos().Line, value.Value.StartPos().Column, value.Value.EndPos().Column, "incorrect map value. "+err.Error()).ErrorLevel(errgen.NORMAL)
		}
	}

	return evaluatedMapType
}
