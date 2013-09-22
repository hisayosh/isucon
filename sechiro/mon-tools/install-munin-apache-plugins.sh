#!/bin/bash
set -x

# Change httpd setting
echo "Add Directive:"
echo "ExtendedStatus On

<Location /server-status>
    SetHandler server-status
    Order deny,allow
    Deny from all
    Allow from localhost
</Location>" | sudo tee /etc/httpd/conf.d/server-status.conf
sudo service httpd graceful

# Install munin plugins
for i in $(ls /usr/share/munin/plugins/apache_*)
do
    sudo ln -snf $i /etc/munin/plugins
done
sudo service munin-node restart
