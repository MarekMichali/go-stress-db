#!/bin/bash

# Check if the number of arguments is correct
if [ "$#" -ne 2 ]; then
    echo "Usage: $0 <inputFile> <outputFile>"
    exit 1
fi

inputFile="$1"
outputFile="$2"

# Initialize a counter
counter=0

# Initialize an array to hold every 5th line's value
fifthLines=()

echo "Input file: $inputFile" >> "$outputFile"

# Read the file line by line
while IFS= read -r line; do
    # Increment counter
    ((counter++))
    
    # Check if the line number is a multiple of 5
    if ((counter % 5 == 0)); then
        # Store the value of the current (5th) line
        fifthLines+=("$line")
    fi
    
    # Check if the line number is one more than a multiple of 5 (i.e., line after every 5th line)
    if ((counter % 5 == 1)) && [ ${#fifthLines[@]} -ne 0 ]; then
        # Get the last value from fifthLines array
        lastValue=${fifthLines[${#fifthLines[@]}-1]}
        
        # Perform division with 5 digits precision
        result=$(echo "scale=5; $lastValue / $line" | bc)
        
        # Print the result
        echo "$result" >> "$outputFile"
        
        # Remove the last element from the array
        unset fifthLines[${#fifthLines[@]}-1]
    fi
done < "$inputFile"