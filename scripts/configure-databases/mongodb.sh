#!/bin/bash

docker start mongodb

docker exec -i mongodb bash <<'EOF'

mongosh <<'EOSQL'

db.getSiblingDB("admin").createRole({
    role: "explainRole",
    privileges: [{
        resource: {
            db: "",
            collection: ""
        },
        actions: [
            "listIndexes",
            "listCollections",
            "dbStats",
            "dbHash",
            "collStats",
            "find"
        ]
    }],
    roles: []
})

db.getSiblingDB("admin").createRole({
    "role": "pbmAnyAction",
    "privileges": [{
        "resource": {
            "anyResource": true
        },
        "actions": [
            "anyAction"
        ]
    }],
    "roles": []
})

db.getSiblingDB("admin").createUser({
    user: "pmm",
    pwd: "pmm",
    roles: [
        { role: "explainRole", db: "admin" },
        { role: "clusterMonitor", db: "admin" },
        { role: "read", db: "local" },
        { "db" : "admin", "role" : "readWrite", "collection": "" },
        { "db" : "admin", "role" : "backup" },
        { "db" : "admin", "role" : "clusterMonitor" },
        { "db" : "admin", "role" : "restore" },
        { "db" : "admin", "role" : "pbmAnyAction" }
    ]
})

db.getSiblingDB("admin").updateUser("pmm", {
    roles: [
        { role: "explainRole", db: "admin" },
        { role: "clusterMonitor", db: "admin" },
        { role: "read", db: "local" },
        { "db" : "admin", "role" : "readWrite", "collection": "" },
        { "db" : "admin", "role" : "backup" },
        { "db" : "admin", "role" : "clusterMonitor" },
        { "db" : "admin", "role" : "restore" },
        { "db" : "admin", "role" : "pbmAnyAction" }
    ]
})

db.videos.drop()
db.setProfilingLevel(0)
db.system.profile.drop()
db.createCollection( "system.profile", { capped: true, size:1024000000 } )
db.setProfilingLevel(2)
EOSQL

EOF