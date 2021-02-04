#!/usr/bin/env bash

set -e

pushd ..
make install
popd

rm -fr .terraform*
terraform init
