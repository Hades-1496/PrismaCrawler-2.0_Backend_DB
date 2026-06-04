# PrismaCrawler 2.0 - Backend API 🐉

Backend genérico orientado a juegos Dungeon Crawler / Roguelike, desarrollado en **Golang**. 
Este proyecto sirve como introducción al desarrollo de APIs en Go, enfocándose en un rendimiento robusto, acceso a base de datos y gestión de estados de partida.

## 🎯 Objetivo
Proveer una API sólida para gestionar usuarios, personajes, inventarios y el estado del ciclo de juego (runs/partidas) interactuando con una base de datos PostgreSQL alojada en Supabase.

## 🛠️ Stack Tecnológico Recomendado

Para cumplir con los plazos (10 días) y asegurar la mejor experiencia de desarrollo:

- **Lenguaje**: Go (Golang) 1.21+
- **Framework Web / Enrutador**: [Gin](https://gin-gonic.com/) (Rápido, excelente documentación y fácil de aprender).
- **Base de Datos**: PostgreSQL (via Supabase).
- **ORM (Equivalente a Prisma)**: [Ent](https://entgo.io/) (Desarrollado por Facebook, utiliza generación de código y esquemas basados en grafos, muy similar a Prisma) o alternativamente **GORM** para una curva de aprendizaje más rápida.
- **Autenticación**: JWT / Supabase Auth.

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

## 📐 Arquitectura del Proyecto (Layered / Capas)

Para mantener la simplicidad sin sacrificar el orden, utilizaremos una arquitectura de carpetas estándar en Go:

```text
PrismaCrawler/
├── cmd/
│   └── api/             # Punto de entrada de la aplicación (main.go)
├── internal/
│   ├── handlers/        # Controladores (HTTP/Gin). Reciben la petición y devuelven JSON
│   ├── services/        # Lógica de negocio (Validar que un jugador puede equipar un item)
│   ├── models/          # Entidades y esquemas (Ent/GORM)
│   └── repository/      # Capa de acceso a la base de datos (Querys)
├── pkg/                 # Código reutilizable (Helpers de JWT, Configuración, etc)
├── go.mod               # Dependencias
└── README.md
```

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

## ✅ Lista de Requisitos (MVP para Entrega)

- [ ] **Setup y Conexión**: Conectar Go con Supabase (PostgreSQL) exitosamente.
- [ ] **Migraciones**: Poder generar las tablas en la BD usando Ent o GORM.
- [ ] **Auth API**: Endpoints para registrar y loguear un usuario (`/auth/register`, `/auth/login`).
- [ ] **Character API**: CRUD de personajes (`/characters`).
- [ ] **Game State API**: 
  - Iniciar una nueva partida (`POST /runs/start`).
  - Guardar el progreso / piso actual (`PUT /runs/{id}/save`).
  - Morir / Finalizar partida (`POST /runs/{id}/die`).

## 🚀 Instalación y Ejecución

1. Clonar el repositorio.
2. Renombrar `.env.example` a `.env` y configurar la URL de Supabase.
3. Instalar dependencias:
   ```bash
   go mod tidy
   ```
4. Ejecutar en modo desarrollo:
   ```bash
   go run cmd/api/main.go
   ```
