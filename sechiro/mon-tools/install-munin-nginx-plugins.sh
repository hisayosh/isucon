#!/bin/bash
set -x

# Install munin plugins
for i in $(ls /usr/share/munin/plugins/nginx_*)
do
    sudo ln -snf $i /etc/munin/plugins
done
sudo service munin-node restart
