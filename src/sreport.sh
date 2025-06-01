#!/bin/bash

MODE=${1:-daily}

get_used() {
  local START=$1
  local END=$2

  sreport cluster UserUtilizationByAccount start=${START} end=${END} 2>/dev/null | \
  awk '
    BEGIN { data = 0 }
    /^[[:space:]]*Cluster[[:space:]]+Login/ { data = 1; next }
    /^[[:space:]]*-+/ { next }
    data == 1 && NF >= 6 {
      print $(NF-1)
      exit
    }
  '
}

if [[ $MODE == "daily" ]]; then
  for i in {5..0}; do
    DATE=$(date -d "$i day ago" +%Y-%m-%d)
    START="${DATE}T00:00:00"
    END="${DATE}T23:59:59"
    USED=$(get_used "$START" "$END")
    echo "$DATE ${USED:-0}"
  done

elif [[ $MODE == "weekly" ]]; then
  for i in {5..0}; do
    # Calculate the Saturday of i weeks ago
    START_DATE=$(date -d "this saturday - $i week" +%Y-%m-%d)
    START="${START_DATE}T00:00:00"

    # Calculate the Friday of that week (6 days after Saturday)
    END_DATE=$(date -d "$START_DATE +6 days" +%Y-%m-%d)
    END="${END_DATE}T23:59:59"

    USED=$(get_used "$START" "$END")
    echo "$START_DATE ${USED:-0}"
  done

elif [[ $MODE == "monthly" ]]; then
  for i in {5..0}; do
    START=$(date -d "$(date +%Y-%m-01) -$i month" +%Y-%m-01T00:00:00)
    END=$(date -d "$START +1 month -1 day" +%Y-%m-%dT23:59:59)
    MONTH_LABEL=$(date -d "$START" +%Y-%m)
    USED=$(get_used "$START" "$END")
    echo "$MONTH_LABEL ${USED:-0}"
  done

else
  echo "Usage: $0 [daily|weekly|monthly]"
  exit 1
fi
