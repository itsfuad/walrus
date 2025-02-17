echo "🧪 Running tests on compiler and LSP modules..."
(cd compiler && go test ./...)
(cd lsp && go test ./...)

echo "✅ All tests passed!"