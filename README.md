# How to launch
- To run it locally:
    1. `make install` to install all dependencies.
    2. `make` to compile the source code.
    3. `GIN_MODE=release config=config/prod/parameters.yaml ./bin/main` to start the server.
- To run it in docker:
    1. you need to have docker installed.
    2. run `make deploy-build` to build the image.
    3. use `make deploy` to start a container instance.

# Features
1. The server is support external configuration file to define the running port and storage root.
2. Endpoints:
    1. `GET    /healthcheck` Health check endpoint to indicate whether server is live
    2. `POST   /file`        User can post a file via here, use form: path to provide file dir, file to upload file. P.S. user can not upload file content via this endpoint.
    3. `PUT    /file`        User can create/update file content here, the file target dir and content is defined in body (JSON format). Example: 
    ```
    {
        "filePath": "anyPath",
        "fileName": "anyFileName",
        "fileContent": "this is a file"
    }
    ```
    4. `GET    /file?path=<file_dir>`    Get file content.
    5. `DELETE /file?path=<file_dir>`    Delete target file. User is not allowed to delete a folder if it is not empty.
    6. `GET    /static?path=<folder_dir/file_dir>`  Get static information about target file/folder. P.S. path parameter can be omitted and storage root path will take as the default path parameter. 
3. The server will prevent write/delete/read file to the level upper than given storage root, it will take the file name and write under storage root.
4. Docker storage folder volume is shared with `./storage` folder in host. So the file change can be persisted.
5. The static information is gathered on demand but highly concurrently.

# Future Improves
- Split server (router and handler), file operation code and statics code into different repositories.
- Better endpoint and rest method organize for create/update file via upload via file upload and upload via request body.  
- Give 400 error if user request static for a non-exist file/folder.
- An endpoint to give file structures of a given dir. 
- User credentials.
- Speed:
    1. We can store static information in a separate meta file aside with files and load them into memory, so when user request the static information it can be gathered much quicker. The event of update static files can relates with the hit of the relative file operation endpoints. 
- Security:
    1. Host docker repository ourselves.
    2. Host glide.sh ourselves.
    3. Should restrict to read the binary files such as jpeg, pdf, zip