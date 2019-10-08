mage_version="1.9.0"
mage_url="https://github.com/magefile/mage/releases/download/v${mage_version}/mage_${mage_version}_Linux-64bit.tar.gz"

mkdir -p ~/bin && curl -L -s -o ~/bin/mage.tar.gz ${mage_url} 
if [ $? -ne 0 ]; then exit 1; fi

tar xzf ~/bin/mage.tar.gz -C ~/bin && chmod +x ~/bin/mage
export PATH=$PATH:~/bin

mage --version