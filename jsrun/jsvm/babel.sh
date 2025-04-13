#!/bin/bash

set -euo pipefail

LATEST_VERSION=$(curl -i https://unpkg.com/@babel/standalone | grep -i Location | grep -Eo '[0-9][.0-9]+[0-9]')
echo "Updating to latest version: ${LATEST_VERSION}"
curl https://unpkg.com/@babel/standalone@${LATEST_VERSION}/babel.min.js | sed 's/# sourceMappingURL=//' > babel.js
