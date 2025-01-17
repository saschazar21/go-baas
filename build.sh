#!/bin/bash
set +eux

# Setting the environment variables
FUNCTIONS_DIR="functions"

# Export Go build environment variables
export GOOS="linux"
export GOARCH="amd64"
export GO111MODULE=on
export CGO_ENABLED=0

echo "Building the project..."

# Create the build directory
if [ ! -d $FUNCTIONS_DIR ]; then
  mkdir $FUNCTIONS_DIR
fi

# Build the functions
for v in $(ls -d $PWD/cmd/*); do
  for e in $(ls -d $v/*); do
    echo -ne "Building $e..."
    cd $e

    # Exctract the folder names
    endpoint=$(basename $e)   # endpoint name
    version=$(basename $v)    # API version number

    go build -o ../../../$FUNCTIONS_DIR/${version}_${endpoint} -ldflags="-s -w" main.go
    
    echo " done"
    cd ..
  done
done