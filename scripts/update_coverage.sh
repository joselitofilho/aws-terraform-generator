#!/bin/bash

# Ensure correct usage
if [ $# -ne 2 ]; then
    echo "Usage: $0 <coverage_file> <readme_path>"
    exit 1
fi

# Get the coverage file name and README path from arguments
coverage_file="$1"
readme_path="$2"

# Extract coverage from the last line of the coverage file
new_coverage=$(awk '/^total:/ {print $NF}' "$coverage_file")

# Remove the percentage sign from the coverage value
new_coverage="${new_coverage%\%}"

# Set color based on coverage percentage
if (( $(echo "$new_coverage > 80" | bc -l) )); then
    new_color="2ea44f" # Green color
elif (( $(echo "$new_coverage > 30" | bc -l) )); then
    new_color="yellow" # Yellow color
else
    new_color="red" # Red color
fi

# Search and replace the old code coverage value and color code with the new ones in README
sed -i "s/\(Coverage-[0-9]*\.[0-9]*%25-\)[0-9a-z]*\(\?*\)/Coverage-$new_coverage%25-$new_color\2/g" "$readme_path"

echo "Code coverage value updated to $new_coverage% with color code $new_color in $readme_path"
