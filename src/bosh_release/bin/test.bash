#!/bin/bash

set -eu
set -o pipefail

bosh_target

go run github.com/onsi/ginkgo/v2/ginkgo ${@}
