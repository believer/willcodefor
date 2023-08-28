#!/bin/bash

# Calculate the SHA256 hash of a file and return the first 6 characters
function calculate_hash() {
    local filename="$1"
    local hash=$(shasum -a 256 "$filename" | cut -c1-6)
    echo "$hash"
}

echo "Running the public hasher"

# Remove the old CSS file
rm ./public/styles.*.css

# Generate the new CSS file
npx tailwindcss -i ./styles.css -o ./public/styles.css

# Add the 6 first characters of the hash to the file name
# of the generated CSS file. This is to bust the cache of the CSS file.

# Get the hash of the generated CSS file
hash=$(calculate_hash ./public/styles.css)

# Rename the generated CSS file
mv ./public/styles.css ./public/styles.$hash.css

# Replace the old CSS file with the new one
# The first argument is an empty string to skip backup in macOS
# Uses a counted range, {0,1}, to simulate a ? operator for the hash
# This matches both styles.css and styles.hash.css
sed -i "" "s/styles\(\.[a-z0-9]\{6\}\)\{0,1\}\.css/styles\.$hash\.css/g" ./views/layouts/main.html

echo "Tailwind generated and updated"

# Fixing HTMX scripts

hash=$(calculate_hash ./public/htmx.*.min.js)
mv ./public/htmx.*.min.js ./public/htmx.$hash.min.js
sed -i "" "s/htmx\.[a-z0-9]\{6\}\.min/htmx\.$hash\.min/g" ./views/layouts/main.html

hash=$(calculate_hash ./public/htmx-head-support.*.js)
mv ./public/htmx-head-support.*.js ./public/htmx-head-support.$hash.js
sed -i "" "s/htmx-head-support\.[a-z0-9]\{6\}/htmx-head-support\.$hash/g" ./views/layouts/main.html

echo "HTMX scripts fixed"

# Fixing Command menu script

hash=$(calculate_hash ./public/command-menu.*.js)
mv ./public/command-menu.*.js ./public/command-menu.$hash.js
sed -i "" "s/command-menu\.[a-z0-9]\{6\}/command-menu\.$hash/g" ./views/layouts/main.html

echo "Command menu script fixed"
