git pull
cd $HOME/hub_intra_proto/src/intra-hub
bee pack
sudo mkdir -p /var/www/beego
sudo rm -rfv /var/www/beego/*
sudo mv -v intra-hub.tar.gz /var/www/beego
cd /var/www/beego
sudo tar -xvf intra-hub.tar.gz
sudo rm -v intra-hub.tar.gz
sudo pkill intra-hub
echo '' > $HOME/nohup.out
sudo nohup ./intra-hub&
