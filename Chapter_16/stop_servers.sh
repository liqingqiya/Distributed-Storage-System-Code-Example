ps auxf|grep server|grep hellofs|awk '{print $2}'|xargs -I {} kill -9 {}