package typechecker

import (
	//Walrus packages
	"walrus/compiler/internal/ast"
	"walrus/compiler/report"
)

func checkMapLiteral(node ast.MapLiteral, env *TypeEnvironment) Tc {

	//get the map definitions
	evaluatedMapType := evaluateTypeName(node.MapType, env)

	//check the key value pairs
	for _, value := range node.Values {
		keyType := parseNodeValue(value.Key, env)
		valueType := parseNodeValue(value.Value, env)

		err := validateTypeCompatibility(evaluatedMapType.(Map).KeyType, keyType)
		if err != nil {
			report.Add(env.filePath, value.Key.StartPos().Line, value.Key.EndPos().Line, value.Key.StartPos().Column, value.Key.EndPos().Column, "incorrect map key. "+err.Error()).SetLevel(report.NORMAL_ERROR)
		}
		err = validateTypeCompatibility(evaluatedMapType.(Map).ValueType, valueType)
		if err != nil {
			report.Add(env.filePath, value.Value.StartPos().Line, value.Value.EndPos().Line, value.Value.StartPos().Column, value.Value.EndPos().Column, "incorrect map value. "+err.Error()).SetLevel(report.NORMAL_ERROR)
		}
	}

	return evaluatedMapType
}
