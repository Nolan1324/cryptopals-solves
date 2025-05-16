set -e

if [ -z "$1" ]; then
  echo "Error: no argument provided"
  exit 1
fi

if [ "$1" = "all" ]; then
  i=1
  while true; do
    dir=$(find sets/*/. -name "$i")
    if [ -z "$dir" ]; then
      break
    fi
    echo ""
    echo "===== Running challenge $i in $dir ====="
    (cd "$dir" && go run .)
    i=$((i + 1))
  done
  exit 0
fi

dir=$(find sets/*/. -name "$1")

if [ -z "$dir" ]; then
  echo "Error: no matching directory found"
  exit 1
fi

cd "$dir"
go run .
