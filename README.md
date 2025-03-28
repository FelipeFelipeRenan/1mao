# 1mao
This is the official 1mao repo


## How to run


### Optional

If you want to use Swagger documentation

Fist install swag utilitary using 
```go install github.com/swaggo/swag/cmd/swag@latest```

After the successful instalation fo Swag, run the command to generate the swagger documentation 

```swag init -g cmd/main.go```

Run the command ```docker compose up --build``` to build and run the application

access the url [link](http://localhost:8080/swagger/index.htm) to use Swagger documentation endpoints