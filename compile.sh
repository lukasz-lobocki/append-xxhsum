#!/bin/bash

set -euo pipefail
IFS=$'\n\t'

# List architectures to build
goarchs=('amd64' 'arm64')

# Go to source directory
pushd ~/Code/golang/append-xxhsum/cmd > /dev/null || exit

for i in "${goarchs[@]}"; do
  export GOARCH=${i}
  rm ../bin/append-xxhsum-"${GOARCH}"
  go build \
    -ldflags="-X 'main.version=$(git describe --abbrev=0 --tags)' \
      -X 'main.buildTime=$(date)' \
      -X 'main.gitCommit=$(git rev-parse --short HEAD)' \
      -s -w" \
    -o ../bin/append-xxhsum-"${GOARCH}" \
  .

  # Display the result's characteristics
  file ../bin/append-xxhsum-"${GOARCH}"

done

# For your local architecture, create default file without architecture name
cp ../bin/append-xxhsum-"$(dpkg --print-architecture)" ../bin/append-xxhsum

popd > /dev/null || exit
