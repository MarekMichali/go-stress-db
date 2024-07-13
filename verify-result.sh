#!/bin/bash

docker update --cpus 2 mongodb
op_type=$1

if [[ "$op_type" == "insert" ]]; then
    result=$(docker exec -i mongodb bash <<'EOF'
    mongosh <<'EOSQL'
    db.system.profile.aggregate([ { $match: { op: "insert", ns: "test.videos" } }, { $group: { _id: null, totalExecutionTimeMillis: { $sum: "$millis" }, count: { $sum: 1 } } }] )
    EOSQL
EOF
    )
elif [[ "$op_type" == "select" ]]; then
    result=$(docker exec -i mongodb bash <<'EOF'
    mongosh <<'EOSQL'
    db.system.profile.aggregate([ { $match: { op: "select", ns: "test.videos" } }, { $group: { _id: null, totalExecutionTimeMillis: { $sum: "$millis" }, count: { $sum: 1 } } }] )
    EOSQL
EOF
    )
elif [[ "$op_type" == "update" ]]; then
    result=$(docker exec -i mongodb bash <<'EOF'
    mongosh <<'EOSQL'
    db.system.profile.aggregate([ { $match: { op: "update", ns: "test.videos" } }, { $group: { _id: null, totalExecutionTimeMillis: { $sum: "$millis" }, count: { $sum: 1 } } }] )
    EOSQL
EOF
    )
elif [[ "$op_type" == "delete" ]]; then
    result=$(docker exec -i mongodb bash <<'EOF'
    mongosh <<'EOSQL'
    db.system.profile.aggregate([ { $match: { op: "remove", ns: "test.videos" } }, { $group: { _id: null, totalExecutionTimeMillis: { $sum: "$millis" }, count: { $sum: 1 } } }] )
    EOSQL
EOF
    )
else
    echo "Invalid operation type"
fi

# Filtering to only lines which start with '{' (indicating JSON output), and removing leading/trailing whitespace.
filtered_result=$(echo $result | grep -o "{.*}" | sed 's/^ *//;s/ *$//')

echo $filtered_result
