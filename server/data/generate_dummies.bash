#!/bin/bash

rm file*
rm -r folder*

for ((i=1; i<=10; i++))
do
    touch "file$i.txt"
    mkdir "folder$i"
done
