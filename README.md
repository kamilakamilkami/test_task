# Запуск проекта

1. **Сборка Docker**  
```bash
docker compose build --no-cache
docker compose up
```

2. Тестирование можно через 3 юзеров
После поднятия контейнеров перходить надо на api/v1/auth/login

- **EMPLOYEE**: kamilafake@gmail.com / fakepassword  
- **ADMIN**: rizafake@gmail.com / fakepassword  
- **MANAGER**: alishfake@gmail.com / fakepassword
