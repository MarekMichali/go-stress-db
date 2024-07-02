# go-stress-db

To initialize the database, follow these steps:

1. Navigate to the `/scripts/init-databases` directory or `/scripts/init-custom-databases`. If you choose custom databases, build the images. The database configuration is built into the image to avoid using a Docker volume. The `/init-databases` directory uses a Docker volume to pass the configuration to the container.

2. Run the initialization script in `/init-databases` for the desired database.

Next, go to the `configure-databases` directory and perform the following:

1. Run the configuration script for the desired database.

Once the above steps are completed, you can execute the `./run-test.sh` script. Detailed instructions on how to use this script are provided in the script manual.