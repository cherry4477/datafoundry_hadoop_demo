#!/usr/bin/env bash
echo "kinit"
date
kinit $BSI_HDFS_HDFSDEMO_USERNAME <<!!
$BSI_HDFS_HDFSDEMO_PASSWORD
!!

echo "start main"
date
tail -f /dev/null
