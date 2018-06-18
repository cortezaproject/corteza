#!/bin/bash
echo $(echo $(sed -n 's/^\(type \|\s*\)\([A-Z][^[:space:]]*\) interface .*/\2/p' "interfaces.go") | tr ' ' ',')
