#!/bin/sh

cd $HOME/Documents/HouseGuard-NetworkAccessController/src

git pull

go clean

go build

if [ -f exeNetworkAccessController ];
then
    echo "FH File found"
    if [ -f $HOME/Documents/Deploy/exeNetworkAccessController ];
    then
        echo "FH old removed"
        rm -f $HOME/Documents/Deploy/exeNetworkAccessController
    fi
    mv exeNetworkAccessController $HOME/Documents/Deploy/exeNetworkAccessController
fi