# 🔍 Análisis del Backend Original (Node.js + Express + Prisma)

Este documento detalla toda la funcionalidad que poseía el backend antiguo escrito en JavaScript, sirviendo como mapa de ruta para nuestra migración a Go.

## 🏗️ Arquitectura General
El proyecto utilizaba una arquitectura en capas clásica (Layered Architecture):
- **Rutas (Routes)**: Definían los endpoints y apuntaban a los controladores.
- **Controladores (Controllers)**: Recibían la petición (`req`, `res`), extraían los datos y llamaban a los servicios.
- **Servicios (Services / Logic)**: Contenían toda la lógica de negocio pesada y hacían las consultas a la base de datos usando Prisma ORM.
- **Middlewares**: Interceptores para autenticación y manejo de errores.

## 🛡️ Middlewares Implementados
1. **`authMiddleware.js`**: Protegía rutas leyendo la cabecera `Authorization`. Extraía el token JWT (quitando la palabra "Bearer "), lo verificaba con `jsonwebtoken` y guardaba el `userId` en `req.user` para su uso posterior.
2. **`errorHandler.js`**: Un manejador global que capturaba cualquier error de la aplicación, unificaba el formato de respuesta (con un `statusCode` y un `message`) y evitaba que el servidor crasheara, devolviendo un JSON limpio al frontend.
3. **CORS y JSON**: En `index.js` se utilizaban los middlewares de Express para parsear `body` (JSON) y permitir peticiones desde orígenes cruzados (CORS).

## 👥 Módulo de Usuarios (`userController` & `userLogic`)
Encargado de la gestión de identidades.

- **Registro (`POST /api/users/register`)**:
  - Comprueba si el email ya existe.
  - Encripta la contraseña usando `bcryptjs` (salt 10).
  - Crea el usuario en BD y genera un JWT de 24h.
- **Login (`POST /api/users/login`)**:
  - Busca el usuario y compara la contraseña encriptada.
  - Devuelve un JWT y los datos públicos del usuario.
- **Perfil (`GET /api/users/profile`)**: *(Ruta protegida)*
  - Obtiene los datos del usuario logueado usando el token.
  - **Característica clave:** Traía automáticamente el "Top 5" de las mejores puntuaciones de ese usuario (`scores` ordenados por `floor` y `xp`).

## ⚔️ Módulo de Juego (`gameController` & `gameLogic`)
Encargado del core del juego y las estadísticas.

- **Obtener Mapa (`GET /api/game/map/:id`)**:
  - Busca un mapa específico en BD.
  - **Importante:** Parsea los campos `layout` y `dictionary` que estaban guardados como cadenas JSON en la base de datos, convirtiéndolos en Arrays y Objetos utilizables para Phaser.
- **Leaderboard Global (`GET /api/game/top-scores`)**:
  - Trae el Top 10 histórico de todas las partidas de la base de datos.
  - El orden de jerarquía era: Piso más alto (`floor` DESC) -> Mayor Experiencia (`xp` DESC) -> Más Kills (`kills` DESC).
- **Catálogo de Items (`GET /api/game/items`)**:
  - Devuelve la lista completa de objetos para que el frontend (Phaser) los reconozca.
- **Guardar Puntuación / Morir (`POST /api/game/score`)**: *(Ruta protegida)*
  - Asocia una nueva estadística final al usuario logueado. Guarda piso, kills, daño hecho, daño recibido y experiencia.
- **Crear Mapas (`POST /api/game/map`)**: *(Supuestamente para Administradores)*
  - Recibe un array de strings (layout) y un diccionario y los inyecta en la BD convirtiéndolos previamente a `JSON.stringify()`.

## 🌱 Seeding y Base de Datos (`seed.js`)
El script de reseteo (`seed.js`) hacía un trabajo fundamental de preparación:

1. **Limpieza**: Vaciaba las tablas de `Score`, `Item` y `Map`.
2. **Usuarios Semilla**:
   - Creaba 3 usuarios con contraseña "admin123":
     - *Kilian Admin* (Rol ADMIN).
     - *Guts The Slayer* (Rol USER - El tryhard).
     - *Sir Pupas* (Rol USER - El novato).
3. **Puntuaciones de Prueba**: 
   - Añadía 3 records al salón de la fama para estos usuarios (Guts en el piso 15, Sir Pupas en el piso 1).
4. **Items**: 
   - Creaba 5 objetos clave con sus `spriteKeys` para Phaser y arrays de efectos:
     - *Espada de Hierro* (+Daño).
     - *Cerveza de Haste* (+Velocidad de ataque).
     - *Poción de Salud* (Curación).
     - *Té Enfermizo* (+Velocidad).
     - *Sack of Weight* (Más HP, Menos velocidad).
5. **Diccionarios de Mapas**:
   - Creaba **6 niveles** completos.
   - Cada mapa usaba una **matriz ASCII** (ej: `####DD####`) donde `#` era pared, `_` vacío, `M` enemigo, `P` jugador, `T` loot.
   - Un diccionario adjunto traducía esos caracteres a entidades del juego (Ej: `M` = Slime con 30 HP).

---

## 📌 Diferencias y Evolución para Go (PrismaCrawler 2.0)
Observando este código antiguo, en nuestra nueva versión en Go estamos haciendo algunas mejoras arquitectónicas:
1. **Estado persistente (Runs)**: En JS solo se guardaba la puntuación *al morir*. En Go estamos guardando la partida *piso a piso* (`GameRun`), lo que permitirá reconectarse si se cierra el navegador.
2. **Inventario de la partida**: En JS los items existían como catálogo, pero los personajes no se los guardaban en base de datos. En Go, hemos creado `run_inventories` para que la partida recuerde qué lleva equipado.
3. **Personajes (Characters)**: En JS, un Usuario era directamente el jugador. En Go hemos separado `User` (la cuenta) de `Character` (el héroe), permitiendo que una cuenta tenga varios héroes distintos (ej: Un Mago y un Guerrero).