#!/bin/bash
set -x
for i in $(ls /usr/share/munin/plugins/mysql_*)
do
    sudo ln -snf $i /etc/munin/plugins
done
sudo service munin-node restart
