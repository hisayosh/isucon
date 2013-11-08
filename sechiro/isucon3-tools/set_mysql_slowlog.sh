#!/bin/bash
set -ue
# http://d.hatena.ne.jp/sh2/20090414
# mysqldumpslow -s t /tmp/test.log

if [ "$1" = ALL ];then
    cat <<EOL
    set global slow_query_log = 1;
    set global slow_query_log_file = '/tmp/test.log';
    set global long_query_time = 0;
EOL

else
    cat <<EOL
set slow_query_log_file = '/var/lib/mysql/slow.log'
set global long_query_time = $1;
EOL
fi

