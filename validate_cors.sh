#!/bin/bash
cd /home/node/.openclaw/workspace/daily-kotoba
go build -o /dev/null ./internal/middleware/... 2>&1 || echo "BUILD FAILED"
echo "Done"
