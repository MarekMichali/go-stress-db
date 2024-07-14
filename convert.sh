#!/bin/bash

# Check if a file name is provided
if [ "$#" -ne 2 ]; then
    echo "Usage: $0 <input_file> <output_file>.csv"
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
        echo "${line}: $time_in_ms" >> "tmp.txt"
    elif [[ "$line" =~ min$ ]]; then
        # Extract the number and convert minutes to millisecondss
        time_in_min=$(echo "$line" | grep -o '[0-9.]*')
        time_in_ms=$(echo "$time_in_min * 60 * 1000" | bc)
        echo "${line}: $time_in_ms" >> "tmp.txt"
    elif [[ "$line" =~ ms$ ]]; then
        # Extract the number and convert milliseconds to milliseconds
        time_in_ms=$(echo "$line" | grep -o '[0-9.]*')
        echo "${line}: $time_in_ms" >> "tmp.txt"
    elif [[ "$line" =~ s$ ]]; then
        # Extract the number and convert seconds to milliseconds
        time_in_sec=$(echo "$line" | grep -o '[0-9.]*')
        time_in_ms=$(echo "$time_in_sec * 1000" | bc)
        echo "${line}: $time_in_ms" >> "tmp.txt"
    else
        # Just print the line if it doesn't match the patterns
        echo "$line" >> "tmp.txt"
    fi
done < "$input_file"

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
            result=$(echo "scale=5; $dividend / $divisor" | bc)
            echo "$dividend / $divisor = $result" >> "tmp2.txt"
        fi
    fi
done < "tmp.txt"

# Add the name of the input file to the output file
echo "Input file: $input_file" >> "$output_file"
# Process each line of the input file
while IFS= read -r line; do
    # Extract the result of the division
    result=$(echo "$line" | awk '{print $NF}')
    # Replace dots with commas in the numbers
    result_with_commas=$(echo "$result" | sed 's/\./,/g')
    # Write the result to the output file
    echo "$result_with_commas" >> "$output_file"
done < "tmp2.txt"

rm tmp.txt
rm tmp2.txt

echo "Results have been written to $output_file"