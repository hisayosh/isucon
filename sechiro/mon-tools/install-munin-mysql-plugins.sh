#!/bin/bash
set +x
for i in $(cat munin-mysql-plugin-names)
do
    sudo ln -snf $i /etc/munin/plugins
done
sudo service munin-node restart
