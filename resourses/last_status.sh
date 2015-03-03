#!/bin/sh

log_file=$1
status_file=$2

grep Test $log_file | tail -n 2 > aux.txt
grep accuracy aux.txt | awk '{print $NF}' > $status_file.acc.txt
grep loss aux.txt | awk '{print $(NF-5)}' > $status_file.loss.txt

rm aux.txt
