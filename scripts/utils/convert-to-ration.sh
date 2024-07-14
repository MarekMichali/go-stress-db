#!/bin/bash

# Path to the input file
input_file="$1"
# Output file name
output_file="$2"

# Read the file line by line
while IFS= read -r line; do
    # Check if the line contains "total_time"
    if [[ $line =~ total_time ]]; then
        # Extract the last number on the line (the divisor)
        divisor=$(echo $line | awk '{print $NF}')
        
        # Read the next line to get the dividend
        IFS= read -r next_line
        dividend=$(echo $next_line)
        
        # Perform the division and print the result
        if [[ $divisor != 0 ]]; then # Check to avoid division by zero
            result=$(echo "scale=2; $dividend / $divisor" | bc)
            echo "$dividend / $divisor = $result" >> "$output_file"
        fi
    fi
done < "$input_file"