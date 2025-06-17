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

## Deploying in Docker Swarm

To deploy the application in Docker Swarm, follow these steps:

1. **Initialize Docker Swarm**  
   First, you need to initialize Docker Swarm. Run the command:
   ```
    docker swarm init
   ```

2. **Plan the Build**  
   Next, execute the script to plan the build:
   ```
    ./build-cloud-docker.sh plan-pocker git
   ```

3. **Create the Environment File**  
   Create an environment file based on the example located in `cloud/test.env`. Copy it and edit as needed:
   ```
    cp cloud/test.env /path/to/file.env
   ```

4. **Deploy the Application**  
   After that, run the script to deploy the application, specifying the path to your environment file:
   ```
    ./swarm-deploy.sh up envfile=/path/to/file.env
   ```

Once these steps are completed, the application will start deploying in Docker Swarm and will be accessible on the specified port.


## Demo
You can try application at https://plan.newpage.xyz