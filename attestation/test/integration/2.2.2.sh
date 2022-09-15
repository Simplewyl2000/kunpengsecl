#!/bin/bash
# this scripts should be run under the root folder of kunpengsecl project
#set -eux
PROJROOT=.
# run number of rac clients to test
NUM=1
# include common part
. ${PROJROOT}/attestation/test/integration/common.sh

# above are common preparation steps, below are specific preparation step, scope includs:
# configure files, input files, environment variables, cmdline paramenters, flow control paramenters, etc.
### Start Preparation
echo "start test preparation..." | tee -a ${DST}/control.txt
pushd $(pwd)
cd ${PROJROOT}/attestation/quick-scripts
echo "clean database" | tee -a ${DST}/control.txt
sh clear-database.sh | tee -a ${DST}/control.txt
popd
### End Preparation

### define some constant
strUUID="9b954212d796863e9f2c04372d4ab7e39fe0b62870c82a9e83c3ec326e5fb9b9"
strCONTAINER="container"
strNAME="testContainer"
strOLDIMA="ima 5a2842c1767f26defc2e96a01e46062524333501 /home/test/abc\n"
strNEWIMA="ima 0000000000000000000000000000000000000000 /home/test/abc\n"

### start launching binaries for testing
echo "start ras..." | tee -a ${DST}/control.txt
( cd ${DST}/ras ; ./ras -T &>${DST}/ras/echo.txt ; ./ras -v &>>${DST}/ras/echo.txt ;)&
echo "sleep 5s" | tee -a ${DST}/control.txt
sleep 5

# start number of rac 
echo "start ${NUM} rac clients..." | tee -a ${DST}/control.txt
(( count=0 ))
for (( i=1; i<=${NUM}; i++ ))
do
    ( cd ${DST}/rac-${i} ; ${DST}/rac/raagent -t -v &>>${DST}/rac-${i}/echo.txt ; )&
    (( count++ ))
    if (( count >= 1 ))
    then
        (( count=0 ))
        echo "start ${i} rac clients at $(date)..." | tee -a ${DST}/control.txt
    fi
done

### start monitoring and control the testing
echo "start to perform test ..." | tee -a ${DST}/control.txt
echo "wait for 5s" | tee -a ${DST}/control.txt
sleep 5

# get cid
echo "get client id" | tee -a ${DST}/control.txt
cid=$(awk '{ if ($1 == "clientid:") { print $2 } }' ${DST}/rac-1/config.yaml)
echo ${cid} | tee -a ${DST}/control.txt

# add container basevalue
AUTHTOKEN=$(grep "Bearer " ${DST}/ras/echo.txt)
echo "start adding new container basevalue which uuid is ${strUUID}..." | tee -a ${DST}/control.txt
curl -X POST -k -H "Authorization: $AUTHTOKEN" -H "Content-Type: application/json" https://localhost:40003/${cid}/newbasevalue --data "{\"uuid\":\"${strUUID}\", \"basetype\":\"${strCONTAINER}\", \"name\":\"${strNAME}\", \"enabled\":true, \"ima\":\"${strOLDIMA}\"}"
echo "wait for 5s" | tee -a ${DST}/control.txt
sleep 5

# get bid
RESPONSE=$(curl -k -H "Content-Type: application/json" https://localhost:40003/${cid}/basevalues)
echo "get basevalue id" | tee -a ${DST}/control.txt
bid=$(echo ${RESPONSE} | jq -r '.' | grep -B 3 ${strUUID} | awk '/"ID"/ {gsub(",","",$2);print $2}')
echo ${bid} | tee -a ${DST}/control.txt

# get the container basevalue for the first time
BASECONTAINER1=$(curl -k -H "Content-Type: application/json" https://localhost:40003/${cid}/basevalues/${bid})
ENABLED1=$(echo ${BASECONTAINER1} | jq -r '.' | awk '/"Enabled"/ {gsub(",","",$2);print $2}')
CONTAINERIMA1=$(echo ${BASECONTAINER1} | jq -r '.' | awk '/"Ima"/ {gsub(",","",$2);print $2}')
if [[ $ENABLED1 == true ]]
then
    echo "set container named ${strNAME} basevalue succeeded!" | tee -a ${DST}/control.txt
    echo "get the container ima value is ${CONTAINERIMA1}" | tee -a ${DST}/control.txt
    echo "take the next step..." | tee -a ${DST}/control.txt
else
    echo "set container named ${strNAME} basevalue failed!" | tee -a ${DST}/control.txt
    pkill -u ${USER} ras
    pkill -u ${USER} raagent
    echo "test DONE!!!" | tee -a ${DST}/control.txt
    exit 1
fi

# post basevalue
echo "start posting container basevalue..." | tee -a ${DST}/control.txt
curl -X POST -k -H "Authorization: $AUTHTOKEN" -H "Content-type: application/json" https://localhost:40003/${cid}/basevalues/${bid} --data '{"enabled":false}'
curl -X POST -k -H "Authorization: $AUTHTOKEN" -H "Content-Type: application/json" https://localhost:40003/${cid}/newbasevalue --data "{\"uuid\":\"${strUUID}\", \"basetype\":\"${strCONTAINER}\", \"name\":\"${strNAME}\", \"enabled\":true, \"ima\":\"${strNEWIMA}\"}"
echo "wait for 5s" | tee -a ${DST}/control.txt
sleep 5

# get bid
RESPONSE2=$(curl -k -H "Content-Type: application/json" https://localhost:40003/${cid}/basevalues)
echo "get basevalue id" | tee -a ${DST}/control.txt
bid2=$(echo ${RESPONSE2} | jq -r '.' | grep -B 3 ${strUUID} | awk '/"ID"/ {gsub(",","",$2);print $2}' | sed -n '2p')
echo ${bid2} | tee -a ${DST}/control.txt

# get the container basevalue for the second time
BASECONTAINER2=$(curl -k -H "Content-Type: application/json" https://localhost:40003/${cid}/basevalues/${bid})
ENABLED2=$(echo ${BASECONTAINER2} | jq -r '.' | awk '/"Enabled"/ {gsub(",","",$2);print $2}')
BASECONTAINER3=$(curl -k -H "Content-Type: application/json" https://localhost:40003/${cid}/basevalues/${bid2})
ENABLED3=$(echo ${BASECONTAINER3} | jq -r '.' | awk '/"Enabled"/ {gsub(",","",$2);print $2}')
CONTAINERIMA2=$(echo ${BASECONTAINER3} | jq -r '.' | awk '/"Ima"/ {gsub(",","",$2);print $2}')
if [[ $ENABLED2 == false && $ENABLED3 == true ]]
then
    echo "modify container named ${strNAME} basevalue succeeded!" | tee -a ${DST}/control.txt
    echo "get the container ima value is ${CONTAINERIMA2}" | tee -a ${DST}/control.txt
    echo "take the next step..." | tee -a ${DST}/control.txt
else
    echo "modify container named ${strNAME} basevalue failed!" | tee -a ${DST}/control.txt
    pkill -u ${USER} ras
    pkill -u ${USER} raagent
    echo "test DONE!!!" | tee -a ${DST}/control.txt
    exit 1
fi

### stop testing
echo "kill all test processes..." | tee -a ${DST}/control.txt
pkill -u ${USER} ras
pkill -u ${USER} raagent

echo "test SUCCEEDED!!!" | tee -a ${DST}/control.txt
exit 0
