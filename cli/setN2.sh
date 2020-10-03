#!/usr/bin/env bash

set -x

echo 11111111 | \
nchcli vm call --from=$(nchcli keys show -a "$1") \
--contract_addr="$2" \
--method=setN2 \
--abi_file="$3/contract/paillier.abi" \
--args="$4" \
--gas=98669000 \
--home="$3/nchcli" \
-b block -y