To start PMM, run the following command in your terminal:

```bash
docker run --detach --restart always \
--publish 443:443 \
-v pmm-data:/srv \
--name pmm-server \
percona/pmm-server:2.41.2
```

After starting PMM, you need to create a Docker network and add both the PMM server and the database you want to monitor to the network.

By default, PMM provides dashboards for MySQL/MariaDB and MongoDB. To add a database service, navigate to the configuration section and select "Add Service". For MongoDB, make sure to enable the "Use QAN MongoDB Profiler" option.

If you want to monitor Redis, you will need to install the Redis plugin. Go to the configuration section and select "Plugins" to install the Redis plugin. Once installed, you can search for a dashboard named "Redis" to monitor your Redis container. It should be in the "General" section.