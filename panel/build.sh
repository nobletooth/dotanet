#!/bin/bash

OUTPUT_BINARY="panel"
GOOS="linux"
GOARCH="amd64"
SERVER="gambron@95.217.125.139"
SERVER_PORT="2233"
PANEL_PORT="8081"

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

  echo "Copying binary and static files to server..."
  scp -P $SERVER_PORT $OUTPUT_BINARY $SERVER:~/panel/
  scp -P $SERVER_PORT -r static $SERVER:~/panel/templates

  if [ $? -eq 0 ]; then
    echo "Binary and static files successfully copied to server"

    echo "Starting panel on server..."
    ssh -p $SERVER_PORT $SERVER "chmod +x ~/panel/$OUTPUT_BINARY && nohup ~/panel/$OUTPUT_BINARY -port $PANEL_PORT > ~/panel/server.log 2>&1 &"
    if [ $? -eq 0 ]; then
      echo "Panel started successfully on server."
    else
      echo "Failed to start panel on server."
      exit 1
    fi
  else
    echo "Failed to copy binary or static files to server"
    exit 1
  fi

else
  echo "Build failed."
  exit 1
fi
