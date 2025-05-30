# Sharing Warehouse Machines

# Build application
```
sudo docker network create custom_netw
docker compose up
```

# PostgreSQL DB (old)
**How to up project databse**:
```
docker compose up -d
```

Run app:
```
go build cmd/app/main.go
./main --config=config/local.yml
```

# Use Cases

## Frontend API Examples

`curl` - утилита для отправки запросов. Если нету `curl`, то отправить запрос можно любым другим способом.
Curl flags which used for request:
- `-d` added to request body data
- `-H` added header to request 
- '-X' added http method to request

### User authorization

```
curl -d '{"phone_number": "88889997766", "password": "boss-password321"}' -X POST "localhost:8080/login"
```
**Success** response example:
```
{"token":"..."}
```

### GET DB Data Methods 
These methods allows only for `Admin` users.
```
# пример получения всех пользователей/машин/сессий
curl -H "Authorization: Bearer <user-token>" -X GET "localhost:8080/get_all_users"
curl -H "Authorization: Bearer <user-token>" -X GET "localhost:8080/get_all_machines"
curl -H "Authorization: Bearer <user-token>" -X GET "localhost:8080/get_all_sessions"

# примеры получения объектов по их id
curl -H "Authorization: Bearer <user-token>" -X GET "localhost:8080/get_user?user_id=1"
curl -H "Authorization: Bearer <user-token>" -X GET "localhost:8080/get_machine?machine_id=1FGH345"
curl -H "Authorization: Bearer <user-token>" -X GET "localhost:8080/get_session?session_id=2"
```

## Unlock/Lock Machine
**Unlock machine example**:
```
curl -H "Authorization: Bearer <user-token>" -d '{"machine_id": "<machine-id>"}' -X POST "localhost:8080/unlock_machine"
```

Success response:
```
{"sessionId":3}
```

Error response:
```
{"error":"<error-msg>"}
```

**Lock machine example**:
```
curl -H "Authorization: Bearer <user-token>" -d '{"machine_id": "<machine-id>"}' -X POST "localhost:8080/lock_machine"
```

Success response:
```
{"msg":"successfullly lock machine"
```

Error response:
```
{"error":"<error-msg>"}
```

## Register Microcontroller
Request to register new machine in system:
```
curl -d '{"machine_id": "NEWM123", "ip_addr": "123.456.78.90"}' -X POST "localhost:8080/register_machine"
```

Response:
```
{"current_state":<state>}
```

state: 0 | 1 
If machine was turned off in session-process then after machine register session will be restore.

# Добавление пользователей в базу данных
Изначально в базе данных нету информации. В веб клиенте не предусмотрена возможность добавления новых пользователей в систему.

Полная инструкция по добавлению пользоваеля доступна здесь [docs/add_user.md](./docs/add_user.md)
