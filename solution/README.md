<br />
<div align="center">
  <a href="https://github.com/github_username/repo_name">
    <img src="img/img.png" alt="Logo" height="80">
  </a>

<h3 align="center">Решение 3 этапа PROD</h3>

  <p align="center">
    Лучший рекламный движок на белом свете
  </p>
</div>



<a id="about_the_project"></a>
## О проекте

![Product Name Screen Shot](img/img_1.png)

Самый крутой движок для таргетированной рекламы, о котором вы только могли бы мечтать.

<a id="built_with"></a>
### Использованные технологии

* ![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white)
* ![Postgres](https://img.shields.io/badge/postgres-%23316192.svg?style=for-the-badge&logo=postgresql&logoColor=white)
* ![Redis](https://img.shields.io/badge/redis-%23DD0031.svg?style=for-the-badge&logo=redis&logoColor=white)
* ![Prometheus](https://img.shields.io/badge/Prometheus-E6522C?style=for-the-badge&logo=Prometheus&logoColor=white)
* ![Grafana](https://img.shields.io/badge/grafana-%23F46800.svg?style=for-the-badge&logo=grafana&logoColor=white)
* ![Docker](https://img.shields.io/badge/docker-%230db7ed.svg?style=for-the-badge&logo=docker&logoColor=white)
* ![YAML](https://img.shields.io/badge/yaml-%23ffffff.svg?style=for-the-badge&logo=yaml&logoColor=151515)
* ![Amazon S3](https://img.shields.io/badge/Amazon%20S3-FF9900?style=for-the-badge&logo=amazons3&logoColor=white)



<a id="getting_started"></a>
## С чего начать?

<a id="prerequisites"></a>
### Требования

Убедитесь, что у вас установлен docker и docker compose, больше ничего не потребуется. Ниже приведен пример установки для Ubuntu

  ```sh
# Add Docker's official GPG key:
sudo apt-get update
sudo apt-get install ca-certificates curl
sudo install -m 0755 -d /etc/apt/keyrings
sudo curl -fsSL https://download.docker.com/linux/ubuntu/gpg -o /etc/apt/keyrings/docker.asc
sudo chmod a+r /etc/apt/keyrings/docker.asc

# Add the repository to Apt sources:
echo \
  "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.asc] https://download.docker.com/linux/ubuntu \
  $(. /etc/os-release && echo "${UBUNTU_CODENAME:-$VERSION_CODENAME}") stable" | \
  sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
sudo apt-get update
  ```

```shell
sudo apt-get install docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin
```

```shell
sudo docker run hello-world
```

<a id="installation"></a>
### Запуск

```shell
docker compose up
```
Вуа-ля! Бэкенд доступен на порту 8080.

Unit тесты:

```shell
go test -v ./...
```

После запуска тг бот будет доступен [здесь](https://t.me/prod_warl0rdd_bot)

<a id="see_also"></a>
## Смотрите также

Вся документация расписана в .md файлах, находящихся в папке solution/doc

* [Схема базы данных](doc/db.md)
* [Обоснование использованных технологий](doc/proof.md)
* [Описание docker compose файла](doc/compose.md)
* [Описание работы алгоритма](doc/algorithm.md)
* [Описание эндпоинтов](doc/usage.md)
