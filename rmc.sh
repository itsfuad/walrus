#!/bin/bash

# Ensure the script is run inside a Git repository
if [ ! -d .git ]; then
    echo "Error: This is not a Git repository."
    exit 1
fi

# Get commit hash from user
read -p "Enter the commit hash to remove: " commit_hash

# Confirm before proceeding
read -p "⚠️ This will rewrite history! Are you sure? (y/N): " confirm
if [[ "$confirm" != "y" && "$confirm" != "Y" ]]; then
    echo "Aborted."
    exit 0
fi

# Find how many commits back the given commit is
commit_count=$(git rev-list --count $commit_hash..HEAD)
if [[ -z "$commit_count" || "$commit_count" -eq 0 ]]; then
    echo "❌ Commit not found in history!"
    exit 1
fi

echo "🔍 Removing commit $commit_hash using interactive rebase..."

# Start interactive rebase to remove the commit
GIT_SEQUENCE_EDITOR="sed -i '/$commit_hash/d'" git rebase -i HEAD~$commit_count

# Force push to remote (WARNING: This rewrites history)
echo "🚀 Force pushing updated history..."
git push --force

echo "✅ Commit $commit_hash removed successfully!"
