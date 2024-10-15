package typechecker

import (
	"fmt"
	"walrus/ast"
	"walrus/errgen"
)

func checkImplStmt(node ast.ImplStmt, env *TypeEnvironment) ValueTypeInterface {

	//if env is not global, throw error
	if env.scopeType != GLOBAL_SCOPE {
		errgen.MakeError(env.filePath, node.StartPos().Line, node.EndPos().Line, node.StartPos().Column, node.EndPos().Column, "impl statement must be in global scope").Display()
	}

	traitMethods := make(map[string]Fn)

	// if trait present, check if it exists
	if node.Trait != nil {
		//check if trait exists
		//convert to Identifier
		trait, ok := node.Trait.(ast.IdentifierExpr)
		if !ok {
			errgen.MakeError(env.filePath, node.Trait.StartPos().Line, node.Trait.EndPos().Line, node.Trait.StartPos().Column, node.Trait.EndPos().Column, "trait name must be an identifier").Display()
		}
		
		traitEnv, err := env.ResolveTrait(trait.Name)
		if err != nil {
			errgen.MakeError(env.filePath, trait.StartPos().Line, trait.EndPos().Line, trait.StartPos().Column, trait.EndPos().Column, err.Error()).Display()
		}

		traitMethods = traitEnv.traits[trait.Name].Methods
	} else {
		// if no trait then the impl must be for a struct
		structNode := node.ImplFor
		//check if struct exists
		_, err := env.ResolveStruct(structNode.Name)
		if err != nil {
			errgen.MakeError(env.filePath, structNode.StartPos().Line, structNode.EndPos().Line, structNode.StartPos().Column, structNode.EndPos().Column, err.Error()).Display()
		}
		structValue := env.types[structNode.Name].(UserDefined).TypeDef.(Struct)
		fmt.Println(structValue.StructName)
		for name, method := range node.Methods {
			if _, ok := structValue.Elements[name]; ok {
				errgen.MakeError(env.filePath, method.StartPos().Line, method.EndPos().Line, method.StartPos().Column, method.EndPos().Column, fmt.Sprintf("'%s' is already defined in struct %s", name, structValue.StructName)).Display()
			} else {
				// if not defined, add to struct
				fnParams := []FnParam{}
				for _, param := range method.Params {
					fnParams = append(fnParams, FnParam{
						Name:       param.Identifier.Name,
						IsOptional: param.IsOptional,
						Type:       EvaluateTypeName(param.Type, env),
					})
				}
				structValue.Methods[name] = StructMethod{
					IsPrivate: 	method.IsPrivate,
					Fn: Fn{
						DataType: FUNCTION_TYPE,
						Params:   fnParams,
						Returns:  EvaluateTypeName(method.ReturnType, env),
						FunctionScope: *NewTypeENV(env, FUNCTION_SCOPE, name, env.filePath),
					},
				}
				fmt.Printf("Added %s to struct %s\n", name, structValue.StructName)
			}
		}
		//update struct in env
		env.types[structNode.Name] = UserDefined{
			DataType: USER_DEFINED_TYPE,
			TypeDef:  structValue,
		}
	}

	sName := node.ImplFor
	
	fmt.Printf("Implentation for: %v\n", env.types[node.ImplFor.Name].DType())

	scope, err := env.ResolveType(sName.Name)
	if err != nil {
		errgen.MakeError(env.filePath, sName.StartPos().Line, sName.EndPos().Line, sName.StartPos().Column, sName.EndPos().Column, err.Error()).Display()
	}

	udType := scope.types[sName.Name].(UserDefined)

	fmt.Printf("User defined type: %v\n", udType.TypeDef.DType())

	// we have impl STRUCT {} or impl TRAIT for STRUCT | primitive {} syntax rules
	// if trait is present, check if all methods are implemented
	if node.Trait != nil {
		for methodName, _ := range traitMethods {
			fmt.Printf("Checking method: %v\n", methodName)
		}
	}

	return nil
}