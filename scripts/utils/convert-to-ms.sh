#!/bin/bash

# Check if a file name is provided
if [ "$#" -ne 2 ]; then
    echo "Usage: $0 <input_file> <output_file>"
    exit 1
fi

input_file="$1"
output_file="$2"

# Check if the input file exists
if [ ! -f "$input_file" ]; then
    echo "Input file not found: $input_file"
    exit 1
fi

# Process the file and save the output to the output file
while IFS= read -r line; do
    if [[ "$line" =~ h$ ]]; then
        # Extract the number and convert hours to milliseconds
        time_in_hours=$(echo "$line" | grep -o '[0-9.]*')
        time_in_ms=$(echo "$time_in_hours * 3600 * 1000" | bc)
        echo "${line}: $time_in_ms" >> "$output_file"
    elif [[ "$line" =~ min$ ]]; then
        # Extract the number and convert minutes to milliseconds
        time_in_min=$(echo "$line" | grep -o '[0-9.]*')
        time_in_ms=$(echo "$time_in_min * 60 * 1000" | bc)
        echo "${line}: $time_in_ms" >> "$output_file"
    elif [[ "$line" =~ s$ ]]; then
        # Extract the number and convert seconds to milliseconds
        time_in_sec=$(echo "$line" | grep -o '[0-9.]*')
        time_in_ms=$(echo "$time_in_sec * 1000" | bc)
        echo "${line}: $time_in_ms" >> "$output_file"
    else
        # Just print the line if it doesn't match the patterns
        echo "$line" >> "$output_file"
    fi
done < "$input_file"