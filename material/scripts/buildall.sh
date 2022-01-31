#!/usr/bin/env bash

# Specify platforms to build for
platforms=("linux/amd64" "linux/arm64" "darwin/amd64" "darwin/arm64", "windows/amd64")

# Clean build directory
builddir="build"
if [ -d $builddir ]; then rm -rf build; fi
mkdir -p $builddir

# Build for all specified platforms
for platform in "${platforms[@]}"
do
    platform_split=(${platform//\// })
    GOOS=${platform_split[0]}
    GOARCH=${platform_split[1]}
    output_name="${builddir}/gotz_${GOOS}_${GOARCH}"
    if [ $GOOS = "windows" ]; then
        output_name+=".exe"
    fi

    echo "Building $output_name..."
    env GOOS=$GOOS GOARCH=$GOARCH go build -o $output_name .
    if [ $? -ne 0 ]; then
        echo "Error occurred during build, continuing..."
    fi
done
