#!/bin/bash

# Define input and output files
if [ "$#" -ne 2 ]; then
    echo "Usage: $0 <input_file> <output_file>.csv"
    exit 1
fi

input_file="$1"
output_file="$2"

# Check every 4th line and extract the required information using egrep and awk
awk 'NR % 4 == 0' $input_file | egrep -o 'totalExecutionTimeMillis: [0-9]+' | awk '{print $2}' > temp_millis.txt
awk 'NR % 4 == 0' $input_file | egrep -o 'count: [0-9]+' | awk '{print $2}' > temp_count.txt

# Combine the two temporary files
paste temp_millis.txt temp_count.txt > "extracted_data.txt"

# Remove temporary files
rm temp_millis.txt temp_count.txt


# Use awk to divide the first column by the second and save the result to a CSV file
echo "Input file: $input_file" >> "$output_file"
awk '{printf "%.5f\n", $1/$2}' extracted_data.txt | sed 's/\./,/' >> $output_file
rm extracted_data.txt