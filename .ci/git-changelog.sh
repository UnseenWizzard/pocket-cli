#!/bin/bash
if [ "$#" -ne 2 ]; then
    echo "Missing required arguments"
    echo "run as: ./git_log_release_notes.sh VERSION GIT_REPO_URL"
fi

echo "# Release $1" > CHANGELOG.md
echo "" >> CHANGELOG.md
git log --pretty=format:"%s ([%h]($2/commits/%h))" | sort >> CHANGELOG.md