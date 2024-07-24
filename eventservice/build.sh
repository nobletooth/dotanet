OUTPUT_BINARY="eventservice"
GOOS="linux"
GOARCH="amd64"
SERVER="gambron@95.217.125.139"
SERVER_PORT="2233"
EVENTSERVICE_PORT="8082"



if [ ! -f "go.mod" ]; then
  echo "go.mod not found"
  exit 1
fi

echo "Building the project for $GOOS/$GOARCH..."
GOOS=$GOOS GOARCH=$GOARCH go build -o $OUTPUT_BINARY .

if [ $? -eq 0 ]; then
  echo "Build successful! Binary created: $OUTPUT_BINARY"

  # Make the binary executable
  chmod +x $OUTPUT_BINARY
  echo "Made $OUTPUT_BINARY executable"

  ssh $SERVER -p $SERVER_PORT
  scp $OUTPUT_BINARY $SERVER:~/eventservice/
  if [ $? -eq 0 ]; then
    echo "Binary successfully copied to server"

    # Run the eventservice
    echo "Starting eventservice..."
    ssh $SERVER -p $SERVER_PORT "./eventservice -port $EVENTSERVICE_PORT"
    ./eventservice 
  else
    echo "Failed to copy binary to server"
    exit 1
  fi

else
  echo "Build failed."
  exit 1
fi