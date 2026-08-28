import { create } from 'zustand'
import { persist } from 'zustand/middleware'

export interface UserInfo {
  login: string
  name?: string
  avatar_url: string
  provider?: string
  email?: string
}

interface AuthStore {
  user: UserInfo | null
  setUser: (user: UserInfo | null) => void
  logout: () => void
}

export const useAuthStore = create<AuthStore>()(
  persist(
    (set) => ({
      user: null,
      setUser: (user) => set({ user }),
      logout: () => set({ user: null }),
    }),
    { name: 'aegis.auth' },
  ),
)