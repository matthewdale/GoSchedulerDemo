#!/bin/bash
jmeter -n -t Starvation.jmx -Jhost=$1 -Jport=$2 -Jthreads=$3
