#!/bin/bash
#===============================================================================
# SCRIPT: setup-environment.sh
#
# DESCRIPTION:
#   This script automates the setup of a server for an application that requires:
#     1. Enabling necessary repositories.
#     2. Updating system packages.
#     3. Setting the system timezone and installing locales.
#     4. Creating a new user with sudo privileges.
#     5. Configuring a firewall for secure operation.
#     6. Installing and configuring MongoDB with a dedicated user, database,
#        and collections for the application.
#     7. Installing and verifying a Tor proxy setup.
#     8. Setting environment variables for system-wide use.
#     9. Rebooting the system to apply changes.
#
# USAGE:
#   1. Copy this script to a fresh server.
#   2. Run as root: `sudo ./setup-environment.sh`
#   3. Wait for the script to complete; the system will reboot automatically.
#
# NOTES:
#   - This script assumes a Debian/Ubuntu-based system with apt and ufw installed.
#   - Make sure to verify that all variables (e.g., MongoDB credentials, proxy settings)
#     align with your application's requirements before running.
#
# LICENSE: MIT
#===============================================================================
set -euo pipefail

# ======================================================================== #
# VARIABLES
# ======================================================================== #

TIMEZONE="Europe/Berlin"             # Timezone to configure on the server
USERNAME="bot"                       # New user to create with sudo privileges

# MongoDB credentials for the application
MONGO_USER=bot_user
MONGO_PASS=bot_password
MONGO_HOST=127.0.0.1
MONGO_PORT=27017
MONGO_DB=bot_db
MONGO_URLS_COLLECTION=urls
MONGO_VACANCY_COLLECTION=vacancies

SOURCE_ALFA_SITEMAP_URL=
SOURCE_BETA_SITEMAP_URL=
SOURCE_GAMMA_SITEMAP_URL=
SOURCE_BATCH_SIZE=1

# Proxy settings
PROXY_HOST=127.0.0.1
PROXY_PORT=9050
PROXY_CONTROL_PASSWORD=password
PROXY_CONTROL_PORT=9051
PROXY_PING_URL=https://api.ipify.org?format=json

# gRPC client/application
VACANCY_SERVER_ADDRESS=api.pulse-finder.mguley.com:64055
AUTH_SERVER_ADDRESS=api.pulse-finder.mguley.com:63055
AUTH_ISSUER=grpc.pulse-finder.bot
ENV=prod

# Export locale to ensure consistent system behavior
export LC_ALL=en_US.UTF-8

# ======================================================================== #
# FUNCTIONS
# ======================================================================== #

# ------------------------------------------------------------------------------
# update_system
#
# Updates and upgrades all system packages to ensure software is up-to-date.
# ------------------------------------------------------------------------------
update_system() {
    echo "Updating system packages..."
    apt-get update -q
    apt-get --yes -o Dpkg::Options::="--force-confnew" upgrade
}

# ------------------------------------------------------------------------------
# enable_repositories
#
# Enables necessary repositories (e.g., universe repository) to support
# installation of required packages.
# ------------------------------------------------------------------------------
enable_repositories() {
    echo "Enabling universe repository..."
    add-apt-repository --yes universe
}

# ------------------------------------------------------------------------------
# setup_time_and_locale
#
# Configures the server timezone and installs locales to ensure compatibility
# with internationalization needs.
# ------------------------------------------------------------------------------
setup_time_and_locale() {
    echo "Setting timezone to ${TIMEZONE}..."
    timedatectl set-timezone "${TIMEZONE}"
    echo "Installing all locales..."
    apt-get --yes install locales-all
}

# ------------------------------------------------------------------------------
# create_user
#
# Creates a new system user with sudo privileges. If root's SSH keys are
# available, they will be copied to the new user's home directory.
# ------------------------------------------------------------------------------
create_user() {
    if id "${USERNAME}" &>/dev/null; then
        echo "User ${USERNAME} already exists. Skipping creation."
    else
        echo "Creating user ${USERNAME}..."
        useradd --create-home --shell "/bin/bash" --groups sudo "${USERNAME}"
        passwd --delete "${USERNAME}"               # Remove password (force password reset on first login)
        chage --lastday 0 "${USERNAME}"             # Expire password immediately

        echo "Copying SSH keys to new user..."
        if [ -d "/root/.ssh" ]; then
            rsync --archive --chown="${USERNAME}:${USERNAME}" /root/.ssh /home/"${USERNAME}"
        else
            echo "No SSH keys found in /root/.ssh. Skipping SSH key copy."
        fi
    fi
}

