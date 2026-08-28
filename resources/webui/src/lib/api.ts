import type { ProblemDetails } from '@/lib/types'

export class ApiRequestError extends Error {
  status: number
  detail: string
  instance: string
  method: string
  host: string

  constructor(err: ProblemDetails) {
    super(err.title)
    this.name = 'ApiRequestError'
    this.status = err.status
    this.detail = err.detail
    this.instance = err.instance
    this.method = err.method
    this.host = err.host
  }
}

export async function api<T>(url: string, options?: RequestInit): Promise<T> {
  const res = await fetch(url, options)

  if (!res.ok) {
    let err: ProblemDetails
    try {
      err = await res.json()
    } catch {
      throw new Error(`请求失败: ${res.status} ${res.statusText}`)
    }

    if (err.status === 401) {
      // 动态导入避免循环依赖
      const { useAuthStore } = await import('@/stores/auth')
      useAuthStore.getState().logout()
      throw new ApiRequestError(err)
    }

    throw new ApiRequestError(err)
  }

  return res.json()
}