#!/usr/bin/env bash

nchcli q vm call $(nchcli keys show -a alice) "$1" \
queryResult ../contract/paillier.abi \
--args="$2" \
--home ./keys