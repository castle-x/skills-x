#!/usr/bin/env bash
set -euo pipefail

usage() {
  cat <<'EOF'
Usage: scripts/release.sh <version>

Example:
  scripts/release.sh 0.2.9
EOF
}

die() {
  echo "ERROR: $*" >&2
  exit 1
}

require_cmd() {
  command -v "$1" >/dev/null 2>&1 || die "Missing required command: $1"
}

if [[ $# -ne 1 ]]; then
  usage
  exit 1
fi

VERSION="$1"
if [[ ! "$VERSION" =~ ^[0-9]+\.[0-9]+\.[0-9]+([.-][0-9A-Za-z.-]+)?$ ]]; then
  die "Invalid version: $VERSION (expected semver like 0.2.9)"
fi

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$REPO_ROOT"

require_cmd git
require_cmd gh
require_cmd npm
require_cmd make

BRANCH="$(git rev-parse --abbrev-ref HEAD)"
if [[ "$BRANCH" != "main" ]]; then
  die "Must be on main branch (current: $BRANCH)"
fi

if git rev-parse "v$VERSION" >/dev/null 2>&1; then
  die "Tag v$VERSION already exists"
fi

if ! gh auth status -h github.com >/dev/null 2>&1; then
  die "gh is not authenticated (run: gh auth login)"
fi

if ! npm whoami >/dev/null 2>&1; then
  die "npm is not authenticated (run: npm login)"
fi

echo "Updating npm version to $VERSION..."
( cd "$REPO_ROOT/npm" && npm version "$VERSION" --no-git-tag-version )

echo "Building local binary..."
make build

TEST_DIR=""
cleanup() {
  if [[ -n "${TEST_DIR}" && -d "${TEST_DIR}" ]]; then
    rm -rf "${TEST_DIR}"
  fi
}
trap cleanup EXIT

echo "Testing all skills install..."
TEST_DIR="$(mktemp -d)"
echo "Testing in: ${TEST_DIR}"
./bin/skills-x init --all --target "${TEST_DIR}"

echo "Building npm release binaries..."
make build-npm

ASSETS=(
  "npm/bin/skills-x-linux-amd64"
  "npm/bin/skills-x-linux-arm64"
  "npm/bin/skills-x-darwin-amd64"
  "npm/bin/skills-x-darwin-arm64"
  "npm/bin/skills-x-windows-amd64.exe"
)

for asset in "${ASSETS[@]}"; do
  [[ -f "$asset" ]] || die "Missing release asset: $asset"
done

if [[ -z "$(git status --porcelain)" ]]; then
  die "No changes to commit"
fi

echo "Committing changes..."
git add -A
git commit -m "chore: release v$VERSION"

echo "Tagging release..."
git tag -a "v$VERSION" -m "v$VERSION"

echo "Pushing main and tag..."
git push origin main
git push origin "v$VERSION"

echo "Creating GitHub release with assets..."
gh release create "v$VERSION" \
  --title "v$VERSION" \
  --notes "Release v$VERSION" \
  "${ASSETS[@]}"

echo "Publishing to npm..."
( cd "$REPO_ROOT/npm" && npm publish --access public )

echo "Release complete: v$VERSION"
