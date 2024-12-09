#!/bin/bash
set -euo pipefail  # Enhanced error handling: -e (exit on error), -u (treat unset variables as errors), -o pipefail (catch errors in pipelines)

# ======================================================================== #
# VARIABLES
# ======================================================================== #

TIMEZONE="Europe/Berlin"             # Set the timezone
USERNAME="bot"                       # Name of the new user to create

# MongoDB credentials (should be in sync with your application requirements)
MONGO_USER="bot_user"
MONGO_PASSWORD="bot_password"

# Export locale to avoid any locale-related errors.
export LC_ALL=en_US.UTF-8

# ======================================================================== #
# FUNCTIONS
# ======================================================================== #

# Update and upgrade system packages
update_system() {
    echo "Updating system packages..."
    apt-get update -q
    apt-get --yes -o Dpkg::Options::="--force-confnew" upgrade
}

# Enable necessary repositories
enable_repositories() {
    echo "Enabling universe repository..."
    add-apt-repository --yes universe
}

# Set timezone and install locales
setup_time_and_locale() {
    echo "Setting timezone to ${TIMEZONE}..."
    timedatectl set-timezone "${TIMEZONE}"
    echo "Installing all locales..."
    apt-get --yes install locales-all
}

# Create a new user with sudo privileges and SSH access
create_user() {
    if id "${USERNAME}" &>/dev/null; then
        echo "User ${USERNAME} already exists. Skipping creation."
    else
        echo "Creating user ${USERNAME}..."
        useradd --create-home --shell "/bin/bash" --groups sudo "${USERNAME}"
        passwd --delete "${USERNAME}"
        chage --lastday 0 "${USERNAME}"

        echo "Copying SSH keys to new user..."
        if [ -d "/root/.ssh" ]; then
            rsync --archive --chown="${USERNAME}:${USERNAME}" /root/.ssh /home/"${USERNAME}"
        else
            echo "No SSH keys found in /root/.ssh. Skipping SSH key copy."
        fi
    fi
}

# Configure firewall
configure_firewall() {
    echo "Configuring firewall to allow services..."
    ufw allow 22           # SSH
    ufw allow 27017        # MongoDB
    ufw allow 9050/tcp     # Tor SocksPort
    ufw allow 9051/tcp     # Tor ControlPort
    ufw --force enable
}

# Install and configure MongoDB
setup_mongodb() {
    echo "Updating package list and installing prerequisites..."
    apt-get update -q
    apt-get install -y gnupg curl

    echo "Importing MongoDB GPG key..."
    curl -fsSL https://www.mongodb.org/static/pgp/server-7.0.asc | gpg --dearmor -o /usr/share/keyrings/mongodb-server-7.0.gpg

    echo "Creating MongoDB source list for Ubuntu 22.04 (Jammy)..."
    echo "deb [ arch=amd64,arm64 signed-by=/usr/share/keyrings/mongodb-server-7.0.gpg ] https://repo.mongodb.org/apt/ubuntu jammy/mongodb-org/7.0 multiverse" | tee /etc/apt/sources.list.d/mongodb-org-7.0.list

    echo "Updating package list..."
    apt-get update -q

    echo "Installing MongoDB..."
    apt-get install -y mongodb-org

    echo "Configuring MongoDB..."
    # Backup the original configuration file
    cp /etc/mongod.conf /etc/mongod.conf.bak
    # Modify the configuration file
    echo -e "\nsecurity:\n  authorization: enabled" >> /etc/mongod.conf

    echo "Starting and enabling MongoDB service..."
    systemctl restart mongod
    systemctl enable mongod

    echo "Waiting for MongoDB service to initialize..."
    sleep 10  # Adjust this delay if needed

    echo "Checking MongoDB service status..."
    if ! systemctl is-active --quiet mongod; then
        echo "MongoDB service failed to start. Check /var/log/mongodb/mongod.log for details. Exiting..."
        exit 1
    fi

    echo "Creating MongoDB user..."
    retry=0
    until mongosh --eval "db.runCommand({ connectionStatus: 1 })" &>/dev/null || [ $retry -ge 5 ]; do
        echo "Waiting for MongoDB to accept connections..."
        sleep 5
        retry=$((retry + 1))
    done

    if [ $retry -ge 5 ]; then
        echo "MongoDB did not become available. Exiting..."
        exit 1
    fi

    mongosh <<EOF
use admin
db.createUser({
  user: "${MONGO_USER}",
  pwd: "${MONGO_PASSWORD}",
  roles: [{ role: "root", db: "admin" }]
})
EOF

    echo "MongoDB setup complete!"
}

# Verify MongoDB connection and user authentication
verify_mongodb() {
    echo "Verifying MongoDB connection and user authentication..."

    mongosh --host 127.0.0.1 --port 27017 -u "$MONGO_USER" -p "$MONGO_PASSWORD" --authenticationDatabase admin <<EOF
use admin
db.runCommand({ connectionStatus: 1 })
EOF

    if [ $? -eq 0 ]; then
        echo "MongoDB connection verified successfully!"
    else
        echo "MongoDB connection verification failed!"
        exit 1
    fi
}

# Install and configure Tor proxy
setup_tor() {
    echo "Installing Tor and necessary dependencies..."
    apt-get update -q
    apt-get install -y tor curl netcat-openbsd

    echo "Configuring Tor..."
    cat <<EOF >/etc/tor/torrc
SocksPort 0.0.0.0:9050
ControlPort 0.0.0.0:9051
HashedControlPassword 16:EC1800A189DA53D6600B08E22D26B20C2A34E24962AA23FC6E5AA8B8F4
EOF

    echo "Starting and enabling Tor service..."
    systemctl restart tor
    systemctl enable tor

    echo "Waiting for Tor to initialize..."
    sleep 15
}

# Validate Tor proxy installation
verify_tor() {
    echo "Verifying Tor proxy installation..."

    echo "Checking exit IP using Tor..."
    exit_ip=$(curl --socks5-hostname localhost:9050 https://httpbin.org/ip --max-time 10 2>/dev/null | jq -r '.origin')
    if [ -z "$exit_ip" ]; then
        echo "Tor proxy validation failed! Unable to fetch exit IP."
        exit 1
    fi

    echo "Exit IP via Tor: $exit_ip"

    echo "Checking Tor control port..."
    echo -e 'authenticate "password"\ngetinfo status/circuit-established' | nc -w 5 localhost 9051 > /tmp/tor_control_status

    if grep -q "250-status/circuit-established=1" /tmp/tor_control_status; then
        echo "Tor proxy is fully operational."
    else
        echo "Tor proxy validation failed! Circuit not established."
        echo "Tor control port response:"
        cat /tmp/tor_control_status
        exit 1
    fi
}

# ======================================================================== #
# MAIN SCRIPT
# ======================================================================== #

main() {
    enable_repositories
    update_system
    setup_time_and_locale
    create_user
    configure_firewall
    setup_mongodb
    verify_mongodb
    setup_tor
    verify_tor

    echo "Script complete! Rebooting..."
    reboot
}

main "$@"