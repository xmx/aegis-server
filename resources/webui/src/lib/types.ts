/* ===== API 数据类型 ===== */

export interface Semver {
  version: string
  number: string
}

export interface AgentProcessInfo {
  inet?: string
  semver?: Semver
  goos?: string
  goarch?: string
  pid?: number
  args?: string[]
  hostname?: string
  workdir?: string
  executable?: string
}

export interface AgentTunnelInfo {
  connected_at?: string
  keepalive_at?: string
  library_name?: string
  library_module?: string
  server_addr?: string
  remote_addr?: string
  receive_bytes?: number
  transmit_bytes?: number
}

export interface AgentRecord {
  id: string
  machine_id: string
  status: boolean
  enabled: boolean
  process_info: AgentProcessInfo
  tunnel_info: AgentTunnelInfo
  created_at: string
  updated_at?: string
}

export interface AgentConnRecord {
  id: string
  agent_id: string
  machine_id: string
  process_info: AgentProcessInfo
  tunnel_info: AgentTunnelInfo
  active_seconds?: number
  disconnected_at?: string
}

export interface PageResponse<T> {
  page: number
  size: number
  total: number
  records: T[]
}

export interface UserRecord {
  id: string
  enabled: boolean
  provider: string
  puid: string
  login: string
  name?: string
  avatar_url?: string
  company?: string
  email?: string
  location?: string
  created_at: string
  updated_at?: string
}

export interface OAuthProviderInfo {
  provider: string
  client_id: string
  redirect_uri: string
  scopes?: string[]
  auth_url: string
}

export interface ProblemDetails {
  status: number
  title: string
  detail: string
  instance: string
  method: string
  host: string
}

export interface BuildInfo {
  goos: string
  goarch: string
  version: string
  revision: string
  username: string
  workdir: string
  module: string
  committed_at?: string
  compiled_at?: string
  build_info?: {
    GoVersion: string
    Path: string
    Main: { Path: string; Version: string }
    Deps: { Path: string; Version: string }[]
    Settings: { Key: string; Value: string }[]
  }
}