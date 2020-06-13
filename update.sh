#!/bin/sh

cd $HOME/Documents/HouseGuard-NetworkAccessController/src

git pull

go clean

go build

if [ -f exeNetworkAccessController ];
then
    echo "FH File found"
    if [ -f $HOME/Documents/Temp/exeNetworkAccessController ];
    then
        echo "FH old removed"
        rm -f $HOME/Documents/Temp/exeNetworkAccessController
    fi
    mv exeNetworkAccessController $HOME/Documents/Temp/exeNetworkAccessController
fi