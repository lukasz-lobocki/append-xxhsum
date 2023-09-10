#!/bin/bash

set -euo pipefail
IFS=$'\n\t'

# Last segment of current path
modulename=${PWD##*/}
modulename=${modulename:-/}  # to correct for the case where PWD=/

if [[ -f ./cmd/"${modulename}.go" ]]; then
  :
else
  echo "No file ${modulename}.go found."
  exit
fi

# List architectures to build
goarchs=('amd64' 'arm64')

for i in "${goarchs[@]}"; do
  export GOARCH=${i}
  rm --force ./bin/"${modulename}"-"${GOARCH}"
  go build \
    -ldflags="-X 'main.version=$(git describe --abbrev=0 --tags)+${GOARCH}.$(git rev-parse --short HEAD)' \
      -s -w" \
    -o ./bin/"${modulename}"-"${GOARCH}" \
  ./cmd

  # Display the result's characteristics
  file ./bin/"${modulename}"-"${GOARCH}"

done

# For your local architecture, create default file without architecture name suffix
cp ./bin/"${modulename}"-"$(dpkg --print-architecture)" ./bin/"${modulename}"

# Local copy to directory of executables
sudo cp ./bin/"${modulename}" /usr/local/bin/"${modulename}"

# Remote copy arm64 build to nextcloudpi
scp ./bin/"${modulename}"-arm64 la_lukasz@nextcloudpi.local:./tmp/"${modulename}"
