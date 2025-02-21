# install uninstall old version of extension
echo Uninstalling old version of extension...
code --uninstall-extension walrus.walrus
# compile
./fmt.sh
# pack
echo Packing language syntax
cd language-support
vsce package
./syn.sh
# install
echo Installing new version of extension...
code --install-extension walrus.walrus-1.0.15.vsix
cd ..
echo Done