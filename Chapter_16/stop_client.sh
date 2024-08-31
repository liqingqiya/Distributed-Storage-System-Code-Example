
umount /mnt/hellofs/

ps auxf|grep client|grep hellofs|awk '{print $2}'|xargs -I {} kill -9 {}