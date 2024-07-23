OUTPUT_BINARY="publisherwebsite" # Name of the output binary
GOOS="linux"
GOARCH="amd64" # Target architecture

if [ ! -f "go.mod" ]; then
  echo "go.mod not found"
  exit 1
fi

echo "Building the project for $GOOS/$GOARCH..."
GOOS=$GOOS GOARCH=$GOARCH go build -o $OUTPUT_BINARY .

if [ $? -eq 0 ]; then
  echo "Build successful! Binary created: $OUTPUT_BINARY"
else
  echo "Build failed."
  exit 1
fi