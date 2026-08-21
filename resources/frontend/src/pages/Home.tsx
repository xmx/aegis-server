import { Button } from "@/components/ui/button"

function Home() {
  return (
    <div className="flex flex-1 flex-col items-center justify-center gap-8 p-8">
      <div className="flex flex-col items-center gap-4">
        <h1 className="text-3xl font-bold tracking-tight">Aegis</h1>
        <p className="text-muted-foreground">
          React + shadcn/ui + Tailwind CSS v4
        </p>
      </div>

      <div className="flex flex-wrap items-center gap-3">
        <Button>默认</Button>
        <Button variant="secondary">次要</Button>
        <Button variant="outline">轮廓</Button>
        <Button variant="ghost">幽灵</Button>
        <Button variant="destructive">危险</Button>
        <Button variant="link">链接</Button>
      </div>
    </div>
  )
}

export default Home