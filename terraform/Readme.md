
---
#### Step 1: Generate an RSA Key Pair

You need to create an RSA key pair (public and private keys) on your local machine (inside the terraform container).

```bash
ssh-keygen -t rsa -b 2048 -f ~/.ssh/id_rsa
```

This command will:
- Generate a private key file: `~/.ssh/id_rsa`
- Generate a public key file: `~/.ssh/id_rsa.pub`

Ensure private key has appropriate file permissions

```bash
chmod 600 ~/.ssh/id_rsa
```

---
#### Step 2: Initialize Infrastructure with Terraform

Initialize configuration, preview change, and apply the changes.

```bash
terraform init
```
```bash
terraform plan
```
```bash
terraform apply
```

---
#### Step 3. Run the Setup Script

Copy the `01.sh` setup script to the remote instance. This script will install necessary software and configure the system.
```bash
rsync -rPv --delete remote/ root@<INSTANCE-IP>:/root/remote/
```

Log into the remote instance as root and run the script
```bash
ssh -t root@<INSTANCE-IP> "bash /root/remote/01.sh"
```

This will:
- set up the firewall
- install required software

---
#### Step 4. Initial Login as `bot` User

After the script completes, log in as the `bot` user for additional setup.
```bash
ssh bot@<INSTANCE-IP>
```
On your first login, you'll be prompted to set a password for the `bot` user.

