# remindme
Remind Me On Recurring Deeds

Run:
```commandline
docker run -p 80:80 -d mpkondrashin/remindme
```

Update:
```commandline
docker pull mpkondrashin/remindme
docker stop remindme
docker rm remindme
docker run -p 80:80 --name remindme -d mpkondrashin/remindme
```