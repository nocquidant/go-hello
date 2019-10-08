mage_version="1.9.0"
mage_url="https://github.com/magefile/mage/releases/download/v${mage_version}/mage_${mage_version}_Linux-64bit.tar.gz"

curl -L -s -o mage.tar.gz ${mage_url} 
if [ $? -ne 0 ]; then exit 1; fi

tar xzf mage.tar.gz
cp mage /usr/local/bin/mage
chmod +x /usr/local/bin/mage