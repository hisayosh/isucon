#!/bin/bash
echo "Install munin server"
sudo yum -y --enablerepo=epel install munin
echo ""
echo "Install munin node"
sudo yum -y --enablerepo=epel install munin-node

echo "###################################"
echo "Do Next":
echo "
# Set password
sudo htpasswd -c /etc/munin/munin-htpasswd Admin

# Set interval (default: 5min)
sudo vi /etc/cron.d/munin
"
