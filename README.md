# Golang Gin Gorm Basket Api

## How to run

### Required

- Mysql

### Conf

You should modify `conf.yml`

```
database:
  dbname: "dbname"
  username: "db_username"
  password: "db_password"
  host: "db_host"
  testdbname: "dbname_test"

server:
  port: ":app_port"
  secret: "jwt_secret_token_key"
```

### Run

```
$ go run main.go
```

### Base Url

```
http://localhost:8080
```

## Routes

Api routes are listed below.

| METHOD | ROUTE            | EXPLANATION           |    AUTHORIZATION     | POST PARAMS                                                       |
| :----: | :--------------- | :-------------------- | :------------------: | :---------------------------------------------------------------- |
| `GET`  | /swagger/\*any   | API Documentation     |                      |                                                                   |
| `POST` | /user/add        | Add a New User        |                      | { "username":"data1", "password":"data2" }                        |
| `POST` | /user/login      | User Login            |                      | { "username":"data1", "password":"data2" }                        |
| `POST` | /category/add    | Add a New Category    | ✔️ (only admin role) | CSV file                                                          |
| `POST` | /category/       | Get All Categories    |                      | { "page":int1, "limit":int2 }                                     |
| `POST` | /product/add     | Add a New Product     | ✔️ (only admin role) | { "name":"data1", "brand":"data2", "categoryId":int1 }            |
| `POST` | /product/        | Get All Products      |                      | { "page":int1, "limit":int2 }                                     |
| `POST` | /product/search  | Product Search        |                      | { "value":"data1" }                                               |
| `POST` | /product/delete  | Product Delete        | ✔️ (only admin role) | { "id":int1 }                                                     |
| `POST` | /product/update  | Product Update        | ✔️ (only admin role) | { "name":"data1", "brand":"data2", "categoryId":int1, "id":int2 } |
| `POST` | /shopping/add    | Add or Update Cart    |          ✔️          | [ { "productId":int1, "quantity":int2 } ]                         |
| `GET`  | /shopping/list   | Get the Cart          |          ✔️          |                                                                   |
| `POST` | /shopping/delete | Delete or Update Cart |          ✔️          | [ { "productId":int1, "quantity":int2 } ]                         |
| `POST` | /order/add       | Create Order          |          ✔️          |                                                                   |
| `GET`  | /order/list      | Get All Orders        |          ✔️          |                                                                   |
| `POST` | /shopping/cancel | Cancel the Order      |          ✔️          | { "orderId":int1 }                                                |
