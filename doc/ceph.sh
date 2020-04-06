# crreate containere network  cephnet
docker network create --subnet 192.168.44.0/24 cephnet
# mon node
docker run -d --net cephnet --ip 192.168.44.11 -v /Users/rh01/temp/etc_ceph/:/etc/ceph -v /Users/rh01/temp/var_lib_ceph/:/var/lib/ceph -e MON_IP=192.168.44.11 -e CEPH_PUBLIC_NETWORK=192.168.44.0/24 --name mon --hostname mon  ceph/daemon:latest-mimic mon
# mgr node
docker run -d --net cephnet --ip 192.168.44.12 -v /Users/rh01/temp/etc_ceph/:/etc/ceph -v /Users/rh01/temp/var_lib_ceph/:/var/lib/ceph --name mgr --hostname mgr ceph/daemon:latest-mimic mgr
# osd
docker run -d --net cephnet --ip 192.168.44.13 -v /Users/rh01/temp/etc_ceph/:/etc/ceph -v /Users/rh01/temp/var_lib_ceph/:/var/lib/ceph --name osd --hostname osd -e OSD_TYPE=directory ceph/daemon:latest-mimic osd

