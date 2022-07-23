#!/bin/bash
if [ "$#" -ne 2 ]; then
    echo "Missing required arguments"
    echo "run as: ./git_log_release_notes.sh VERSION GIT_REPO_URL"
    exit 1
fi

NEW_VERSION=$1 

# Get the last release tag (reverse sort by date, limited to one, return only refname, if pattern matches a hardening tag)
LATEST_RELEASE=$(git for-each-ref --sort=-creatordate  --count=1 --format="%(refname)" --no-contains=HEAD refs/tags/release/*)

#  Get extra short gitlog from laste release to current HEAD
GIT_CHANGELOG=$(git log $LATEST_RELEASE..HEAD --format="- %s ([%h]($2/commits/%h))" | sort)

FEATURES=""
FIXES=""
OTHER=""

prefix=""
while IFS= read -r line; do
    current_prefix=$(echo "$line" | sed -E 's/-\s*(.*?):\s.*/\1/g')
    if [ "$prefix" != "$current_prefix" ]; then
        if [ "feat" == "$current_prefix" ]; then
          FEATURES+="\n## 🚀 Features\n"
        elif [ "fix" == "$current_prefix" ]; then
          FIXES+="\n## 🐛 Fixes\n"
        else
          OTHER+="\n### $current_prefix\n"
        fi
        
        prefix=$current_prefix
    fi
    if [ "feat" == "$prefix" ]; then
        FEATURES+="$line\n"
    elif [ "fix" == "$prefix" ]; then
        FIXES+="$line\n"
    else
        OTHER+="$line\n"
    fi
done <<< "$GIT_CHANGELOG"

printf "# Release $NEW_VERSION\n$FEATURES$FIXES\n## 🛠️ Other Changes\n$OTHER" > CHANGELOG.md