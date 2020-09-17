#!/bin/bash

while getopts pr option
do
case "${option}"
in
p) REPLACE=0; shift;;
r) REPLACE=1; shift;;
esac
done

for FILE in "$@"
do
  if [ $REPLACE -eq 0 ]; then  
    echo "Printing faulty lines in $FILE"
    cat "$FILE" | sed -n '/[^[:print:]]/p'
  elif [ $REPLACE -eq 1 ]; then 
    echo "Deleting faulty lines in $FILE"
    sed -i '/[^[:print:]]/d' "$FILE"
  fi
done


