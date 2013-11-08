#!/bin/bash
set -ux
USERS="sechiro hisayosh qtakamitsu"
PASSWORD="password###"
SALT=1
for i in $USERS
do
    sudo useradd $i
    sudo usermod -p $(/usr/bin/perl -e 'print crypt(${ARGV[0]}, ${ARGV[1]})' $PASSWORD} ${SALT}) $i
    sudo mkdir /home/$i/.ssh
    sudo cp /home/$USER/.ssh/authorized_keys /home/$i/.ssh
    sudo chown $i:$i /home/$i/.ssh/authorized_keys
done
