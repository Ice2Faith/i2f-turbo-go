AppName=dist
echo begin deploy ${AppName} ...

_p_now=$(date "+%Y%m%d%H%M%S")
_p_bak=${AppName}.${_p_now}
_p_bak_dir=backup.${AppName}

echo backup ...
mkdir -p ${_p_bak_dir}
mv ${AppName} ${_p_bak_dir}/${_p_bak}

echo unzip ...
unzip ${AppName}.zip -d ${AppName} > /dev/null

echo verify directory...
if [ -d "${AppName}/${AppName}" ];then
    rm -rf tmp.${AppName}
    mv ${AppName} tmp.${AppName}
    mv tmp.${AppName}/${AppName} .
    rm -rf tmp.${AppName}
fi

echo done.
