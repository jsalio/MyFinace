# Contexto del Proyecto: Aplicación de Finanzas Personales

## Descripción General
Aplicación backend desarrollada en Go (Golang) para gestionar finanzas personales, con autenticación de usuarios y gestión de billeteras digitales.

## Estructura del Proyecto

### 1. Dominio (Domains/)
- **Puertos (Ports/)**: Define las interfaces principales del sistema
  - `Repository.go`: Interfaz base para operaciones de repositorio
  - `UserUseCase.go`: Casos de uso para gestión de usuarios
  - `WalletUseCase.go`: Casos de uso para gestión de billeteras
  - `WalletRepositoryExtend.go`: Extensiones específicas para repositorio de billeteras

### 2. Modelos (Models/)
- **Base de Datos**
  - `db/User.go`: Estructura de datos para usuarios
  - `db/wallet.go`: Estructura de datos para billeteras

- **DTOs (Data Transfer Objects)**
  - `dtos/auth_request.go`: Estructura para autenticación
  - `dtos/create_account_request.go` y `create_account_response.go`: Para creación de cuentas
  - `dtos/create_wallet_request.go`: Para creación de billeteras
  - `dtos/update_*_request.go`: Para actualización de entidades
  - `dtos/error_response.go` y `success_response.go`: Respuestas estándar

### 3. Interfaces (interfaces/)
- **Controladores**
  - `controllers/`: Manejo de peticiones HTTP
    - `auth_controller.go`: Autenticación y registro
    - `wallet_controller.go`: Gestión de billeteras
    - `account_controller.go`: Gestión de cuentas
    - `base_controller.go`: Controlador base con funcionalidad común

- **Middleware**
  - `middleware/auth.go`: Autenticación JWT
  - `middleware/AuthConfig.go`: Configuración de autenticación

### 4. Casos de Uso (UseCases/)
- Implementación de la lógica de negocio

### 5. Infraestructura (infrastructure/)
- Implementación de repositorios y conexiones a base de datos

## Características Implementadas

### Autenticación y Autorización
- Registro de usuarios
- Inicio de sesión con JWT
- Middleware de autenticación

### Gestión de Billeteras
- Creación de billeteras
- Actualización de información
- Eliminación de billeteras

### Gestión de Cuentas
- Creación de cuentas
- Actualización de información
- Eliminación de cuentas

## Tecnologías Utilizadas
- **Lenguaje**: Go (Golang)
- **Framework Web**: Gin
- **Documentación**: Swagger/OpenAPI
- **Autenticación**: JWT (JSON Web Tokens)

## Documentación de la API
La documentación interactiva de la API está disponible en `/swagger/index.html` cuando el servidor está en ejecución.

## Estado Actual
- Estructura base del proyecto establecida
- Autenticación implementada
- Gestión básica de billeteras y cuentas
- Documentación de API con Swagger

## Próximos Pasos
1. Implementar operaciones financieras (transferencias, depósitos, retiros)
2. Añadir categorización de transacciones
3. Implementar reportes y análisis financiero
4. Añadir más validaciones de seguridad
5. Mejorar el manejo de errores

## Configuración
El servidor se ejecuta por defecto en `http://localhost:8085`

## Requisitos
- Go 1.16 o superior
- Base de datos configurada (ver configuración en `infrastructure/`)