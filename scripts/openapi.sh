#!/bin/bash

set -e

$(realpath .)/scripts/openapi-go.sh
$(realpath .)/scripts/openapi-ts.sh
