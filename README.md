# Prueba Técnica

Version Go: 1.24.4

## Docker
Para crear el contenedor y poner en funcionamiento los servicios ejecutar el comando
```
 docker-compose up --build
```

## API
### Servicio de autenticación

Este servicio recibe un nombre y regresa un token temporal de maximo 5 usos y con una duración maxima de 10 minutos

#### Generar token

```
  POST http://localhost:8080/token
```

| Parámetro | Tipo     | Descripción                |
| :-------- | :------- | :------------------------- |
| `name` | `string` | Ingresa un nombre |

Ejemplo: Petición en formato JSON
```
{
    "name":"Fabi"
}
```

Y su respuesta es en JSON
```
{
    "token": "0b80f7c7-de32-40f0-8245-181ffb9a4c1f"
}
```

### Servicio de consulta de datos

Este servicio retorna desde la API publica The Rick and Morty un listado de personajes, pero necesita el token generado previamente, este token se pone como un Bearer Token

#### Obtener personajes

```
  GET http://localhost:8080/characters
```

| Parámetro | Tipo     | Descripción                |
| :-------- | :------- | :------------------------- |
| `token` | `string` | Ingresa un token valido |

Y su respuesta es en JSON
```
[
    {
        "id": 1,
        "name": "Rick Sanchez",
        "image": "https://rickandmortyapi.com/api/character/avatar/1.jpeg"
    },
    {
        "id": 2,
        "name": "Morty Smith",
        "image": "https://rickandmortyapi.com/api/character/avatar/2.jpeg"
    },
    {
        "id": 3,
        "name": "Summer Smith",
        "image": "https://rickandmortyapi.com/api/character/avatar/3.jpeg"
    },
    ....
]

```

## Tests

Si se desea ejecutar los tests el comando es el siguiente

```
go test ./tests/ -v
```


## Autor

Fabián Alejandro Otálora Silva

Desarrollador de Software 🇨🇴

[@fabian-otalora](https://www.github.com/fabian-otalora)

