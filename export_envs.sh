#!/bin/bash

# Input file and environment name
file=$1     
env_name=$2

echo $file
echo $env_name

# Known environment variable keys
keys=("aws_access_key_id" "aws_secret_access_key" "aws_session_token" "aws_security_token")

# Initialize an empty string for export commands
export_string=""

# Process the file line by line
in_section=false
while IFS= read -r line || [[ -n "$line" ]]; do
    # Check if the line starts the desired section
    if [[ "$line" == "[$env_name]" ]]; then
        in_section=true
        continue
    fi

    # If in the desired section, process valid key-value pairs
    if $in_section; then
        # Stop if we encounter a blank line or another section
        if [[ -z "$line" || "$line" =~ ^\[ ]]; then
            break
        fi

        # Check if the line starts with one of the known keys
        for key in "${keys[@]}"; do
            if [[ "$line" == "$key ="* ]]; then
                # Extract the value after '='
                value="${line#*=}"
                value=$(echo "$value" | xargs)  # Trim whitespace
                key_upper=$(echo "$key" | tr '[:lower:]' '[:upper:]')  # Convert key to uppercase
                export_string+="export $key_upper=\"$value\""$'\n'
                break
            fi
        done
    fi
done < "$file"

# Execute the export commands
eval "$export_string"

# Print the exports for confirmation
echo "Exported variables:"
printf "%s" "$export_string"
