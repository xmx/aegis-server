import { useEffect, useRef } from "react"
import { useSearchParams, useNavigate } from "react-router-dom"
import { api } from "@/lib/api"
import { showError } from "@/lib/toast"
import { useAuth } from "@/components/AuthProvider"
import type { User } from "@/components/AuthProvider"
import { Spinner, makeStyles } from "@fluentui/react-components"
import { useDocumentTitle } from "@/hooks/useDocumentTitle"

const useStyles = makeStyles({
  root: {
    display: "flex",
    minHeight: "100vh",
    alignItems: "center",
    justifyContent: "center",
  },
})

function GitHubCallback() {
  const styles = useStyles()
  useDocumentTitle("正在登录")
  const [searchParams] = useSearchParams()
  const navigate = useNavigate()
  const { setUser } = useAuth()
  const code = searchParams.get("code")
  const called = useRef(false)

  useEffect(() => {
    if (!code || called.current) return
    called.current = true

    api<User>("/api/oauth/github", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ code, redirect_uri: window.location.origin + window.location.pathname }),
    })
      .then((user) => {
        setUser(user)
        navigate("/", { replace: true })
      })
      .catch((err) => {
        showError(err)
        navigate("/login", { replace: true })
      })
  }, [])

  return (
    <div className={styles.root}>
      <Spinner label="正在登录..." />
    </div>
  )
}

export default GitHubCallback