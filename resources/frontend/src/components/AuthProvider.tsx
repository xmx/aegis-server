import { createContext, useContext, useState } from "react"

interface User {
  login: string
  name?: string
  avatar_url: string
}

interface AuthState {
  user: User | null
  setUser: (user: User | null) => void
}

const AuthContext = createContext<AuthState>({
  user: null,
  setUser: () => {},
})

function AuthProvider({ children }: { children: React.ReactNode }) {
  const [user, setUser] = useState<User | null>(null)

  return (
    <AuthContext.Provider value={{ user, setUser }}>
      {children}
    </AuthContext.Provider>
  )
}

function useAuth() {
  return useContext(AuthContext)
}

export { AuthProvider, useAuth }
export type { User }