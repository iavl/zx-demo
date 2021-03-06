#!/usr/bin/env bash

set -x

echo 11111111 | \
nchcli vm call --from=$(nchcli keys show -a "$1" --home="$3/nchcli") \
--contract_addr="$2" \
--method=paillierAdd \
--abi_file="$3/contract/paillier.abi" \
--args="$4 $5" \
--gas=37207089 \
--home="$3/nchcli" \
-y
#-b block -y