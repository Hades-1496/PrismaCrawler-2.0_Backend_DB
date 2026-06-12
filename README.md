# PrismaCrawler 2.0 - Backend API 🐉

Backend genérico orientado a juegos Dungeon Crawler / Roguelike, desarrollado en **Golang**. 
Este proyecto sirve como introducción al desarrollo de APIs en Go, enfocándose en un rendimiento robusto, acceso a base de datos y gestión de estados de partida.

## 📑 Índice
1. [Objetivo](#-objetivo)
2. [Stack Tecnológico](#️-stack-tecnológico-recomendado)
3. [Comparativa JS vs Go](#-comparativa-evolución-de-nodejs-a-go)
4. [Arquitectura y Flujo de Datos](#-arquitectura-del-proyecto-layered--capas)
5. [Decisiones Arquitectónicas](#-decisiones-arquitectónicas)
6. [Clean Code Aplicado](#-principios-de-clean-code-aplicados)
7. [Rutas de la API (Endpoints)](#️-rutas-de-la-api-endpoints)
8. [Esquema de Base de Datos](#️-esquema-de-base-de-datos-core-5-tablas)
9. [Instalación y Ejecución](#-instalación-y-ejecución)
10. [Despliegue y Testing](#-despliegue-en-producción-servidor-linux--vps)

## 🎯 Objetivo
Proveer una API sólida para gestionar usuarios, personajes, inventarios y el estado del ciclo de juego (runs/partidas) interactuando con una base de datos PostgreSQL alojada en Supabase.

## 🛠️ Stack Tecnológico Recomendado

Para cumplir con los plazos (10 días) y asegurar la mejor experiencia de desarrollo:

- **Lenguaje**: Go (Golang) 1.21+
- **Framework Web / Enrutador**: [Gin](https://gin-gonic.com/) (Rápido, excelente documentación y fácil de aprender).
- **Base de Datos**: PostgreSQL (via Supabase).
- **ORM (Equivalente a Prisma)**: [Ent](https://entgo.io/) (Desarrollado por Facebook, utiliza generación de código y esquemas basados en grafos, muy similar a Prisma) o alternativamente **GORM** para una curva de aprendizaje más rápida.
- **Autenticación**: JWT / Supabase Auth.

### 📦 Dependencias Integradas
Hasta el momento se han añadido al proyecto:
- `github.com/gin-gonic/gin` para el servidor y enrutamiento HTTP.
- `github.com/joho/godotenv` para cargar variables de entorno seguras.
- `gorm.io/gorm` y `gorm.io/driver/postgres` para la conexión y el ORM de la base de datos.
- `golang.org/x/crypto/bcrypt` para la encriptación segura de contraseñas.
- `github.com/golang-jwt/jwt/v5` para la generación y validación de tokens de sesión.
- `golang.org/x/time/rate` para proteger la API con limitación de peticiones (Rate Limiting).

## 🌟 Comparativa: Evolución de Node.js a Go

Esta versión en Go no es solo una traducción del código original, sino una **evolución arquitectónica** hacia un verdadero motor de estado persistente:

| Característica | Node.js (Anterior) | Golang (Nuevo) |
| :--- | :--- | :--- |
| **Paradigma** | "Arcade" (Stateless, guarda score al morir) | "Roguelike" (Guarda progreso piso a piso) |
| **Identidad** | `User` == Jugador / Avatar | Se separa la cuenta (`User`) de los héroes (`Character`) |
| **Mapas** | Matriz ASCII estática en BD | Generación Procedural delegada por `Seed` aleatoria |
| **Inventario** | Simple catálogo de consulta | El backend gestiona qué llevas equipado (`run_inventory`) |
| **Rendimiento** | Single-thread (Bloqueos en CPU intensiva) | Compilado, multi-hilo nativo (`Goroutines`) rapidísimo |
| **Código** | Scripting monolítico | Tipado estricto, Clean Code, Inyección de Dependencias |
| **Arquitectura** | Todo en un mismo servidor | Go actúa de Orquestador, delegando la IA a microservicios |

## 📐 Arquitectura del Proyecto (Layered / Capas)
El proyecto utiliza una arquitectura por capas. El siguiente diagrama muestra cómo fluye una petición y cómo se comunican los componentes:

```text
   [ Frontend / Cliente Phaser ]
                 |
            (Petición HTTP)
                 v
           [ Router (Gin) ] -----------> Verifica Middlewares (Auth, CORS, Ping)
                 |
                 v
     [ Handlers (Controladores) ] -----> Valida el JSON y extrae IDs
                 |
                 v
       [ Services (Negocio) ] ---------> Aplica lógica del juego y lanza eventos a la IA (Observer)
                 |
                 v
   [ Repository (Acceso a Datos) ] ----> Prepara las consultas abstrayendo el motor de BD
                 |
                 v
         [ GORM (Motor ORM) ]
                 |
                 v
      [ PostgreSQL (Supabase) ]
```
### Estructura de Carpetas
```text
PrismaCrawler/
├── cmd/
│   └── api/             # Punto de entrada de la aplicación (main.go)
├── internal/
│   ├── handlers/        # Controladores (HTTP/Gin). Reciben la petición y devuelven JSON
│   ├── middlewares/     # Interceptores de seguridad (Auth, Rate Limiter, CORS)
│   ├── services/        # Lógica de negocio (Ej: Validaciones antes de guardar partida)
│   ├── models/          # Entidades de la Base de Datos (GORM)
│   ├── repository/      # Capa de abstracción de base de datos
│   └── routes/          # Centralización del enrutamiento
├── pkg/                 # Código reutilizable (Helpers)
│   ├── aiclient/        # Cliente HTTP para el microservicio de IA
│   ├── db/              # Conexión a la base de datos y seeder
│   └── utils/           # Utilidades genéricas (JWT, random seed)
├── tests/               # Pruebas automatizadas (TDD)
├── .env.example         # Plantilla de variables de entorno
└── README.md
```

## 🧠 Decisiones Arquitectónicas y Consideraciones (Para el Equipo)

Al tener un tiempo límite de **10 días** y ser una primera toma de contacto con Go, hemos definido las siguientes directrices:

### 1. PostgreSQL vs GraphQL
Es común confundir estos términos, pero cumplen roles completamente distintos y no son excluyentes:
- **PostgreSQL**: Es el *Motor de Base de Datos Relacional*. Es donde se guardan físicamente los datos (usuarios, items, partidas).
- **GraphQL**: Es un *Lenguaje de Consultas para APIs* (una alternativa a REST). Es la forma en la que el cliente (Frontend) le pide los datos al Backend.

**Nuestra decisión:** Supabase utiliza **PostgreSQL** por debajo, por lo que esa será nuestra base de datos. Supabase provee GraphQL automáticamente de cara al frontend, pero nuestro Backend en Go se conectará directamente a PostgreSQL para gestionar la lógica pesada. Expondremos endpoints REST (con Gin) porque en un plazo de 10 días, configurar un servidor GraphQL propio en Go añadiría una complejidad innecesaria para el MVP.

### 2. Gestión de Base de Datos: El equivalente a Prisma en Go
Prisma es excelente en TypeScript, pero en el ecosistema Go tenemos alternativas muy sólidas. Debemos elegir la que mejor se adapte a nuestro tiempo:
- **Ent:** (Recomendado a largo plazo). Es el más similar a Prisma. Define esquemas como código (basado en grafos) y genera un cliente fuertemente tipado. *Contra:* Configurar la generación de código lleva algo de tiempo.
- **GORM:** (Recomendado para el MVP de 10 días). El ORM más popular en Go. Configuración casi instantánea basada en *structs* (clases).
- **sqlc:** Escribes SQL puro y genera código Go seguro. Increíble rendimiento pero requiere dominio de SQL.

**Nuestra decisión:** Utilizaremos **GORM** (o Ent si el equipo asume la curva de aprendizaje) para acelerar el desarrollo del MVP y centrarnos en la lógica del juego.

### 3. Alcance Realista del Juego (Restricciones del MVP)
El backend no debe sobrecargarse de lógica en tiempo real para esta entrega inicial. Go es rápido, pero los WebSockets requerirían demasiada infraestructura.

**Lo que NO haremos:**
- Combate calculado en tiempo real en el servidor.
- Movimiento casilla por casilla validado por red (síncrono).

**Lo que SÍ haremos (Enfoque Transaccional):**
- El cliente/frontend maneja el *gameplay* de la mazmorra. Nuestro backend en Go actuará como el "servidor de guardado y validación": Autenticará al usuario, entregará los datos del personaje al inicio de la *Run*, y recibirá actualizaciones del progreso (Ej: "El jugador superó el piso 3, guarda este estado y este nuevo inventario").

## 🧼 Principios de Clean Code Aplicados

Para asegurar la mantenibilidad y escalabilidad del proyecto, se han aplicado los siguientes principios de diseño de software:

- **KISS (Keep It Simple, Stupid)**: Se ha evitado la sobre-ingeniería. Por ejemplo, se centralizaron las rutas en un solo archivo en lugar de crear una estructura de carpetas excesiva para el tamaño del MVP.
- **DRY (Don't Repeat Yourself)**: Se han creado funciones `helper` (como `utils.SendError` y `utils.GetUserID`) para evitar la duplicación de código en los controladores.
- **YAGNI (You Ain't Gonna Need It)**: Se han evitado capas innecesarias para CRUDs simples. Por ejemplo, en lugar de crear un `EconomyService` complejo, el handler gestiona las peticiones simples del monedero del jugador.
- **SoC (Separation of Concerns)**: Se han separado las responsabilidades, aislando toda la configuración de endpoints en el paquete `routes` y delegando eventos asíncronos en middlewares u observers.
- **SRP (Single Responsibility Principle)**: Cada componente tiene una única razón para cambiar. Por ejemplo, la lógica de negocio (cómo muere un personaje) se movió de los controladores HTTP al propio modelo `GameRun`, dejando que el controlador solo se ocupe del transporte de datos.
- **LoD (Law of Demeter)**: Se ha evitado el acoplamiento excesivo. Por ejemplo, en lugar de que un controlador acceda a `run.Character.UserID`, se ha creado un método `run.IsOwnedBy(userID)` para que el controlador solo "hable" con el objeto `run`.
- **OCP (Open/Closed Principle)**: Mediante el patrón Observer (`GameObserver`), el servicio principal de partidas puede notificar a la IA y a otros futuros microservicios de los hitos del juego sin tener que modificar su código base.
- **DIP (Dependency Inversion Principle)**: Los módulos de alto nivel (handlers) no dependen de los de bajo nivel (repositorios), sino de abstracciones (interfaces). Esto se logra mediante la Inyección de Dependencias en `main.go`.
- **Testing (TDD) y Fiabilidad**: Se implementó una suite de pruebas unitarias y comprobaciones de health check (`/ping`) para facilitar los despliegues en producción.

## 🛣️ Rutas de la API (Endpoints)
### Autenticación (Públicas):
- POST /auth/register: Registra un nuevo usuario (requiere email y password).
- POST /auth/login: Inicia sesión y devuelve un token JWT.

### Generales (Públicas):
- GET /api/leaderboard: Devuelve el Top 10 de mejores partidas globales.
- GET /ping: Health check para verificar el estado del servidor.
- GET /game/items: Alias de retrocompatibilidad para el catálogo de objetos en Phaser.

### Core del Juego (Protegidas por JWT en /api):
- POST /api/characters: Crea un nuevo personaje (Guerrero, Mago) para el usuario logueado.
- GET /api/characters: Obtiene la lista de personajes del usuario activo.
- GET /api/characters/:id: Obtiene los detalles de un personaje específico.
- PUT /api/characters/:id: Actualiza el nombre/clase de un personaje.
- DELETE /api/characters/:id: Elimina un personaje.
- POST /api/runs/start: Inicia una partida, verificando que el personaje esté vivo, y genera la Seed procedural.
- PUT /api/runs/save: Actualiza el progreso de la partida (piso, score, vida restante) o mata al personaje si HP <= 0.
- GET /api/runs: Historial de todas las partidas del usuario.
- GET /api/runs/:id: Detalles de una partida específica y los objetos equipados.
- GET /api/profile: Obtiene los datos del usuario logueado y su Top 5 de mejores partidas.
- GET /api/items: Devuelve el catálogo completo de objetos del juego.
- GET /api/enemies: Devuelve el catálogo completo (bestiario) de enemigos.
- GET /api/maps: Lista los mapas generados del juego.
- GET /api/maps/:id: Detalles de un mapa específico para el frontend.
- GET/PUT /api/wallet: Persistencia offline para las monedas y gemas.
- GET/PUT /api/garden: Persistencia offline para el minijuego de las plantas.

### Administración (Protegidas por JWT y Rol ADMIN en /api/admin):
- PUT /api/admin/role: Cambia el rol de un usuario (requiere `user_id` y `role`).
- POST /api/admin/discord/changelog: Envía notas del parche al servidor de Discord.
- POST /api/admin/discord/test-webhook: Verifica la conectividad con el bot de Discord.
- POST /api/admin/knowledge: Añade nueva información al cerebro RAG.
- GET /api/admin/knowledge: Lista la base de conocimientos almacenada.
- DELETE /api/admin/knowledge/:id: Borra un fragmento de conocimiento.

### Comunicación Interna (Microservicios / IA):
*(Rutas protegidas por header `X-Internal-Token`, exclusivas para el Backend de IA)*
- GET /api/internal/rankings/top10: Devuelve el top 10 para el scheduler de Discord.
- POST /api/internal/rag/search: Búsqueda vectorial semántica (`pgvector`).
- GET /api/internal/stats: Estadísticas globales (jugadores, partidas, muertes).
- GET /api/internal/items: Catálogo de objetos exportado para análisis de economía.
- GET /api/internal/enemies: Bestiario exportado para el observador de la IA.
- GET /api/internal/wallet/:user_id: Consulta del monedero por ID.
- GET /api/internal/garden/:user_id: Consulta del estado del jardín por ID.

## 🗄️ Esquema de Base de Datos (Core 5 Tablas)
Para un Dungeon Crawler genérico en un plazo realista, necesitamos 5 tablas principales:

1. **`users` (Usuarios)**
   - `id`, `email`, `password_hash`, `created_at`, `last_login`
   - *Nota: Si usas Supabase Auth, esta tabla puede sincronizarse con sus usuarios.*
2. **`characters` (Personajes / Héroes)**
   - `id`, `user_id`, `name`, `class` (Guerrero, Mago...), `level`, `base_hp`, `is_alive`
3. **`items` (Catálogo global del juego - Diccionario)**
   - `id`, `name`, `description`, `type` (Weapon, Armor, Potion), `stats_modifier` (JSONB)
4. **`game_runs` (Partidas / Ciclo Roguelike)**
   - `id`, `character_id`, `seed` (Para generación procedural), `current_floor`, `score`, `status` (In_Progress, Won, Dead)
5. **`run_inventory` (Inventario de la partida actual)**
   - `id`, `run_id`, `item_id`, `quantity`, `equipped`
6. **`knowledge_chunks` (Base de Conocimiento / RAG)**
   - `id`, `keywords`, `content`

## ✅ Lista de Requisitos (MVP para Entrega)

- [x] **Setup y Conexión**: Conectar Go con Supabase (PostgreSQL) exitosamente.
- [x] **Migraciones**: Poder generar las tablas en la BD usando Ent o GORM.
- [x] **Auth API**: Endpoints para registrar y loguear un usuario (`/auth/register`, `/auth/login`).
- [x] **Character API**: CRUD de personajes (`/api/characters`) y perfiles de usuario.
- [x] **Game State API**: 
  - Iniciar una nueva partida (`POST /api/runs/start`).
  - Guardar el progreso / piso actual y morir si HP <= 0 (`PUT /api/runs/save`).
  - Consultar Leaderboard (`GET /api/leaderboard`).

## 🗃️ Archivo `.gitignore`
Antes de hacer un commit, asegúrate de tener un archivo `.gitignore` en la raíz del proyecto para no exponer secretos ni subir código compilado:
```text
# Secretos y Variables de Entorno (¡NUNCA SUBIR EL .env!)
.env

# Binarios compilados de Go
*.exe
*.out
/bin/
```

## 🚀 Instalación y Ejecución

1. Clonar el repositorio.
2. Renombrar `.env.example` a `.env` y configurar la `DATABASE_URL` y el `PORT`.
3. Instalar dependencias:
   ```bash
   go mod tidy
   ```
4. Ejecutar en modo desarrollo:
   ```bash
   go run cmd/api/main.go
   ```
 - 💡 Nota sobre el comando: A diferencia de JS donde ejecutas el script de entrada directamente (y herramientas como nodemon lo reinician), go run compila temporalmente toda tu aplicación en memoria y genera un binario para ejecutarla. La ruta cmd/api/main.go apunta a tu archivo principal. Al ser Go un lenguaje compilado, si modificas cualquier archivo .go, deberás detener el servidor en la terminal (Ctrl+C) y volver a ejecutar el comando para aplicar los cambios.

## 🌍 Despliegue en Producción (Servidor Linux / VPS)

Una de las mayores ventajas de Go frente a Node.js es que no necesitas instalar el lenguaje ni descargar dependencias (`node_modules`) en el servidor de producción. Todo se empaqueta en un único archivo binario nativo.

### Paso 1: Compilar para Linux (Cross-Compilation)
En tu ordenador (Windows/Mac), abre la terminal en la raíz del proyecto y ejecuta:
```bash
GOOS=linux GOARCH=amd64 go build -o prismacrawler cmd/api/main.go
```
Esto generará un archivo ejecutable llamado `prismacrawler` en la raíz de tu proyecto.

### Paso 2: Subir al servidor
Sube ese único archivo `prismacrawler` y tu archivo `.env` (con las URLs de tu Supabase de producción) al servidor Linux usando herramientas como FileZilla o mediante el comando `scp`. No hace falta subir nada más.

### Paso 3: Ejecutar como Servicio (systemd)
Para que la API nunca se apague (incluso si el servidor se reinicia o tú cierras la terminal), crea un servicio en Linux.
En el servidor, crea el archivo `/etc/systemd/system/prismacrawler.service`:
```ini
[Unit]
Description=PrismaCrawler API en Go
After=network.target

[Service]
User=root
WorkingDirectory=/ruta/a/tu/carpeta
ExecStart=/ruta/a/tu/carpeta/prismacrawler
Restart=always

[Install]
WantedBy=multi-user.target
```
Actívalo con `sudo systemctl enable prismacrawler` y enciéndelo con `sudo systemctl start prismacrawler`.

## 🧪 Testing (TDD)
Para ejecutar la suite de pruebas unitarias y de integración (E2E), asegúrate de haber configurado `TEST_DATABASE_URL` en tu `.env` y ejecuta:
```bash
go test ./tests/... -v
```