#!/bin/bash

docker update --cpus 4 mongodb
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
    db.system.profile.aggregate([ { $match: { op: "query", ns: "test.videos" } }, { $group: { _id: null, totalExecutionTimeMillis: { $sum: "$millis" }, count: { $sum: 1 } } }] )
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

filtered_result=$(echo $result | grep -o "{.*}" | sed 's/^ *//;s/ *$//')

echo $filtered_result
