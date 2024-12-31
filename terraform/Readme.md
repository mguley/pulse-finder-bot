
---
#### Step 1: Initialize Infrastructure with Terraform

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
#### Step 2. Run the Setup Script

Copy the `01.sh` setup script to the remote instance. This script will install necessary software and configure the system.
```bash
rsync -rPv --delete remote/ root@<INSTANCE-IP>:/root/remote/
```

Log into the remote instance as root and run the script
```bash
ssh -t root@<INSTANCE-IP> "bash /root/remote/setup/01.sh"
```

This will:
- set up the firewall
- install required software

---
#### Step 3. Initial Login as `bot` User

After the script completes, log in as the `bot` user for additional setup.
```bash
ssh bot@<INSTANCE-IP>
```
On your first login, you'll be prompted to set a password for the `bot` user.

