# Contributing to Financial App

¡Gracias por tu interés en contribuir a Financial App! Este documento te guiará a través del proceso de configuración del proyecto en tu máquina local.

## Requisitos Previos

- Go 1.23.9 o superior
- Git
- Una cuenta de Supabase (para la base de datos)
- Node.js y npm (para dependencias de desarrollo, si es necesario)

## Configuración del Proyecto

### 1. Clonar el Repositorio

```bash
git clone https://github.com/tu-usuario/Financial_app.git
cd Financial_app
```

### 2. Configurar Variables de Entorno

Crea un archivo `.env` en la raíz del proyecto con las siguientes variables:

```env
SUPABASE_URL=tu_url_de_supabase
SUPABASE_KEY=tu_clave_secreta_de_supabase
```

### 3. Instalar Dependencias

El proyecto utiliza Go Modules para la gestión de dependencias. Las dependencias se descargarán automáticamente al compilar el proyecto.

En la terminal :
```bash
npm i
```
y despues :
```bash 
npx supabase login
```
y para configurar la conexion a la base de datos
```bash
npx supabase link --project-ref <you-projec-ref>
```

para subir migraciones :
```bash
npx supabase db push
```

### 4. Ejecutar la Aplicación

Para iniciar el servidor de desarrollo:

```bash
go run main.go
```

La aplicación estará disponible en `http://localhost:8080`.

## Estructura del Proyecto

```
Financial_app/
├── Domains/           # Definición de dominios y puertos
│   └── ports/        # Interfaces para casos de uso y repositorios
├── Models/            # Modelos de datos
├── UseCases/          # Lógica de negocio
├── infrastructure/    # Implementaciones concretas (ej: conexión a Supabase)
├── interfaces/        # Controladores y rutas de la API
│   ├── controllers/   # Controladores de la API
│   └── dtos/         # Objetos de transferencia de datos
├── docs/              # Documentación de la API (Swagger)
└── main.go           # Punto de entrada de la aplicación
```

## Convenciones de Código

- **Nombrado**: Usa nombres descriptivos en inglés.
- **Formato**: Sigue el formato estándar de Go (`gofmt`).
- **Comentarios**: Documenta funciones y tipos exportados.
- **Mensajes de Commit**: Sigue el formato convencional de commits.

## Proceso de Contribución

1. Crea un fork del repositorio
2. Crea una rama para tu característica (`git checkout -b feature/amazing-feature`)
3. Haz commit de tus cambios (`git commit -m 'Add some amazing feature'`)
4. Haz push a la rama (`git push origin feature/amazing-feature`)
5. Abre un Pull Request

## Documentación de la API

La documentación de la API está disponible en formato Swagger. Después de iniciar la aplicación, puedes acceder a:

- Documentación Swagger UI: `http://localhost:8080/swagger/index.html`
- Esquema Swagger JSON: `http://localhost:8080/swagger/doc.json`

## Pruebas

Para ejecutar las pruebas del proyecto:

```bash
go test ./...
```

## migraciones

generar con :
```bash
npx supabase migration new your-migration-name  
```

esto generara un archivo sql en `/supabase/migrations`.

publicar cambio con :
```bash
npx supabase db push      
```

## Soporte

Si tienes preguntas o necesitas ayuda, por favor abre un issue en el repositorio.

## Licencia

Este proyecto está bajo la licencia MIT. Ver el archivo `LICENSE` para más detalles.
