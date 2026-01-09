Список запросов в сервер чтоб получить:

Cписок всех студентов
curl http://localhost:8080/students

Данные о студенте по id, тут это 1
curl http://localhost:8080/students/1

Cписок всех груп
curl http://localhost:8080/groups

Расписание определенной групы по id, тут это 2
curl http://localhost:8080/schedule/group/2

Cписок всего расписания
curl http://localhost:8080/schedule
