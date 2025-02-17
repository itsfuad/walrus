#!/bin/bash

set -e  # Exit immediately if a command exits with a non-zero status

# Ensure the script is run from the project root
if [ ! -d .git ]; then
    echo "❌ Error: This is not a Git repository."
    exit 1
fi

# Format the codebase
echo "🛠 Running code formatter..."
cd compiler
go fmt ./...
cd ..
cd lsp
go fmt ./...
cd ..

echo "✅ Code formatted successfully!"

# Run tests
echo "🧪 Running tests on compiler and LSP modules..."
(cd compiler && go test ./...)
(cd lsp && go test ./...)

echo "✅ All tests passed!"

# Commit changes
read -p "✏️  Enter commit message: " commit_message
git commit -am "$commit_message"

echo "📤 Pushing changes to remote..."
git push

echo "🚀 Done!"