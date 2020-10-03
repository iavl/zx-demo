#!/usr/bin/env bash

set -x

echo 11111111 | \
nchcli vm call --from=$(nchcli keys show -a "$1") \
--contract_addr="$2" \
--method=clear \
--abi_file="$3/contract/paillier.abi" \
--args="$4" \
--gas=3727089 \
--home="$3/nchcli" \
-b block -y
