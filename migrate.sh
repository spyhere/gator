#!/bin/bash

if [[ "$#" -eq 0 ]]; then
  echo "no arguments passed. Expected 'up'/'down'"
  exit 1
fi

DB_URL="postgres://spyhere@localhost:5432/gator"
if [[ "$1" == "up" ]]; then
  cd ./sql/schema/ && goose postgres $DB_URL up
elif [[ "$1" == "down" ]]; then
  cd ./sql/schema/ && goose postgres $DB_URL down
fi

