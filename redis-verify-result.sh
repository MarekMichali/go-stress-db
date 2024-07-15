#!/bin/bash


docker update --cpus 4 redis
op_type=$1

# First, get the slow log from redis
SLOWLOG=$(docker exec -i redis redis-cli slowlog get 700)
#echo -e "\n\n\n\n$SLOWLOG" > redisfull.txt
echo -e "$SLOWLOG" > redisfull.txt
# Call the test.py script
python test.py "$op_type"

rm redisfull.txt

# Extract the values before the "get" word and save to another file

# Define the input file

# Define the output file
#output_file="output.txt"

# Initialize the count variable
#count=0

# Use awk to process every 7th line starting from the first line for select
#awk 'NR % 7 == 0 { sum += $1; count++; print $0 >> "'"$output_file"'" } END { printf "Sum: %.5f\nCount: %d\n", sum, count }' redisfull.txt

# Use awk to process every 8th line starting from the first line for
#awk 'NR % 8 == 0 { sum += $1; count++; print $0 >> "'"$output_file"'" } END { printf "Sum: %.5f\nCount: %d\n", sum, count }' redisfull.txt

#rm redisfull.txt
#rm output.txt
