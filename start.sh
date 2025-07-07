#!/bin/bash
cd ~/repositories/golang-midtrans-service
nohup ./golang-midtrans-service > ~/golang-midtrans.log 2>&1 &
echo "App started. Logs: ~/golang-midtrans.log"