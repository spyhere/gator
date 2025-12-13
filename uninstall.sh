#!/bin/bash

set -e
read -n 1 -r -p "Are you sure you want to uninstall gator? y/n " ANSWER
echo

if [[ $ANSWER == "y" || $ANSWER == "Y" ]]; then
  echo "Uninstalling..."
  BIN_LOCATION="$HOME/go/bin/gator"

  if [[ -f $BIN_LOCATION ]]; then
    rm "$BIN_LOCATION"
    echo "gator has been removed"
  else
    echo "gator hasn't been found at the expected ~/go/bin/ directory"
  fi
else
  echo "Abort"
fi

