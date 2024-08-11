#!/bin/bash

# Set the project root relative to the script's location
PROJECT_ROOT=$(dirname "$(realpath "$0")")

# Create a backup of the existing Nginx configuration
echo "Backing up existing Nginx configuration..."
sudo cp /usr/local/etc/nginx/nginx.conf /usr/local/etc/nginx/nginx.conf.bak

# Create the servers directory if it doesn't exist
NGINX_SERVERS_DIR="/usr/local/etc/nginx/servers"
if [ ! -d "$NGINX_SERVERS_DIR" ]; then
    echo "Creating Nginx servers directory..."
    sudo mkdir -p "$NGINX_SERVERS_DIR"
fi

# Write the Nginx configuration
MVC_CONF="$NGINX_SERVERS_DIR/mvc.conf"
echo "Writing Nginx configuration to $MVC_CONF..."
sudo tee "$MVC_CONF" > /dev/null <<EOL
server {
    listen 80;
    server_name chetak.sdslabs.mvc;

    root "$PROJECT_ROOT";

    location / {
        proxy_pass http://127.0.0.1:8000;
        proxy_set_header Host \$host;
        proxy_set_header X-Real-IP \$remote_addr;
        proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto \$scheme;
    }

    location /static/ {
        alias "$PROJECT_ROOT/static/";
    }

    # Deny access to hidden files
    location ~ /\. {
        deny all;
    }
}
EOL

# Test the Nginx configuration
# echo "Testing Nginx configuration..."
# sudo nginx -t

# Restart Nginx to apply the new configuration
echo "Starting Nginx..."
sudo nginx

echo "Virtual hosting setup completed."
