
ssh-copy-id fpcyredis002
ssh-copy-id fpcyredis003
ssh-copy-id fpcyredis004
ssh-copy-id fpcyredis005
ssh-copy-id fpcyredis006
ssh-copy-id fpcyredis007
ssh-copy-id fpcyredis008
ssh-copy-id fpcyredis009
ssh-copy-id fpcyredis010
ssh-copy-id fpcyredis011
ssh-copy-id fpcyredis012
ssh-copy-id fpcyredis013
ssh-copy-id fpcyredis014
ssh-copy-id fpcyredis015
ssh-copy-id fpcyredis016
ssh-copy-id fpcyredis017
ssh-copy-id fpcyredis018
ssh-copy-id fpcyredis019
ssh-copy-id fpcyredis020
ssh-copy-id fpcyredis021
ssh-copy-id fpcyredis022
ssh-copy-id fpcyredis023
ssh-copy-id fpcyredis024
ssh-copy-id fpcyredis025
ssh-copy-id fpcyredis026
ssh-copy-id fpcyredis027

for h in $(cat hosts);do ssh fpcyredis027 ps -ef|grep redis;done

for h in $(cat hosts);do mkdir -p /app/fpcy/dedup/;done

for h in $(cat hosts);do scp dedup dedup.conf $h":/app/fpcy/dedup/";done

for h in $(cat hosts);do ssh $h chmod +x /app/fpcy/dedup/dedup;done

for h in $(cat hosts);do ssh $h "cd /app/fpcy/dedup;nohup /app/fpcy/dedup/dedup > /app/fpcy/dedup/dedup.log 2>&1 &";done

for h in $(cat hosts);do ssh $h ls -l /app/fpcy/dedup/;done

for h in $(cat hosts);do ssh $h pkill dedup;done

for h in $(cat hosts);do ssh $h "rm -f /app/fpcy/dedup/*.log";done


nohup /app/fpcy/dedup/dedup > /app/fpcy/dedup/dedup.log 2>&1 &

