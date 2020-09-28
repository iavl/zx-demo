#!/usr/bin/env bash

set -x

nchcli q vm call $(nchcli keys show -a alice) "$1" \
queryResult ./contract/paillier.abi \
--args="$2"