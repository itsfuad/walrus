#!/bin/bash

set -e  # Exit immediately if a command exits with a non-zero status

# Ensure the script is run from the project root
if [ ! -d .git ]; then
    echo "❌ Error: This is not a Git repository."
    exit 1
fi

# Format the codebase
./fmt.sh

echo "✅ Code formatted successfully!"

# Run tests
./test.sh

# Commit changes
read -p "✏️  Enter commit message: " commit_message
git commit -am "$commit_message"

echo "📤 Pushing changes to remote..."
git push

echo "🚀 Done!"