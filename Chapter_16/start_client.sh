mkdir -p output/bin/
mkdir -p output/log/

nohup ./output/bin/client --mountpoint=/mnt/hellofs 2>&1 >> ./output/log/client_1.log &