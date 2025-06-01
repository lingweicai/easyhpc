#!/bin/bash

# Output final JSON
echo "["

# Get all node hostnames
nodes=$(sinfo -h -N -o "%N")

first=1
for node in $nodes; do
    # Gather node details from scontrol
    node_info=$(scontrol show node "$node")

    # CPU usage
    CPUs=$(echo "$node_info" | awk '/CPUTot/ { for(i=1;i<=NF;i++) if($i ~ /CPUTot=/) { split($i,a,"="); print a[2] } }')
    CPUsAlloc=$(echo "$node_info" | awk '/CPUAlloc/ { for(i=1;i<=NF;i++) if($i ~ /CPUAlloc=/) { split($i,a,"="); print a[2] } }')

    cpuUsage=0
    if [[ "$CPUs" -gt 0 ]]; then
        cpuUsage=$((CPUsAlloc * 100 / CPUs))
    fi

    # Memory usage
    RealMemory=$(echo "$node_info" | awk '/RealMemory=/ { for(i=1;i<=NF;i++) if($i ~ /RealMemory=/) { split($i,a,"="); print a[2] } }')
    AllocMem=$(echo "$node_info" | awk '/AllocMem=/ { for(i=1;i<=NF;i++) if($i ~ /AllocMem=/) { split($i,a,"="); print a[2] } }')

    memUsage=0
    if [[ "$RealMemory" -gt 0 ]]; then
        memUsage=$((AllocMem * 100 / RealMemory))
    fi

    # Storage usage (requires SSH to node)
    # You can tweak this if needed (e.g., /home or /scratch)
    diskUsage=$(ssh -o BatchMode=yes -o ConnectTimeout=2 "$node" "df --output=pcent / | tail -1 | tr -dc '0-9'") || diskUsage=0

    # Output JSON entry
    if [[ $first -eq 0 ]]; then
        echo ","
    fi
    echo "  {"
    echo "    \"node\": \"$node\","
    echo "    \"cpuUsage\": $cpuUsage,"
    echo "    \"memoryUsage\": $memUsage,"
    echo "    \"storageUsage\": $diskUsage"
    echo "  }"
    first=0
done

echo "]"
