#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

go_version_minimum="1.17.9"

go_version() {
    go version | sed -nE -e 's/[^0-9.]+([0-9.]+).+/\1/p'
}

version_lt() {
    # Return true if $1 is a lower version than than $2,
    local ver1=$1
    local ver2=$2
    # Reverse sort the versions, if the 1st item != ver1 then ver1 < ver2
    if  [[ $(echo -e -n "$ver1\n$ver2\n" | sort -rV | head -n1) != "$ver1" ]]; then
        return 0
    else
        return 1
    fi
}

if version_lt "$(go_version)" "$go_version_minimum"; then
    echo "SubnetEVM requires Go >= $go_version_minimum, Go $(go_version) found." >&2
    exit 1
fi

# Avalanche root directory
SUBNET_EVM_PATH=$( cd "$( dirname "${BASH_SOURCE[0]}" )"; cd .. && pwd )

# Load the versions
source "$SUBNET_EVM_PATH"/scripts/versions.sh

# Load the constants
source "$SUBNET_EVM_PATH"/scripts/constants.sh

if [[ $# -eq 1 ]]; then
    binary_path=$1
else
    echo "Invalid arguments to build subnet-evm. Requires one arguments to specify binary location."
    exit 1
fi

# Check if SUBNET_EVM_COMMIT is set, if not retrieve the last commit from the repo.
# This is used in the Dockerfile to allow a commit hash to be passed in without
# including the .git/ directory within the Docker image.
subnet_evm_commit=${SUBNET_EVM_COMMIT:-$( git rev-list -1 HEAD )}
swimmer_id=srSGD5JeYhL8GLx4RUw53VN5TcoBbax6EeCYmy5S3DiteJhdF

# Build Subnet EVM, which is run as a subprocess
echo "Building Subnet EVM Version: $subnet_evm_version $swimmer_version; GitCommit: $subnet_evm_commit"
go build -ldflags "-X github.com/ava-labs/subnet-evm/plugin/evm.GitCommit=$subnet_evm_commit -X github.com/ava-labs/subnet-evm/plugin/evm.Version=$subnet_evm_version-$swimmer_version $static_ld_flags" -o "$binary_path" "plugin/"*.go
