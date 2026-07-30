## GO Backend Course

This is just an exemplification of progress through using a project-base learning approach. This is my work, notes and approach to learning the basics of Backend development focusing on API RESTFUL architecture, middleware chain with auth, rate limiting, repository pattern, and PostgreSQL with migrations. The core idea is to simulate a real industry project and understand how backend works in production.

# Learning Notes

### Flujo del request

```
┌─────────────────────────────────────────────────────────────────┐
│                    FLUJO COMPLETO DE UN REQUEST                  │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  CLIENTE                                                         │
│    │                                                             │
│    │ POST /v1/posts {"content":"hola"}                          │
│    │ Authorization: Bearer <token>                              │
│    ▼                                                             │
│ ┌─────────────────────────────────────────────────────────────┐  │
│ │ 1. chi.Router recibe el request                             │  │
│ │    - Asigna un RequestID único                              │  │
│ │    - Logea el request (middleware.Logger)                   │  │
│ │    - Si hay panic, lo atrapa (Recoverer)                    │  │
│ │    - Si tarda >60s, lo corta (Timeout)                      │  │
│ └─────────────────────────────────────────────────────────────┘  │
│    │                                                             │
│    ▼                                                             │
│ ┌─────────────────────────────────────────────────────────────┐  │
│ │ 2. rateLimiterMiddleware                                    │  │
│ │                                                             │  │
│ │    Saca la IP del request → "192.168.1.5"                   │  │
│ │    Pregunta: Allow("192.168.1.5")?                          │  │
│ │                                                             │  │
│ │    ┌──────────────┐  ┌──────────────┐                       │  │
│ │    │ Balde lleno  │  │ Balde vacío  │                       │  │
│ │    │ (tokens > 1) │  │ (tokens < 1) │                       │  │
│ │    │ → SIGUE      │  │ → 429        │                       │  │
│ │    └──────────────┘  └──────────────┘                       │  │
│ └─────────────────────────────────────────────────────────────┘  │
│    │                                                             │
│    ▼                                                             │
│ ┌─────────────────────────────────────────────────────────────┐  │
│ │ 3. authMiddleware                                           │  │
│ │                                                             │  │
│ │    Extrae header: "Bearer <token>"                          │  │
│ │    Busca token en DB: auth_tokens WHERE token = $1          │  │
│ │                    AND expires_at > NOW()                   │  │
│ │                                                             │  │
│ │    ┌──────────────┐  ┌──────────────┐                       │  │
│ │    │ Token válido │  │ Token invál. │                       │  │
│ │    │ Busca user   │  │ → 401        │                       │  │
│ │    │ por user_id  │  └──────────────┘                       │  │
│ │    │              │                                         │  │
│ │    │ Guarda user  │                                         │  │
│ │    │ en context   │                                         │  │
│ │    │ → SIGUE      │                                         │  │
│ │    └──────────────┘                                         │  │
│ └─────────────────────────────────────────────────────────────┘  │
│    │                                                             │
│    ▼ (si la ruta requiere ownership)                             │
│ ┌─────────────────────────────────────────────────────────────┐  │
│ │ 4. requirePostOwnership (solo para PATCH/DELETE /posts/{id})│  │
│ │                                                             │  │
│ │    Toma {id} de la URL → Busca post en DB:                  │  │
│ │    SELECT * FROM posts WHERE id = $1                        │  │
│ │                                                             │  │
│ │    Compara: post.UserID == context.User.ID                   │  │
│ │                                                             │  │
│ │    ┌──────────────┐  ┌──────────────────┐                   │  │
│ │    │ Es dueño     │  │ No es dueño      │                   │  │
│ │    │ Guarda post  │  │ → 403            │                   │  │
│ │    │ en context   │  └──────────────────┘                   │  │
│ │    │ → SIGUE      │                                         │  │
│ │    └──────────────┘                                         │  │
│ └─────────────────────────────────────────────────────────────┘  │
│    │                                                             │
│    ▼                                                             │
│ ┌─────────────────────────────────────────────────────────────┐  │
│ │ 5. Handler (ej: createPostHandler)                          │  │
│ │                                                             │  │
│ │    Saca user del context                                     │  │
│ │    Decodifica JSON del body                                  │  │
│ │    Crea struct Post{UserID, Content}                         │  │
│ │    Llama: store.Posts().Create(ctx, post)                    │  │
│ │                                                             │  │
│ │    ┌────────────────────────────────────────────────────┐   │  │
│ │    │ store.Posts().Create():                            │   │  │
│ │    │   INSERT INTO posts (user_id, content)             │   │  │
│ │    │   VALUES ($1, $2) RETURNING id, created_at         │   │  │
│ │    └────────────────────────────────────────────────────┘   │  │
│ │                                                             │  │
│ │    Responde: 201 + JSON del post                            │  │
│ └─────────────────────────────────────────────────────────────┘  │
│    │                                                             │
│    ▼                                                             │
│  CLIENTE recibe {"id":1,"user_id":1,"content":"hola",...}       │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

### Interfaces -> Repo Pattern

```
┌─────────────────────────────────────────────────────────────┐
│                     ¿QUÉ CONECTA?                            │
│                                                              │
│  Handler necesita GUARDAR un post                            │
│       │                                                      │
│       ▼                                                      │
│  Llama a: store.Posts().Create(ctx, post)                    │
│       │                                                      │
│       ▼                                                      │
│  Posts() devuelve un PostRepository (la INTERFAZ)            │
│       │                                                      │
│       ▼                                                      │
│  La implementación real es PostStore (el SQL)                │
│                                                              │
│  ┌──────────────────────────────────────────────────────┐    │
│  │  INTERFAZ (storage.go)                               │    │
│  │  type PostRepository interface {                     │    │
│  │      Create(ctx, post) error                         │    │
│  │      GetByID(ctx, id) (*Post, error)                  │    │
│  │      Update(ctx, post) error                         │    │
│  │      Delete(ctx, id, userID) error                   │    │
│  │  }                                                   │    │
│  │                                                      │    │
│  │  func (s *Store) Posts() PostRepository {            │    │
│  │      return s.PostStore                              │    │
│  │  }                                                   │    │
│  │                                                      │    │
│  │  IMPLEMENTACIÓN (posts.go)                           │    │
│  │  type PostStore struct { db *pgxpool.Pool }          │    │
│  │  func (s *PostStore) Create(...) { SQL... }          │    │
│  └──────────────────────────────────────────────────────┘    │
└─────────────────────────────────────────────────────────────┘
```

### Middleware Chain

```
        ┌─────────────────────────────────────┐
        │  rateLimiterMiddleware              │
        │  ┌───────────────────────────────┐  │
        │  │  authMiddleware               │  │
        │  │  ┌─────────────────────────┐  │  │
        │  │  │  ownershipMiddleware    │  │  │
        │  │  │  ┌───────────────────┐  │  │  │
        │  │  │  │  handler          │  │  │  │
        │  │  │  │                   │  │  │  │
        │  │  │  │  (trabaja)        │  │  │  │
        │  │  │  └───────────────────┘  │  │  │
        │  │  └─────────────────────────┘  │  │
        │  └───────────────────────────────┘  │
        └─────────────────────────────────────┘
