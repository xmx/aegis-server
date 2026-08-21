import { CheckIcon, CopyIcon } from "lucide-react"
import { toast } from "@/components/ui/toast"
import { ApiRequestError } from "@/lib/api"

function showError(err: unknown) {
  let title = "错误"
  let detail = ""

  if (err instanceof ApiRequestError) {
    title = err.message
    detail = err.detail
  } else if (err instanceof Error) {
    detail = err.message
  } else {
    detail = "未知错误"
  }

  let copied = false

  const handleCopy = (id: string) => {
    if (copied) return
    copied = true
    navigator.clipboard.writeText(detail)
    toast.update(id, {
      actionProps: {
        children: <CheckIcon className="size-4 text-emerald-500" />,
      },
    })
    setTimeout(() => {
      copied = false
      toast.update(id, {
        actionProps: {
          onClick: () => handleCopy(id),
          children: <CopyIcon className="size-4" />,
        },
      })
    }, 1500)
  }

  const id = toast.add({
    type: "error",
    title,
    description: detail,
    actionProps: {
      onClick: () => handleCopy(id),
      children: <CopyIcon className="size-4" />,
    },
  })
}

export { showError }