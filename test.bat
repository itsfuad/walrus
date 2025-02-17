@echo off
echo starting tests
echo Compiler tests
cd compiler
go test ./...
cd ..
echo Compiler tests complete
echo LSP tests
cd lsp
go test ./...
echo tests complete