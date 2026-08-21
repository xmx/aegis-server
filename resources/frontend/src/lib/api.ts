interface ApiError {
  status: number
  title: string
  detail: string
  instance: string
  method: string
  host: string
}

class ApiRequestError extends Error {
  status: number
  detail: string
  instance: string
  method: string
  host: string

  constructor(err: ApiError) {
    super(err.title)
    this.name = "ApiRequestError"
    this.status = err.status
    this.detail = err.detail
    this.instance = err.instance
    this.method = err.method
    this.host = err.host
  }
}

async function api<T>(url: string, options?: RequestInit): Promise<T> {
  const res = await fetch(url, options)

  if (!res.ok) {
    let err: ApiError
    try {
      err = await res.json()
    } catch {
      throw new Error(`请求失败: ${res.status} ${res.statusText}`)
    }

    if (err.status === 401) {
      window.location.href = "/login"
      throw new ApiRequestError(err)
    }

    throw new ApiRequestError(err)
  }

  return res.json()
}

export { api, ApiRequestError }
export type { ApiError }