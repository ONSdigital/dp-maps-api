#!/bin/bash -eux

pushd dp-maps-api
  make build
  cp build/dp-maps-api Dockerfile.concourse ../build
popd
