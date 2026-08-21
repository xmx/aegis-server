import { useSearchParams } from "react-router-dom"
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card"

function GitHubCallback() {
  const [searchParams] = useSearchParams()
  const code = searchParams.get("code")
  const state = searchParams.get("state")

  if (!code) {
    return (
      <div className="flex min-h-screen items-center justify-center">
        <p className="text-muted-foreground">缺少认证参数</p>
      </div>
    )
  }

  return (
    <div className="flex min-h-screen items-center justify-center p-4">
      <Card className="w-full max-w-sm">
        <CardHeader>
          <CardTitle className="text-base font-medium">GitHub 认证回调</CardTitle>
          <CardDescription>后端正在开发中，以下是 GitHub 返回的临时授权码</CardDescription>
        </CardHeader>
        <CardContent className="flex flex-col gap-4">
          <div className="flex flex-col gap-1.5">
            <label className="text-sm font-medium text-muted-foreground">授权码</label>
            <code className="rounded-md bg-muted px-3 py-2 text-sm break-all font-mono">
              {code}
            </code>
          </div>
          {state && (
            <div className="flex flex-col gap-1.5">
              <label className="text-sm font-medium text-muted-foreground">Nonce</label>
              <code className="rounded-md bg-muted px-3 py-2 text-sm break-all font-mono">
                {state}
              </code>
            </div>
          )}
        </CardContent>
      </Card>
    </div>
  )
}

export default GitHubCallback