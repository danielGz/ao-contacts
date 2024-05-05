#!/bin/bash

# Default Base URL for the API
DEFAULT_API_URL="http://localhost:8000/contacts"

# Check if an environment variable is set for API_URL
API_URL="${API_URL:-$DEFAULT_API_URL}"

# Define arrays of first and last names
FIRST_NAMES=("Alice" "Bob" "Carol" "David" "Eve" "Frank" "Grace" "Hank" "Ivy" "Jack")
LAST_NAMES=("Johnson" "Smith" "Williams" "Jones" "Brown" "Davis" "Miller" "Wilson" "Moore" "Taylor")

# Get number of first and last names
NUM_FIRST_NAMES=${#FIRST_NAMES[@]}
NUM_LAST_NAMES=${#LAST_NAMES[@]}

# Loop to create 100 users
for i in {1..100}
do
    # Generate a random first name and last name
    RANDOM_FIRST=${FIRST_NAMES[$RANDOM % NUM_FIRST_NAMES]}
    RANDOM_LAST=${LAST_NAMES[$RANDOM % NUM_LAST_NAMES]}
    NAME="${RANDOM_FIRST} ${RANDOM_LAST}"
    EMAIL="${RANDOM_FIRST,,}.${RANDOM_LAST,,}${i}@example.com"  # Convert to lowercase and add index to email

    # Create JSON data for the POST request
    JSON_DATA="{\"name\": \"${NAME}\", \"email\": \"${EMAIL}\"}"

    # Send the POST request to the API
    curl -X POST "$API_URL" \
         -H "Content-Type: application/json" \
         -d "$JSON_DATA" \
         -o /dev/null -s -w "Contact ${i}: %{http_code}\n"
done
