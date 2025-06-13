#!/bin/bash
# git_ops.sh - Script for managing Git operations for Yinsen

# Function to add all changes and commit with a message
commit() {
  if [ -z "$1" ]; then
    echo "Error: Please provide a commit message"
    exit 1
  fi
  
  git add .
  git commit -m "$1"
  echo "Changes committed with message: $1"
}

# Function to push changes to remote repository
push() {
  git push origin main
  echo "Changes pushed to remote repository"
}

# Function to pull latest changes from remote repository
pull() {
  git pull origin main
  echo "Latest changes pulled from remote repository"
}

# Function to check status of the repository
status() {
  git status
}

# Parse command line arguments
case "$1" in
  commit)
    commit "$2"
    ;;
  push)
    push
    ;;
  pull)
    pull
    ;;
  status)
    status
    ;;
  *)
    echo "Usage: $0 {commit|push|pull|status} [commit message]"
    exit 1
    ;;
esac

exit 0
