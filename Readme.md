# Financial App

Aplicación de gestión financiera desarrollada en Go con Gin y Supabase.

## Características

- Gestión de usuarios
- Manejo de cuentas financieras
- API RESTful documentada con Swagger
- Autenticación con JWT
- Base de datos en Supabase

## Requisitos

- Go 1.23.9 o superior
- Cuenta en Supabase
- Git

## Instalación

1. Clona el repositorio:
   ```bash
   git clone https://github.com/tu-usuario/Financial_app.git
   cd Financial_app
   ```

2. Configura las variables de entorno:
   ```bash
   cp .env.example .env
   # Edita el archivo .env con tus credenciales
   ```

3. Instala las dependencias:
   ```bash
   go mod download
   ```

## Uso

Inicia el servidor de desarrollo:

```bash
go run main.go
```

La aplicación estará disponible en `http://localhost:8080`

## Documentación de la API

Accede a la documentación interactiva en:
- Swagger UI: http://localhost:8080/swagger/index.html
- Esquema JSON: http://localhost:8080/swagger/doc.json

## Estructura del Proyecto

```
Financial_app/
├── Domains/           # Definición de dominios y puertos
├── Models/            # Modelos de datos
├── UseCases/          # Lógica de negocio
├── infrastructure/    # Implementaciones concretas
├── interfaces/        # Controladores y rutas
└── docs/              # Documentación Swagger
```

## Contribución

Por favor lee [CONTRIBUTING.md](CONTRIBUTING.md) para detalles sobre nuestro código de conducta y el proceso para enviar pull requests.

## Licencia

Este proyecto está bajo la Licencia MIT - ver el archivo [LICENSE](LICENSE) para más detalles.