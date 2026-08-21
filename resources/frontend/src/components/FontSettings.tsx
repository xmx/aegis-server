import { Settings } from "lucide-react"
import { Button } from "@/components/ui/button"
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog"
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
} from "@/components/ui/select"

const CN_FONTS = [
  { value: "system", label: "系统默认" },
  { value: "microsoft-yahei", label: "微软雅黑" },
  { value: "simsun", label: "宋体" },
  { value: "simhei", label: "黑体" },
  { value: "fangsong", label: "仿宋" },
  { value: "kaiti", label: "楷体" },
  { value: "noto-sans-sc", label: "思源黑体" },
  { value: "noto-serif-sc", label: "思源宋体" },
  { value: "pingfang-sc", label: "苹方" },
]

const EN_FONTS = [
  { value: "space-mono", label: "Space Mono" },
  { value: "jetbrains-mono", label: "JetBrains Mono" },
  { value: "fira-code", label: "Fira Code" },
  { value: "cascadia-code", label: "Cascadia Code" },
  { value: "source-code-pro", label: "Source Code Pro" },
  { value: "ibm-plex-mono", label: "IBM Plex Mono" },
  { value: "consolas", label: "Consolas" },
  { value: "menlo", label: "Menlo" },
  { value: "courier-new", label: "Courier New" },
]

const CN_FONT_MAP: Record<string, string> = {
  system: "",
  "microsoft-yahei": '"Microsoft YaHei", "微软雅黑", sans-serif',
  simsun: 'SimSun, "宋体", serif',
  simhei: 'SimHei, "黑体", sans-serif',
  fangsong: 'FangSong, "仿宋", serif',
  kaiti: 'KaiTi, "楷体", serif',
  "noto-sans-sc": '"Noto Sans SC", sans-serif',
  "noto-serif-sc": '"Noto Serif SC", serif',
  "pingfang-sc": '"PingFang SC", "苹方", sans-serif',
}

const EN_FONT_MAP: Record<string, string> = {
  "space-mono": '"Space Mono", Menlo, Consolas, monospace',
  "jetbrains-mono": '"JetBrains Mono", Menlo, Consolas, monospace',
  "fira-code": '"Fira Code", Menlo, Consolas, monospace',
  "cascadia-code": '"Cascadia Code", Menlo, Consolas, monospace',
  "source-code-pro": '"Source Code Pro", Menlo, Consolas, monospace',
  "ibm-plex-mono": '"IBM Plex Mono", Menlo, Consolas, monospace',
  consolas: 'Consolas, Menlo, monospace',
  menlo: 'Menlo, Consolas, monospace',
  "courier-new": '"Courier New", monospace',
}

export function getFontStyle(cnFont: string, enFont: string): string {
  const en = EN_FONT_MAP[enFont] ?? EN_FONT_MAP["space-mono"]
  const cn = CN_FONT_MAP[cnFont]
  return cn ? `${en},${cn}` : en
}

function cnFontFace(key: string): string {
  return CN_FONT_MAP[key] ?? ""
}

function enFontFace(key: string): string {
  return EN_FONT_MAP[key] ?? ""
}

function cnLabel(key: string): string {
  return CN_FONTS.find((f) => f.value === key)?.label ?? key
}

function enLabel(key: string): string {
  return EN_FONTS.find((f) => f.value === key)?.label ?? key
}

interface FontSettingsProps {
  cnFont: string
  enFont: string
  onCnFontChange: (value: string) => void
  onEnFontChange: (value: string) => void
}

function FontSettings({
  cnFont,
  enFont,
  onCnFontChange,
  onEnFontChange,
}: FontSettingsProps) {
  return (
    <Dialog>
      <DialogTrigger
        render={
          <Button variant="outline" size="icon">
            <Settings className="size-4" />
          </Button>
        }
      />
      <DialogContent className="sm:max-w-md">
        <DialogHeader>
          <DialogTitle>字体设置</DialogTitle>
          <DialogDescription>选择页面显示的中文字体和英文字体</DialogDescription>
        </DialogHeader>
        <div className="flex flex-col gap-5 py-4">
          <div
            className="rounded-lg border bg-muted/50 p-4"
            style={{ fontFamily: getFontStyle(cnFont, enFont) }}
          >
            <p className="text-sm leading-relaxed">
                Default:
                <br />
                abcdefghijklmnopqrstuvwxyz
                <br />
                ABCDEFGHIJKLMNOPQRSTUVWXYZ
                <br />
                0123456789 (){}[]
                <br />
                + - * / = .,;:!? #&$%@|^
                <br />
                <br />
                <span className="font-bold">Bold:</span>
                <br />
                <span className="font-bold">abcdefghijklmnopqrstuvwxyz</span>
                <br />
                <span className="font-bold">ABCDEFGHIJKLMNOPQRSTUVWXYZ</span>
                <br />
                <span className="font-bold">0123456789 (){}[]</span>
                <br />
                <span className="font-bold">+ - * / = .,;:!? #&$%@|^</span>
                <br />
                <br />
                &lt;!-- != := === &gt;= &gt;- &gt;=&gt; |-&gt; -&gt; &lt;$&gt;
                <br />
                &lt;/&gt; #[ |||&gt; |= ~@
              </p>
          </div>
          <div className="flex flex-col gap-2">
            <label className="text-sm font-medium">中文字体</label>
            <Select value={cnFont} onValueChange={(v) => onCnFontChange(v ?? "system")}>
              <SelectTrigger>
                <span style={{ fontFamily: cnFontFace(cnFont) }}>{cnLabel(cnFont)}</span>
              </SelectTrigger>
              <SelectContent>
                {CN_FONTS.map((font) => (
                  <SelectItem key={font.value} value={font.value} style={{ fontFamily: cnFontFace(font.value) }}>
                    {font.label}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>
          <div className="flex flex-col gap-2">
            <label className="text-sm font-medium">英文字体</label>
            <Select value={enFont} onValueChange={(v) => onEnFontChange(v ?? "space-mono")}>
              <SelectTrigger>
                <span style={{ fontFamily: enFontFace(enFont) }}>{enLabel(enFont)}</span>
              </SelectTrigger>
              <SelectContent>
                {EN_FONTS.map((font) => (
                  <SelectItem key={font.value} value={font.value} style={{ fontFamily: enFontFace(font.value) }}>
                    {font.label}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>
        </div>
      </DialogContent>
    </Dialog>
  )
}

export { FontSettings, CN_FONTS, EN_FONTS }