#!/bin/bash

mkdir -p /var/run/parceldrop
cp parceldrop /usr/bin/
cp parceldrop.service /etc/systemd/system/
#cp service.env /etc/parceldrop

