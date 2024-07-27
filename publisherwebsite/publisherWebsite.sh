#!/bin/bash

# Variables
PROJECT_DIR="/home/amin/Desktop/yellowbloom/dotanet/publisherwebsite/"
EXECUTABLE_NAME="publisherWebSite"
STATIC_FILES_DIR="/home/amin/Desktop/yellowbloom/dotanet/publisherwebsite/html/"
SERVER_USER="gambron"
SERVER_IP="95.217.125.139"
SERVER_port="2233"
SERVER_DIR="/home/gambron"
LOG_FILE="/home/amin/Desktop/yellowbloom/dotanet/file.log"
SERVER_PASSWORD="Oops123"
PROJECT_URL="0.0.0.0:8084"

# Timestamp function for logging
timestamp() {
    date +"%Y-%m-%d %T"
}

# Log function
log() {
    echo "$(timestamp): $*" | tee -a $LOG_FILE
}

# Step 1: Navigate to the project directory
log "Navigating to project directory: $PROJECT_DIR"
cd $PROJECT_DIR || { log "Failed to navigate to project directory"; exit 1; }

# Step 2: Build the project
log "Building the project"
if go build -o $EXECUTABLE_NAME; then
    log "Build successful"
else
    log "Build failed"
    exit 1
fi

# Step 3: Copy executable and static files to the server
log "Copying executable and static files to the server"
if sshpass -p $SERVER_PASSWORD scp -P $SERVER_port $PROJECT_DIR$EXECUTABLE_NAME $SERVER_USER@$SERVER_IP:$SERVER_DIR && sshpass -p $SERVER_PASSWORD scp -r -P $SERVER_port $STATIC_FILES_DIR $SERVER_USER@$SERVER_IP:$SERVER_DIR; then
    log "Files copied successfully"
else
    log "Failed to copy files"
    exit 1
fi

# Step 4: Run the executable on the server
log "Running the executable on the server"
if sshpass -p $SERVER_PASSWORD ssh -t -p $SERVER_port $SERVER_USER@$SERVER_IP "cd $SERVER_DIR && ./$EXECUTABLE_NAME -publisherservice $PROJECT_URL "; then
    log "Executable ran successfully"
else
    log "Failed to run the executable"
    exit 1
fi

log "Script execution completed"
