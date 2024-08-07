OUTPUT_BINARY="Adserver"
SERVER="gambron@95.217.125.139"
SERVER_PORT="2233"
PROJECT_URL="0.0.0.0:8081"
PANEL_URL="http://localhost:8085"
PROJECT_DIR="./"
LOG_FILE="./file.log"
SERVER_PASSWORD="Oops123"
SERVER_DIR="/home/gambron"
EVENTSERVICE_URL="http://95.217.125.139:8082"

export GOOS=linux
export GOARCH=amd64
export GIN_MODE=release

GO_VERSION="1.22.5"
go version | grep -q "go$GO_VERSION" || {
    echo "Go version problem!"
    exit 1
}

timestamp() {
    date +"%Y-%m-%d %T"
}

log() {
    echo "$(timestamp): $*" | tee -a $LOG_FILE
}

log "Building Linux binary"
go build -o $OUTPUT_BINARY

if [ $? -ne 0 ]; then
    log "Build failed."
    exit 1
fi

log "Build successful! Binary created: $OUTPUT_BINARY"

# Make the binary executable
chmod +x $OUTPUT_BINARY
log "Made panel executable"

log "Copying binary to server..."
sshpass -p $SERVER_PASSWORD scp -P $SERVER_PORT $PROJECT_DIR$OUTPUT_BINARY $SERVER:$SERVER_DIR

if [ $? -eq 0 ]; then
  log "Binary successfully copied to server"
else
  log "Failed to copy binary to server"
  exit 1
fi

log "Starting adserver..."


sshpass -p $SERVER_PASSWORD ssh -t -p $SERVER_PORT $SERVER "cd $SERVER_DIR && ./$OUTPUT_BINARY -adserverurl $PROJECT_URL -panelurl $PANEL_URL -eventserviceurl $EVENTSERVICE_URL -newAdTreshold 5 -newAdProb 0.25 -expAdProb 0.75"
if [ $? -eq 0 ]; then
  log "adserver started on $PROJECT_URL"
  log "Deployment completed."
else
  log "Failed to start."
  exit 1
fi