# ------------------------------------------------------------------------------
# configure_firewall
#
# Configures the firewall (ufw) to allow specific services and enable the firewall.
# ------------------------------------------------------------------------------
configure_firewall() {
    echo "Configuring firewall to allow services..."
    ufw allow 22           # SSH
    ufw allow 9050/tcp     # Tor SocksPort
    ufw allow 9051/tcp     # Tor ControlPort
    ufw --force enable
}

# ------------------------------------------------------------------------------
# set_environment_variables
#
# Sets essential environment variables in /etc/environment for global access.
# ------------------------------------------------------------------------------
set_environment_variables() {
      echo "Adding environment variables to /etc/environment..."
      {
        # MongoDB
        echo "MONGO_USER=${MONGO_USER}"
        echo "MONGO_PASS=${MONGO_PASS}"
        echo "MONGO_HOST=${MONGO_HOST}"
        echo "MONGO_PORT=${MONGO_PORT}"
        echo "MONGO_DB=${MONGO_DB}"
        echo "MONGO_URLS_COLLECTION=${MONGO_URLS_COLLECTION}"
        echo "MONGO_VACANCY_COLLECTION=${MONGO_VACANCY_COLLECTION}"
        # Sources
        echo "SOURCE_ALFA_SITEMAP_URL=${SOURCE_ALFA_SITEMAP_URL}"
        echo "SOURCE_BETA_SITEMAP_URL=${SOURCE_BETA_SITEMAP_URL}"
        echo "SOURCE_GAMMA_SITEMAP_URL=${SOURCE_GAMMA_SITEMAP_URL}"
        echo "SOURCE_BATCH_SIZE=${SOURCE_BATCH_SIZE}"
        # Proxy
        echo "PROXY_HOST=${PROXY_HOST}"
        echo "PROXY_PORT=${PROXY_PORT}"
        echo "PROXY_CONTROL_PASSWORD=${PROXY_CONTROL_PASSWORD}"
        echo "PROXY_CONTROL_PORT=${PROXY_CONTROL_PORT}"
        echo "PROXY_PING_URL=${PROXY_PING_URL}"
        # gRPC
        echo "VACANCY_SERVER_ADDRESS=${VACANCY_SERVER_ADDRESS}"
        echo "AUTH_SERVER_ADDRESS=${AUTH_SERVER_ADDRESS}"
        echo "AUTH_ISSUER=${AUTH_ISSUER}"
        echo "ENV=${ENV}"
      } >> /etc/environment
}

# ------------------------------------------------------------------------------
# setup_mongodb
#
# Installs and configures MongoDB, enabling authorization and creating the
# necessary MongoDB admin user for the application.
# ------------------------------------------------------------------------------
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
    cp /etc/mongod.conf /etc/mongod.conf.bak
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
  pwd: "${MONGO_PASS}",
  roles: [{ role: "root", db: "admin" }]
})
EOF

    echo "MongoDB setup complete!"
}

# ------------------------------------------------------------------------------
# initialize_mongodb
#
# Creates the application database and initializes collections.
# ------------------------------------------------------------------------------
initialize_mongodb() {
    echo "Creating database '${MONGO_DB}' and initializing collections..."
    mongosh -u "$MONGO_USER" -p "$MONGO_PASS" --authenticationDatabase admin <<EOF
use ${MONGO_DB}
db.createCollection("${MONGO_URLS_COLLECTION}")
db.createCollection("${MONGO_VACANCY_COLLECTION}")
EOF
    echo "Database and collections initialized successfully!"
}

# ------------------------------------------------------------------------------
# verify_mongodb
#
# Verify MongoDB connection and user authentication.
# ------------------------------------------------------------------------------
verify_mongodb() {
    echo "Verifying MongoDB connection and user authentication..."

    mongosh --host 127.0.0.1 --port 27017 -u "$MONGO_USER" -p "$MONGO_PASS" --authenticationDatabase admin <<EOF
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

# ------------------------------------------------------------------------------
# setup_tor
#
# Installs and configures the Tor proxy for the application.
# ------------------------------------------------------------------------------
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

# ------------------------------------------------------------------------------
# verify_tor
#
# Validates Tor proxy installation.
# ------------------------------------------------------------------------------
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
    set_environment_variables
    setup_mongodb
    verify_mongodb
    initialize_mongodb
    setup_tor
    verify_tor

    echo "Script complete! Rebooting..."
    reboot
}

main "$@"