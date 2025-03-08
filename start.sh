#!/bin/bash

SERVER_NAME="BlessBot"

Start_Bot() {

#  nohup ./$SERVER_NAME >/dev/null 2> server.log &
  nohup ./"${SERVER_NAME}" > server.log  2>&1 &

  echo "${SERVER_NAME} Server is started."
  >nohup.out
  tail -10f server.log
}

Stop_Bot() {
  # kill $1 if it exists.
  PID_LIST=$(ps -ef | grep $SERVER_NAME | grep -v grep | awk '{printf "%s ", $2}')
  for PID in $PID_LIST; do
    if kill -9 $PID; then
      echo "Process $one($PID) was stopped at " $(date)
      echo "${SERVER_NAME} Server is stoped."
    fi
  done
}

Status_Bot() {
  PID_NUM=$(ps -ef | grep $SERVER_NAME | grep -v grep | wc -l)
  if [ $PID_NUM -gt 0 ]; then
    {
      echo "${SERVER_NAME} server is started."
    }
  else
    {
      echo "${SERVER_NAME} server is stoped."
    }
  fi
}

case "$1" in
'start')
  Start_Bot
  ;;
'stop')
  Stop_Bot
  ;;
'restart')
  Stop_Bot
  Start_Bot
  ;;
'status')
  Status_Bot
  ;;
*)
  echo "Usage: $0 {start|stop}"
  echo "  start : To start the application of ${SERVER_NAME}"
  echo "  stop  : To stop the application of ${SERVER_NAME}"
  echo "  restart  : To restart the application of ${SERVER_NAME}"
  echo "  status  : To view status the application of ${SERVER_NAME}"
  RETVAL=1
  ;;
esac

exit 0
