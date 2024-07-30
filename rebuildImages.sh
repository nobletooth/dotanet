SERVER="gambron@95.217.125.139"
SERVER_PORT="2233"
SERVER_PASSWORD="Oops123"

sshpass -p $SERVER_PASSWORD ssh -t -p $SERVER_PORT $SERVER 'cd /home/gambron/Desktop/dotanet && docker-compose down && docker-compose up --build --force-recreate'

