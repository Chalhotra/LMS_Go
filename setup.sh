#!/bin/bash

isInstalled() {
    if command -v "$1" &> /dev/null; then
        echo "$1 is installed."
    else 
        echo "$1 is not installed. Aborting..."
        exit 1
    fi
}

echo $'Checking for dependencies...\n'

isInstalled "go"
isInstalled "mysql"

echo $'Dependencies are installed.\n'

echo "Kindly enter your MySQL username: "
read username

echo "Kindly enter your MySQL password: "
read -s password

echo "Kindly enter your DB name: "
read dbName

mysql -u "$username" -p"$password" -e "CREATE DATABASE IF NOT EXISTS $dbName;"

echo "Created database 'mvc' successfully!"

migrate -path ./database/migration/ -database "mysql://$username:$password@tcp(localhost:3306)/$dbName" -verbose up

echo "Database migrated successfully!"

echo "Database configuration file created successfully!"

go mod vendor
go mod tidy

echo "Kindly run the following command to start the server: go run main.go"