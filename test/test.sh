user=`date +%s%N | md5sum | head -c 10`
echo "start to operate on " $user
sleep 1

out=`curl -s -X POST -F file=@scripts.sh localhost:8080/$user/task`
echo "task id created is: " $out 

curl localhost:8080/$user/task/$out
echo "\n"

curl -X PUT localhost:8080/$user/task/$out
echo "trigger build\n"

curl localhost:8080/$user/task/$out/files
echo "\n"

curl localhost:8080/$user/task/$out/files/output.txt
echo "\n"

out2=`curl -s -X POST -F file=@scripts.sh localhost:8080/$user/task?output=$user-0.0.1.iso`
echo "task id created is: " $out2

curl localhost:8080/$user/task/$out2
echo "\n"
curl localhost:8080/$user/task
echo "\n"

curl -s -X DELETE localhost:8080/$user/task/$out2
echo "\n"
curl localhost:8080/$user/task/$out2
echo "\n"
curl -s -X DELETE localhost:8080/$user/task/$out2
echo "\n"

echo "\n"
