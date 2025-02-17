echo "🛠 Running code formatter..."
cd compiler
go fmt ./...
cd ..
cd lsp
go fmt ./...
cd ..