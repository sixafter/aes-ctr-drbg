#!/bin/bash
# Copyright (c) 2024 Six After, Inc.
#
# This source code is licensed under the Apache 2.0 License found in the
# LICENSE file in the root directory of this source tree.

__dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "${__dir}"/os-type.sh

# Windows
if is_windows; then
    echo "[ERROR] Windows is not currently supported." >&2
    exit 1
fi

# govulncheck
# Ref: https://github.com/golang/vuln
go install golang.org/x/vuln/cmd/govulncheck@latest

go install golang.org/x/perf/cmd/benchstat@latest
