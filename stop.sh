#!/bin/bash
PID=$(ps aux | grep golang-midtrans-service | grep -v grep | awk '{print $2}')
if [ -n "$PID" ]; then
  kill $PID
  echo "App (PID $PID) stopped."
else
  echo "App not running."
fi
