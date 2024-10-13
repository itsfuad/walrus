package typechecker

import "walrus/ast"

func checkTraitDeclaration(traitStmt ast.TraitDeclStmt, env *TypeEnvironment) ValueTypeInterface {


	methods := make(map[string]Fn)

	for name, method := range traitStmt.Methods {
		methods[name] = checkFunctionSignature(name, method.FunctionType, env)
	}

	trait := Trait{
		DataType: TRAIT_TYPE,
		TraitName: traitStmt.Trait.Name,
		Methods: methods,
	}

	env.DeclareTrait(traitStmt.Trait.Name, trait)

	return nil
}