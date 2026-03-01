# CallFlowManager - Sistema de Call Center

Sistema profesional de gestión de call center construido con Go, MongoDB y Bootstrap 5.

## Características

- **Programación de Llamadas**: Agenda llamadas con clientes
- **Gestión de Agentes**: Administra tu equipo de agentes
- **Gestión de Clientes**: CRM de clientes
- **Dashboard**: Métricas en tiempo real
- **Historial**: Registro de todas las llamadas

## Tecnologías

- Go 1.21, MongoDB 7.0, Bootstrap 5, Chart.js
- Server Side Rendering con html/template
- JWT Authentication

## Requisitos

- Go 1.21+
- MongoDB 7.0 (local o remoto)

## Instalación y Ejecución Local

```bash
# Clonar el repositorio
git clone <repo-url>
cd PruebaGit

# Asegúrate de que MongoDB esté ejecutándose en localhost:27017
# O actualiza MONGODB_URI en el archivo .env

# Instalar dependencias
go mod tidy

# Ejecutar el servidor
go run cmd/server/main.go

# La aplicación estará disponible en http://localhost:8084
```

### Con Docker (Alternativo)

```bash
docker-compose up --build
```

## Colecciones MongoDB

- agents - Agentes del call center
- customers - Clientes
- calls - Llamadas programadas
- call_logs - Historial de llamadas
- schedules - Horarios

## API Endpoints

### Autenticación
- `POST /api/auth/register` - Registrar usuario
- `POST /api/auth/login` - Iniciar sesión

### Agentes
- `GET /api/agents` - Listar agentes
- `POST /api/agents` - Crear agente

### Clientes
- `GET /api/customers` - Listar clientes
- `POST /api/customers` - Crear cliente

### Llamadas
- `GET /api/calls` - Listar llamadas
- `POST /api/calls` - Programar llamada
- `PUT /api/calls/:id/status` - Actualizar estado
- `GET /api/calls/stats` - Estadísticas

## Estructura del Proyecto

```
/
├── cmd/server/           # Punto de entrada
├── internal/
│   ├── config/         # Configuración
│   ├── handlers/       # Controladores HTTP
│   ├── services/      # Lógica de negocio
│   ├── repositories/  # Acceso a datos
│   ├── models/        # Modelos de datos
│   ├── middlewares/   # Middlewares HTTP
│   └── utils/         # Utilidades
├── web/
│   ├── templates/     # Plantillas HTML
│   └── static/        # CSS, JS
├── .env
├── .env.example
├── Dockerfile
├── docker-compose.yml
└── go.mod
```

## Seguridad

- JWT Authentication
- Roles: admin, user
- Hash de contraseñas con bcrypt

---

## 👨‍💻 Desarrollado por Isaac Esteban Haro Torres

**Ingeniero en Sistemas · Full Stack · Automatización · Data**

- 📧 Email: zackharo1@gmail.com
- 📱 WhatsApp: 098805517
- 💻 GitHub: https://github.com/ieharo1
- 🌐 Portafolio: https://ieharo1.github.io/portafolio-isaac.haro/

---

© 2026 Isaac Esteban Haro Torres - Todos los derechos reservados.
