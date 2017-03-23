#!/bin/sh

git checkout --orphan tmp
git add -A
git commit -am "first commit"
git branch -D master
git branch -m master
git push -f origin master
