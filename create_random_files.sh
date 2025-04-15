#!/bin/bash

# Directory where files will be created
DIR_NAME=${1:-"random_files"}

# Create directory if it doesn't exist
mkdir -p "$DIR_NAME/small"
mkdir -p "$DIR_NAME/medium"
mkdir -p "$DIR_NAME/large"

echo "Creating small files..."
for i in {1..20}; do
    for j in {1..20}; do
        # Create the subdirectory
        subdir="$DIR_NAME/small/dir-$i/dir-$j"
        echo "Building $subdir"
        mkdir -p "$subdir"


        for k in {1..1000}; do
            dd if=/dev/urandom of="$subdir/file_$k.bin" bs=1000 count=1 status=none
        done
    done
done

echo "Creating medium files..."
for i in {1..50}; do
    dd if=/dev/urandom of="$DIR_NAME/medium/file_$i.bin" bs=1000000 count=1 status=none
done

echo "Creating large files..."
for i in {1..10}; do
    dd if=/dev/urandom of="$DIR_NAME/large/file_$i.bin" bs=1000000000 count=1 status=none
done

echo "Done! All files created successfully."
