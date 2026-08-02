import { useState, useEffect } from 'react'
import './App.css'

function App() {
  const [token, setToken] = useState(localStorage.getItem('token') || '')
  const [user, setUser] = useState(null)
  const [view, setView] = useState(localStorage.getItem('token') ? 'feed' : 'login')

  //login
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')

  //Register
  const [username, setUsername] = useState('')

  //post
  const [posts, setPosts] = useState([])
  const [content, setContent] = useState('')
  const [imageFile, setImageFile] = useState(null)


  const handleLogin = async (e) => {
    e.preventDefault()
    const res = await fetch('/v1/login', {
      method: 'POST',
      headers: {'Content-Type': 'application/json'},
      body: JSON.stringify({email, password})
    })

    const data = await res.json()
    if (res.ok) {
      setToken(data.token)
      localStorage.setItem('token', data.token)
      setView('feed')
    } else {
      alert(data.error)
    }
  }


  const handleRegister = async (e) => {
    e.preventDefault()
    const res = await fetch('/v1/users', {
      method: 'POST',
      headers: {'Content-Type': 'application/json'},
      body: JSON.stringify({username, email, password})
    })
    if (res.ok) {
      setView('login')
      alert('User Creted. Login')
    } else {
      const data = await res.json()
      alert(data.error)
    }
  }


  //post load
  const loadPost = async () => {
    const res = await fetch('/v1/posts', {
      headers: {Authorization: `Bearer ${token}`}
    })
    const posts = await res.json()
    setPosts(posts)
  }

  //create post
  const handleCreatePost = async (e) => {
    e.preventDefault()

    const formData = new FormData()
    formData.append('content', content)
    if (imageFile) formData.append('image', imageFile)

    const res = await fetch('/v1/posts', {
      method: 'POST',
      headers: { Authorization: `Bearer ${token}` },
      body: formData
    })
    if (res.ok) {
      setContent('')
      setImageFile(null)
      loadPost()
    } else {
      const data = await res.json()
      alert(data.error)
    }
  }

  //logout
  const handleLogout = () => {
    localStorage.removeItem('token')
    setToken('')
  }

  useEffect(() => {
    if (token && view === 'feed') {
      loadPost()
    }
  }, [token, view])

  // ===================== AUTH VIEW =====================
  if (!token) {
    return (
      <div className="min-vh-100 d-flex align-items-center justify-content-center" style={{ background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)' }}>
        <div className="container" style={{ maxWidth: '420px' }}>
          <div className="text-center mb-4">
            <h1 className="text-white fw-bold mb-1">Social API</h1>
            <p className="text-white-50 mb-0">Go backend demo</p>
          </div>

          {view === 'login' ? (
            <div className="card border-0 shadow-lg rounded-4">
              <div className="card-body p-4">
                <h2 className="text-center mb-4 fw-bold">Login</h2>
                <form onSubmit={handleLogin}>
                  <div className="mb-3">
                    <label className="form-label">Email</label>
                    <input
                      type="email"
                      className="form-control form-control-lg"
                      placeholder="you@example.com"
                      value={email}
                      onChange={(e) => setEmail(e.target.value)}
                    />
                  </div>
                  <div className="mb-4">
                    <label className="form-label">Password</label>
                    <input
                      type="password"
                      className="form-control form-control-lg"
                      placeholder="••••••••"
                      value={password}
                      onChange={(e) => setPassword(e.target.value)}
                    />
                  </div>
                  <button type="submit" className="btn btn-primary btn-lg w-100 fw-semibold">
                    Entrar
                  </button>
                </form>
                <p className="text-center mt-3 mb-0 text-secondary">
                  No tenés cuenta?{' '}
                  <button className="btn btn-link p-0 fw-semibold" onClick={() => setView('register')}>
                    Registrate
                  </button>
                </p>
              </div>
            </div>
          ) : (
            <div className="card border-0 shadow-lg rounded-4">
              <div className="card-body p-4">
                <h2 className="text-center mb-4 fw-bold">Registro</h2>
                <form onSubmit={handleRegister}>
                  <div className="mb-3">
                    <label className="form-label">Username</label>
                    <input
                      type="text"
                      className="form-control form-control-lg"
                      placeholder="sam"
                      value={username}
                      onChange={(e) => setUsername(e.target.value)}
                    />
                  </div>
                  <div className="mb-3">
                    <label className="form-label">Email</label>
                    <input
                      type="email"
                      className="form-control form-control-lg"
                      placeholder="you@example.com"
                      value={email}
                      onChange={(e) => setEmail(e.target.value)}
                    />
                  </div>
                  <div className="mb-4">
                    <label className="form-label">Password</label>
                    <input
                      type="password"
                      className="form-control form-control-lg"
                      placeholder="••••••••"
                      value={password}
                      onChange={(e) => setPassword(e.target.value)}
                    />
                  </div>
                  <button type="submit" className="btn btn-success btn-lg w-100 fw-semibold">
                    Crear usuario
                  </button>
                </form>
                <p className="text-center mt-3 mb-0 text-secondary">
                  Ya tenés cuenta?{' '}
                  <button className="btn btn-link p-0 fw-semibold" onClick={() => setView('login')}>
                    Login
                  </button>
                </p>
              </div>
            </div>
          )}
        </div>
      </div>
    )
  }

  // ===================== FEED VIEW =====================
  return (
    <div className="min-vh-100 bg-light">
      {/* Navbar */}
      <nav className="navbar navbar-dark bg-primary shadow-sm">
        <div className="container">
          <span className="navbar-brand fw-bold">🦫 Social API</span>
          <div className="d-flex align-items-center gap-2">
            <span className="badge bg-light text-primary">{email}</span>
            <button className="btn btn-outline-light btn-sm" onClick={handleLogout}>
              Logout
            </button>
          </div>
        </div>
      </nav>

      {/* Create post */}
      <div className="container py-4" style={{ maxWidth: '600px' }}>
        <div className="card border-0 shadow-sm rounded-4 mb-4">
          <div className="card-body">
            <form onSubmit={handleCreatePost}>
              <div className="d-flex gap-2 mb-2">
                <input
                  type="text"
                  className="form-control"
                  placeholder="Qué estás pensando?"
                  value={content}
                  onChange={(e) => setContent(e.target.value)}
                />
                <button type="submit" className="btn btn-primary px-4">Post</button>
              </div>
              <div className="d-flex align-items-center gap-2">
                <label
                  htmlFor="image-upload"
                  className="btn btn-outline-secondary btn-sm d-inline-flex align-items-center"
                  role="button"
                  title="Subir imagen"
                >
                  {/* Traditional paperclip icon */}
                  <svg
                    xmlns="http://www.w3.org/2000/svg"
                    width="16"
                    height="16"
                    viewBox="0 0 24 24"
                    fill="none"
                    stroke="currentColor"
                    strokeWidth="2"
                    strokeLinecap="round"
                    strokeLinejoin="round"
                  >
                    <path d="m21.44 11.05-9.19 9.19a6 6 0 0 1-8.49-8.49l8.57-8.57A4 4 0 1 1 18 8.84l-8.59 8.57a2 2 0 0 1-2.83-2.83l8.49-8.48" />
                  </svg>
                </label>
                <input
                  id="image-upload"
                  type="file"
                  accept="image/jpeg,image/png,image/webp"
                  className="d-none"
                  onChange={(e) => setImageFile(e.target.files[0] || null)}
                />
                {imageFile ? (
                  <small className="text-muted text-truncate">{imageFile.name}</small>
                ) : (
                  <small className="text-muted">Agregar imagen</small>
                )}
              </div>
            </form>
          </div>
        </div>

        {/* Feed */}
        <div className="d-flex flex-column gap-3">
          {posts.map((post) => (
            <div key={post.id} className="card border-0 shadow-sm rounded-4">
              <div className="card-body">
                <div className="d-flex align-items-center mb-2">
                  <div className="rounded-circle bg-primary text-white d-flex align-items-center justify-content-center me-2" style={{ width: 36, height: 36 }}>
                    #{post.user_id}
                  </div>
                  <div>
                    <div className="fw-semibold">user #{post.user_id}</div>
                    <small className="text-muted">{new Date(post.created_at).toLocaleString()}</small>
                  </div>
                </div>
                <p className="mb-0 fs-5">{post.content}</p>
                {post.image_url && (
                  <img
                    src={post.image_url}
                    className="img-fluid rounded mt-2"
                    alt="post"
                  />
                )}
              </div>
            </div>
          ))}
          {posts.length === 0 && (
            <div className="text-center py-5">
              <p className="text-muted fs-5">No hay posts todavía.</p>
              <p className="text-muted">Creá el primero!</p>
            </div>
          )}
        </div>
      </div>
    </div>
  )
}

export default App
