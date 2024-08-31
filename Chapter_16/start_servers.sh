mkdir -p output/bin/
mkdir -p output/data/
mkdir -p output/data/hellofs_data_{1..6}/
mkdir -p output/log/

nohup ./output/bin/server --datapath="./output/data/hellofs_data_1/" --port=37001 2>&1 >> ./output/log/server_1.log &
nohup ./output/bin/server --datapath="./output/data/hellofs_data_2/" --port=37002 2>&1 >> ./output/log/server_2.log &
nohup ./output/bin/server --datapath="./output/data/hellofs_data_3/" --port=37003 2>&1 >> ./output/log/server_3.log &
nohup ./output/bin/server --datapath="./output/data/hellofs_data_4/" --port=37004 2>&1 >> ./output/log/server_4.log &
nohup ./output/bin/server --datapath="./output/data/hellofs_data_5/" --port=37005 2>&1 >> ./output/log/server_5.log &
nohup ./output/bin/server --datapath="./output/data/hellofs_data_6/" --port=37006 2>&1 >> ./output/log/server_6.log &