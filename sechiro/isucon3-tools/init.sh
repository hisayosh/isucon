#!/bin/bash
set -uex
script_dir=$(cd $(dirname $0) && pwd)

# vim
mkdir -p ~/repos/public
if [ ! -d ~/repos/public/dotfiles ]; then
    git clone https://github.com/sechiro/dotfiles.git ~/repos/public/dotfiles
else
    (cd ~/repos/public/dotfiles;git pull)
fi
( cd ~/repos/public/dotfiles
    ./init_vimrc.sh
)

# dstat
sudo yum -y install dstat

if ! grep -e '### init ###' ~/.bashrc >/dev/null ;then
    cat <<EOL >> ~/.bashrc
### init ###
export HISTSIZE=1000000
alias dstat-extra='dstat -Tclmdrn --top-cpu --top-cputime --top-io --top-latency'
alias dstat-full='dstat -Tclmdrn'
alias dstat-mem='dstat -Tclm'
alias dstat-cpu='dstat -Tclr'
alias dstat-net='dstat -Tclnd'
alias dstat-disk='dstat -Tcldr'
EOL
fi

# epel tools
if ! yum repolist | grep ^epel > /dev/null ; then
    wget http://ftp-srv2.kddilabs.jp/Linux/distributions/fedora/epel/6/x86_64/epel-release-6-8.noarch.rpm -P /tmp
    sudo yum localinstall /tmp/epel-release-6-8.noarch.rpm
fi

sudo yum -y install bash-completion etckeeper mosh python-pip