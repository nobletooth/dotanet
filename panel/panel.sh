OUTPUT_BINARY="Panel"
SERVER="gambron@95.217.125.139"
SERVER_PORT="2233"
PROJECT_URL="0.0.0.0:8085"
PROJECT_DIR="./"
TEMPLATES_DIR="./templates"
PUBLISHER_SCRIPT_DIR="./publisher"
LOG_FILE="./file.log"
SERVER_PASSWORD="Oops123"
SERVER_DIR="/home/gambron"
ADSERVERURL="http://95.217.125.139:8081"

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
  log "Copying additional files..."
  sshpass -p $SERVER_PASSWORD scp -P $SERVER_PORT -r $TEMPLATES_DIR $SERVER:$SERVER_DIR
  if [ $? -eq 0 ]; then
    log "templates files successfully copied to server"
  else
    log "Failed to copy templates to server"
      exit 1
  fi

  sshpass -p $SERVER_PASSWORD scp -P $SERVER_PORT -r $PUBLISHER_SCRIPT_DIR $SERVER:$SERVER_DIR
  if [ $? -eq 0 ]; then
  log "publisher scripts files successfully copied to server"
  else
    log "Failed to copy publisher scripts files to server"
    exit 1
  fi



else
  log "Failed to copy binary to server"
  exit 1
fi

log "Starting panel..."


sshpass -p $SERVER_PASSWORD ssh -t -p $SERVER_PORT $SERVER "cd $SERVER_DIR && ./$OUTPUT_BINARY -dbuser user -dbpassword password -dbname dotanet -dbport 5432 -dbhost localhost -panelurl $PROJECT_URL -adserverurl $ADSERVERURL"
if [ $? -eq 0 ]; then
  log "Panel started on port $PROJECT_URL"
  log "Deployment completed."
else
  log "Failed to start panel."
  exit 1
fi
