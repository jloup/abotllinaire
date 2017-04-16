#!/usr/bin/env bash
v=$(node --version 2>&1); echo node $v
v=$(npm --version 2>&1); echo npm $v
v=$(babel --version 2>&1); echo babel $v
v=$(browserify --version 2>&1); echo browserify $v

npm install
rc=$?; if [[ $rc != 0 ]]; then exit $rc; fi

exec bower install
