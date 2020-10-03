#!/usr/bin/env bash

set -x

nchcli q vm call $(nchcli keys show -a alice) "$1" \
queryResult $2/contract/paillier.abi \
--args="$3" --home="$2/nchcli"