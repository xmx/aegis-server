import { useState } from "react"
import {
  useToastController,
  Toast,
  ToastTitle,
  ToastBody,
  Toaster,
  Button,
  type ToastIntent,
} from "@fluentui/react-components"
import { CopyRegular, CheckmarkRegular } from "@fluentui/react-icons"
import { ApiRequestError } from "@/lib/api"

function CopyButton({ text }: { text: string }) {
  const [copied, setCopied] = useState(false)

  const handleCopy = async () => {
    await navigator.clipboard.writeText(text)
    setCopied(true)
    setTimeout(() => setCopied(false), 1500)
  }

  return (
    <Button
      appearance="subtle"
      size="small"
      icon={copied ? <CheckmarkRegular /> : <CopyRegular />}
      onClick={handleCopy}
    />
  )
}

let dispatchToast: ((intent: ToastIntent, title: string, body: string) => void) | null = null

function FluentToaster() {
  const { dispatchToast: dispatch } = useToastController()

  dispatchToast = (intent: ToastIntent, title: string, body: string) => {
    dispatch(
      <Toast>
        <ToastTitle>{title}</ToastTitle>
        {body ? (
          <ToastBody>
            <div style={{ display: "flex", alignItems: "flex-start", gap: "8px" }}>
              <span style={{ flex: 1 }}>{body}</span>
              <CopyButton text={body} />
            </div>
          </ToastBody>
        ) : null}
      </Toast>,
      { intent }
    )
  }

  return <Toaster />
}

function showError(err: unknown) {
  if (!dispatchToast) return

  if (err instanceof ApiRequestError) {
    dispatchToast("error", err.message, err.detail)
    return
  }

  if (err instanceof Error) {
    dispatchToast("error", "错误", err.message)
    return
  }

  dispatchToast("error", "错误", "未知错误")
}

export { FluentToaster, showError }