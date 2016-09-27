#!/usr/bin/env bash
echo "kinit"
date
kinit $BSI_HDFS_HDFSDEMO_USERNAME <<!!
$BSI_HDFS_HDFSDEMO_PASSWORD
!!

echo "36.110.131.65 hadoop-1.jcloud.local" >> /etc/hosts
echo "36.110.132.55 hadoop-2.jcloud.local" >> /etc/hosts
echo "start main"
date
tail -f /dev/null
