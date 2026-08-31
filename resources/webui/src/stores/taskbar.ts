import { create } from 'zustand'
import { persist } from 'zustand/middleware'

interface TaskbarStore {
  /** 从默认固定中取消固定的 appId */
  unpinned: string[]
  /** 动态固定到任务栏的 appId */
  extraPinned: string[]
  pin: (id: string) => void
  unpin: (id: string) => void
}

export const useTaskbarStore = create<TaskbarStore>()(
  persist(
    (set) => ({
      unpinned: [],
      extraPinned: [],
      pin: (id) =>
        set((s) => ({
          unpinned: s.unpinned.filter((x) => x !== id),
          extraPinned: [...s.extraPinned.filter((x) => x !== id), id],
        })),
      unpin: (id) =>
        set((s) => ({
          unpinned: [...s.unpinned.filter((x) => x !== id), id],
          extraPinned: s.extraPinned.filter((x) => x !== id),
        })),
    }),
    { name: 'aegis.taskbar' },
  ),
)