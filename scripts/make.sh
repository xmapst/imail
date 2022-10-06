#!/bin/bash

_os=`uname`
_path=`pwd`
_dir=`dirname $_path`

sed "s:{APP_PATH}:${_dir}:g" $_dir/scripts/init.d/imail.tpl > $_dir/scripts/init.d/imail
chmod +x $_dir/scripts/init.d/imail


sed "s:{APP_PATH}:${_dir}:g" $_dir/scripts/init.d/imail.service.tpl > $_dir/scripts/init.d/imail.service

if [ -d /etc/init.d ];then
	cp $_dir/scripts/init.d/imail /etc/init.d/imail
	chmod +x /etc/init.d/imail
fi

echo `dirname $_path`