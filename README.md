# Назначение
Данный Go-модуль реализует клиент к NamedPipe-API десктопного приложения 1С-Connect

# Cборка примера

```
git clone https://github.com/ros-tel/1c-connect-pipe.git 1c-connect-pipe
cd 1c-connect-pipe/examples/simple_client/
go install
```

Для кроскомпиляции предварительно добавить в переменные окружения GOOS=windows и требуемую архитектуру (GOARCH)
```
export GOOS=windows
export GOARCH="amd64" (или "386") 
```
