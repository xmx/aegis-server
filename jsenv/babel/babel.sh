#!/bin/bash

LATEST_VERSION=$(curl -i https://unpkg.com/@babel/standalone | grep -i Location | grep -Eo '[.0-9]+')
curl https://unpkg.com/@babel/standalone@$LATEST_VERSION/babel.min.js | sed 's/# sourceMappingURL=//' > babel.min.js