```


### Rule of Dependency
```
cmd/api/  ───→  internal/     ───→  base de datos
(HTTP)          (lógica)            (SQL)

NUNCA al revés. internal/ NO sabe que cmd/api/ existe.
Eso es ARQUITECTURA LIMPIA / HEXAGONAL.
```

### Stateless API (Token Auth)

El token es la única llave. El server no guarda "user x  está logueado" en memoria.
Cada request trae el token en el header, el server lo valida contra la DB.

    Login → server genera token → lo guarda en auth_tokens → lo devuelve
    Request → server busca token en DB → si existe y no expiró → deja pasar

Esto hace que cualquier instancia del server pueda validar cualquier request.
No importa a qué instancia caiga. No sticky sessions. No Redis obligatorio.


### Status Codes

| Código | Cuándo | Por qué |
|---|---|---|
| 200 | GET/POST exitoso | El recurso existe y se devuelve |
| 201 | POST que crea (user, post, comment) | El recurso fue creado |
| 204 | DELETE exitoso | No hay nada que devolver |
| 400 | JSON inválido o falta campo | El cliente mandó cualquier cosa |
| 401 | Token faltante o inválido | No estás autenticado |
| 403 | No sos el dueño del recurso | Estás autenticado pero no autorizado |
| 404 | Recurso no existe en DB | El ID no se encontró |
| 429 | Rate limit excedido | Muy rápido, esperá |
| 500 | Error inesperado del server | Algo explotó del lado nuestro |