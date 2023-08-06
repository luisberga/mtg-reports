#!/bin/bash

cat <<EOL > config.yaml
api:
  db:
    user: "root"
    password: "root"
    host: "localhost or db (if it's running in a docker container)"
    port: "3306"
    database: "MTGREPORTS"
    commitSize: 1000
  port: ":8080"
  log:
    level: "debug"

conciliatejob:
  db:
    user: "root"
    password: "root"
    host: "localhost or db (if it's running in a docker container)"
    port: "3306"
    database: "MTGREPORTS"
    commitSize: 1000
  timeout: "1h"
  log:
    level: "debug"
  exchange:
    url: "https://v6.exchangerate-api.com/v6/your_key/latest/USD"

reportjob:
  db:
    user: "root"
    password: "root"
    host: "localhost or db (if it's running in a docker container)"
    port: "3306"
    database: "MTGREPORTS"
  timeout: "1h"
  log:
    level: "debug"
  email:
    host: "smtp.your_host.com"
    username: "your_user@email.com"
    password: "your_password"
    to: "your_destination@email.com"
    port: "587"
EOL

echo "config.yaml generated successfully."
