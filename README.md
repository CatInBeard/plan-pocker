# Plan-pocker

Scrum plan pocker web app

## Run local env
Local env using docker compose  

### Commands
To build and run:  
```
make rebuild
```  
To stop  
```
make down
```  
To start  
```
make up  
```
To migrate db  
```  
make migrate
```   

### First build
At first build you should run
```
make rebuild
```

After all containers started try to run
```
make migrate
```
to update database


Application will start at http://localhost, admin panel run at http://localhost/admin/, default user `admin`, password `password`


## Demo
You can try application at https://plan.newpage.xyz