1. Make

2. Release
look .env
set -a
source .env
set +a
goreleaser release --clean
