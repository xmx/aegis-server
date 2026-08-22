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
import { GitHubIcon } from "@/components/Icons"

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