#!/bin/bash
oldversion=`git describe --tags \`git rev-list --tags --max-count=1\``
newversion=`git describe --tags \`git rev-list --tags --max-count=1\` | awk -F. -v OFS=. 'NF==1{print ++$NF}; NF>1{if(length($NF+1)>length($NF))$(NF-1)++; $NF=sprintf("%0*d", length($NF), ($NF+1)%(10^length($NF))); print}'`

echo "${oldversion} -> ${newversion}"

git tag "${newversion}" && git push --tags
goreleaser
rm -rf dist/
