clear
rm peer
go build ./...
./peer 127.0.0.1:8080 &
./peer 127.0.0.1:8081 &
./peer 127.0.0.1:8082 &

# Wait for user input
echo "Press Enter to terminate all processes"
read -r

# Terminate all background processes
pkill -f "peer 127.0.0.1"