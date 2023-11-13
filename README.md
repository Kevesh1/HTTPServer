# Lab 1 Group 42

Brief description of the project.

## Starting

### local env
In a local environment you only need to run the command 
```
go run .
```

After this you need to give the server the port 8080
and the proxt 8081
This will result in the following:
```
(base) ➜  HTTPServer git:(main) ✗ go run .
Server: Enter what port to listen from: 
8080
Server: Running on port:  8080
Proxy: Enter what port to start proxy server from: 
8081
Proxy: Running on port:  8081
Proxy: BEFORE HANDLE
```

### docker env
For a docker environment you will run two commands.
One for building the image and the other to run the container

```
docker build --tag http-server .
```
```
docker run -e MAIN-PORT=8080 -e PROXY-PORT=8081 --publish 8080:8080 --publish 8081:8081 httpserver
//This will expose the containers 8080 port to our local port 8080 and same with 8081
//and will set the .env variables to our desired ports in the go servers
```


## Testing

### POST
To test the functionality of the proxy server and server we can curl files
from the ./test folder. Unfortunatly readable files are limited to that folder.

When curling for example test.jpg to the server it will be stored in the ./files folder.

If there was an error with a request the error will be displayed in either the server or client window
depending on what error you create.

### GET 
You can test this either through the proxy adress or the server adress.
The results will be the same. A file can only be fetched if it exists in the ./files folder
