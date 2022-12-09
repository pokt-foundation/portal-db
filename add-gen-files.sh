#!/bin/bash

make

if [[ `git diff-index --quiet HEAD` ]]; then
    echo "files were modified by this hook. please add changed files and re-commit"
    exit 1
fi
