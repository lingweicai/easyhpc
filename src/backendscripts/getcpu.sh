#!/bin/bash

echo -e "NodeName\tCPUs\tAlloc\tLoad\tAlloc%\tLoad%"

scontrol show nodes | awk '
/NodeName=/ {
    split($1, a, "="); node=a[2]
}
/CPUTot=/ {
    for (i=1; i<=NF; i++) {
        split($i, a, "=")
        if (a[1]=="CPUTot") cputot=a[2]
        if (a[1]=="CPUAlloc") cpualloc=a[2]
        if (a[1]=="CPULoad") cpuload=a[2]
    }
    alloc_percent = (cpualloc / cputot) * 100
    load_percent = (cpuload / cputot) * 100
    printf "%s\t%d\t%d\t%.2f\t%.1f%%\t%.1f%%\n", node, cputot, cpualloc, cpuload, alloc_percent, load_percent
}'