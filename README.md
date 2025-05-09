# H-project

### Локальный запуск

1. Запуск локально
```
make run
```


2. Взять файл .env.example из deployments/deployment и создать свой .env


3. Запуск через docker-compose 

```
cd deployments/development
docker-compose up -d
```

Задание
- Принимать данные в api -> записать в бд postgres -> сохранить на s3 файлик -> потом отправить и вычитать consumer -> и отдать
- Описание команий(определиться с форматом и составить JSON)

Api Handler
    - /companies/import (POST JSON)
    -  usecase companies в GO (бизнес логика) валидация данных
    - Общий FLOW
- Принять в JSON HANDLER описание компании и передать в слой use case - разпарсить и валидиировать


### S 3
    -  принять файл с json(побайтово marshall)
    - minio s 3 библиотека
    - поднять minio в docker-compose
    - s3 client 
    - внутри minio - mc help
    - информация о компании в s3 по id


- использовать slog для логирования
- заранее создать топик для получения данных по файлу с s3
- ibm-sarama-golang - библиотека - клиент для кафки