# GO and MongoDB

This project is a set API to perform CRUD operations on a MongoDB cluster using the Go Programming Language with JWT authorization.

## Upload json dataset into database (movies.json)

```bash
mongoimport --uri mongodb+srv://<ROOT>:<PASSWORD>@<CLUSTER_NAME>.cmlur.mongodb.net/<DATABASE> --collection <COLLECTION> --type <FILETYPE> --file <FILENAME> --jsonArray
```

## Run the project

Execute this command in your terminal:

```bash
go run main.go
```

```bash
You can use the GET, POST, PUT, DELETE methods to make requests to the localhost and perform the basic CRUD operations on your database Cluster.
```

## Folder Structure

- handlers : Contains all the functions for handling the request made to the localhost.
- helper : Used to check the authentication of a user.
- models : Custom native stucts to store the data.
- utils : Utility fuunctions for database connection and other stuff.
- main.go : point of execution
- movies_data.json : dummy movies data used for this project.

## Resources

While implementing these API's, i got a lot of help from Youtube tutorials and StackOverflow Q/A's

https://www.youtube.com/watch?v=JNr5noDp6EM&t=119s (whole series) : [Nic Raboy]\
https://www.youtube.com/watch?v=hWmR8YtlFlE (JWT implementation with GO) : [Daily Code Buffer]\

https://docs.mongodb.com/manual/reference/method \
https://golang.org/doc/
