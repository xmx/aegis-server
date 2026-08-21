import { useState } from "react"
import { Button } from "@/components/ui/button"
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
} from "@/components/ui/card"
import { api } from "@/lib/api"
import { showError } from "@/lib/toast"

function GitHubIcon({ className }: { className?: string }) {
  return (
    <svg
      className={className}
      viewBox="0 0 24 24"
      fill="currentColor"
      xmlns="http://www.w3.org/2000/svg"
    >
      <path d="M12 0C5.374 0 0 5.373 0 12c0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23A11.509 11.509 0 0112 5.803c1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576C20.566 21.797 24 17.3 24 12c0-6.627-5.373-12-12-12z" />
    </svg>
  )
}

interface ProviderResponse {
  provider: string
  client_id: string
  redirect_uri: string
  auth_url: string
  scopes?: string[]
}

function Login() {
  const [loading, setLoading] = useState(false)

  const handleGitHubLogin = async () => {
    setLoading(true)
    try {
      const params = new URLSearchParams({ provider: "github" })
      const data = await api<ProviderResponse>(`/api/oauth/provider?${params}`)

      const authParams = new URLSearchParams({
        client_id: data.client_id,
        redirect_uri: data.redirect_uri,
        scope: (data.scopes ?? []).join(" "),
      })

      window.location.href = `${data.auth_url}?${authParams.toString()}`
    } catch (err) {
      showError(err)
      setLoading(false)
    }
  }

  return (
    <div className="flex min-h-screen items-center justify-center bg-background p-4">
      <Card className="w-[300px]">
        <CardHeader className="text-center">
          <CardTitle className="text-base font-medium">选择登录方式</CardTitle>
        </CardHeader>
        <CardContent className="flex flex-col gap-4">
          <Button
            variant="outline"
            size="lg"
            className="w-full"
            disabled={loading}
            onClick={handleGitHubLogin}
          >
            <GitHubIcon className="size-5" />
            {loading ? "请求中..." : "GitHub"}
          </Button>
        </CardContent>
      </Card>
    </div>
  )
}

export default Login