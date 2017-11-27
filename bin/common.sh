# common.sh is meant to be sourced by the other bin scripts

# get repo root by looking at this script's parent's parent dir.
DIR="$(readlink -m "${0}/../.." 2>/dev/null || pwd)"

echo "=> RUNNING ${DIR}/${0} [$(date)]"