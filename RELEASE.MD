1. Make

2. Release
set -a
source .env
set +a
goreleaser release --clean
