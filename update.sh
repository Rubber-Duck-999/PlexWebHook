#!/bin/sh

cd $HOME/Documents/HouseGuard-NetworkAccessController/src

git pull

go clean

go build

if [ -f exeNetworkAccessController ];
then
    echo "NAC File found"
    if [ -f $HOME/Documents/Deploy/exeNetworkAccessController ];
    then
        echo "NAC old removed"
        rm -f $HOME/Documents/Deploy/exeNetworkAccessController
    fi
    mv exeNetworkAccessController $HOME/Documents/Deploy/exeNetworkAccessController
fi