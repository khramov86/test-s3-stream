# Тестирование потоковой передачи из S3

Как пользоваться 

Генерация файла
```bash
dd if=/dev/zero of=10gb_file bs=1G count=10
```
Поднимаем

```bash
docker compose build
docker compose up -d
```

Заходим в minio http://localhost:9000

Вводим логин/пароль

Создаем bucket

Cоздаем access + secret keys

Загружаем сгенерированный файлик

Запускаем сбор статистики
```bash
docker compose stats
```
Скачиваем
```
curl "localhost:8080/download?file=10gb_file" -O /temp/
```

Видно, что утилизация по памяти не поднимается
