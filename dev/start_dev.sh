#!/bin/bash
set -e

export SECRET=4ef5b29539b70233dd40c02a1799d25079595565e05a193b09da2c3e60ada1cd

export SIMPLE_AUTH_ENABLE=false
export OIDC_ENABLE=true
export OIDC_ISSUER=http://localhost:8082
export OIDC_CLIENT_SECRET=secret

export ARANGO_DB_HOST=http://localhost:8529
export ARANGO_DB_PASSWORD=foobar
export S3_HOST=http://localhost:9000
export S3_PASSWORD=minio123

export AUTH_BLOCK_NEW=false
export AUTH_DEFAULT_ROLES=analyst,admin

export EXTERNAL_ADDRESS=http://localhost
export CATALYST_ADDRESS=http://host.docker.internal
export INITIAL_API_KEY=d0169af94c40981eb4452a42fae536b6caa9be3a

go run ../cmd/catalyst-dev/*.go
