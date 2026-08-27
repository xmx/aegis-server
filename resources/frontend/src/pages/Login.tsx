import { useState } from "react"
import { Text, makeStyles, tokens, Spinner } from "@fluentui/react-components"
import { api } from "@/lib/api"
import { showError } from "@/lib/toast"
import { GitHubIcon, GoogleIcon, GitLabIcon } from "@/components/Icons"
import { useDocumentTitle } from "@/hooks/useDocumentTitle"

const useStyles = makeStyles({
  root: {
    display: "flex",
    minHeight: "100vh",
    flexDirection: "column",
    alignItems: "center",
    justifyContent: "center",
    backgroundColor: tokens.colorNeutralBackground2,
    gap: tokens.spacingVerticalXXL,
  },
  header: {
    display: "flex",
    flexDirection: "column",
    alignItems: "center",
    gap: tokens.spacingVerticalS,
  },
  providers: {
    display: "flex",
    gap: tokens.spacingHorizontalL,
    justifyContent: "center",
  },
  providerButton: {
    display: "flex",
    flexDirection: "column",
    alignItems: "center",
    gap: tokens.spacingVerticalS,
    cursor: "pointer",
    background: "none",
    border: "none",
    padding: 0,
    outline: "none",
  },
  providerIcon: {
    display: "flex",
    alignItems: "center",
    justifyContent: "center",
    width: "64px",
    height: "64px",
    borderRadius: "50%",
    backgroundColor: tokens.colorNeutralBackground1,
    border: `2px solid ${tokens.colorNeutralStroke1}`,
    transition: "border-color 0.15s, box-shadow 0.15s, transform 0.1s",
    color: tokens.colorNeutralForeground1,
    fontSize: "36px",
  },
})

interface ProviderResponse {
  provider: string
  client_id: string
  redirect_uri: string
  auth_url: string
  scopes?: string[]
}

function Login() {
  const styles = useStyles()
  const [loadingProvider, setLoadingProvider] = useState<string | null>(null)
  useDocumentTitle("登录")

  const doLogin = async (provider: string) => {
    const params = new URLSearchParams({ provider })
    const data = await api<ProviderResponse>(`/api/oauth/provider?${params}`)

    const authParams = new URLSearchParams({
      client_id: data.client_id,
      redirect_uri: data.redirect_uri,
      scope: (data.scopes ?? []).join(" "),
    })

    window.location.href = `${data.auth_url}?${authParams.toString()}`
  }

  const handleGitHubLogin = async () => {
    setLoadingProvider("github")
    try {
      await doLogin("github")
    } catch (err) {
      showError(err)
      setLoadingProvider(null)
    }
  }

  const handleGoogleLogin = async () => {
    setLoadingProvider("google")
    try {
      await doLogin("google")
    } catch (err) {
      showError(err)
      setLoadingProvider(null)
    }
  }

  const handleGitLabLogin = async () => {
    setLoadingProvider("gitlab")
    try {
      await doLogin("gitlab")
    } catch (err) {
      showError(err)
      setLoadingProvider(null)
    }
  }

  return (
    <div className={styles.root}>
      <div className={styles.header}>
        <Text size={400} style={{ color: tokens.colorNeutralForeground3 }}>
          请选择登录方式
        </Text>
      </div>

      <div className={styles.providers}>
        <button
          className={styles.providerButton}
          onClick={handleGitHubLogin}
          title="GitHub"
        >
<div className={styles.providerIcon}>
              {loadingProvider === "github" ? <Spinner size="small" /> : <GitHubIcon />}
            </div>
          <Text size={200}>GitHub</Text>
        </button>

        <button
          className={styles.providerButton}
          onClick={handleGoogleLogin}
          title="Google"
        >
<div className={styles.providerIcon}>
              {loadingProvider === "google" ? <Spinner size="small" /> : <GoogleIcon />}
            </div>
          <Text size={200}>Google</Text>
        </button>

        <button
          className={styles.providerButton}
          onClick={handleGitLabLogin}
          title="GitLab"
        >
<div className={styles.providerIcon}>
              {loadingProvider === "gitlab" ? <Spinner size="small" /> : <GitLabIcon />}
            </div>
          <Text size={200}>GitLab</Text>
        </button>
      </div>

      </div>
  )
}

export default Login