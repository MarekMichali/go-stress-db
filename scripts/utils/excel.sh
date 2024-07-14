#!/bin/bash

# Define the input and output files
input_file="ready.txt"
output_file="results.csv"

# Check if the input file exists
if [ ! -f "$input_file" ]; then
    echo "Input file not found: $input_file"
    exit 1
fi


# Process each line of the input file
while IFS= read -r line; do
    # Extract the result of the division
    result=$(echo "$line" | awk '{print $NF}')
    # Replace dots with commas in the numbers
    result_with_commas=$(echo "$result" | sed 's/\./,/g')
    # Write the result to the output file
    echo "$result_with_commas" >> "$output_file"
done < "$input_file"

echo "Results have been written to $output_file"