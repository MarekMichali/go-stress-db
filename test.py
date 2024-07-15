import sys

# Define the path to the input file
input_file_path = 'redisfull.txt'

# Check for command-line arguments for operation type
operation = 'set'  # Default operation
if len(sys.argv) > 1:
    if sys.argv[1] in ['set', 'get', 'del']:
        operation = sys.argv[1]
    else:
        print("Invalid argument. Defaulting to 'set'.")
else:
    print("No operation specified. Defaulting to 'set'.")

# Initialize a variable to hold the previous line
previous_line = ''
sum_of_previous_values = 0.0
occurences =0
# Open the input file to read
with open(input_file_path, 'r') as input_file:
    for line in input_file:
        # Strip the newline character from the end of the line
        line = line.strip()
        # Check if the current line contains the operation ('set' or 'get')
        if line == operation:
            # Print the previous line (value before the operation)
            if previous_line.isdigit():
                occurences += 1
                sum_of_previous_values += float(previous_line)
                #print(previous_line)
        # Update the previous line
        previous_line = line
print(sum_of_previous_values)
print(occurences)
#print(f"Values before '{operation}' have been printed.")
