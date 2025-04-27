
if [ -z "$1" ]; then
  echo "Error: no argument provided"
  exit 1
fi

dir=$(find sets/*/. -name "$1")

if [ -z "$dir" ]; then
  echo "Error: no matching directory found"
  exit 1
fi

cd "$dir"
go run .
