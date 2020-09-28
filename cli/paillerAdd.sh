#!/usr/bin/env bash

set -x

echo 11111111 | \
nchcli vm call --from=$(nchcli keys show -a "$1") \
--contract_addr="$2" \
--method=paillierAdd \
--abi_file="./contract/paillier.abi" \
--args="$3 $4" \
--gas=37207089  -b block -y