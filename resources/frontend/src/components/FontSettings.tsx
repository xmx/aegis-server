import { useEffect } from "react"
import {
  Button,
  Dialog,
  DialogTrigger,
  DialogSurface,
  DialogTitle,
  DialogBody,
  DialogActions,
  Select,
  makeStyles,
  tokens,
  Text,
} from "@fluentui/react-components"
import { SettingsRegular } from "@fluentui/react-icons"

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
  consolas: "Consolas, Menlo, monospace",
  menlo: "Menlo, Consolas, monospace",
  "courier-new": '"Courier New", monospace',
}

export function getFontStyle(cnFont: string, enFont: string): string {
  const en = EN_FONT_MAP[enFont] ?? EN_FONT_MAP["space-mono"]
  const cn = CN_FONT_MAP[cnFont]
  return cn ? `${en},${cn}` : en
}

const useStyles = makeStyles({
  preview: {
    borderRadius: tokens.borderRadiusMedium,
    border: `1px solid ${tokens.colorNeutralStroke1}`,
    backgroundColor: tokens.colorNeutralBackground2,
    padding: tokens.spacingHorizontalM,
  },
})

interface FontSettingsProps {
  cnFont: string
  enFont: string
  onCnFontChange: (value: string) => void
  onEnFontChange: (value: string) => void
}

function FontSettings({ cnFont, enFont, onCnFontChange, onEnFontChange }: FontSettingsProps) {
  const styles = useStyles()
  const fontStyle = getFontStyle(cnFont, enFont)

  // 按需加载非默认英文字体
  useEffect(() => {
    const defaultEnFonts = ["space-mono"]
    if (defaultEnFonts.includes(enFont)) return

    if (!document.getElementById("font-optional-styles")) {
      const link = document.createElement("link")
      link.id = "font-optional-styles"
      link.rel = "stylesheet"
      link.href = "https://fonts.loli.net/css2?family=JetBrains+Mono:ital,wght@0,400;0,700;1,400;1,700&family=Fira+Code:wght@400;700&family=Cascadia+Code:wght@400;700&family=Source+Code+Pro:ital,wght@0,400;0,700;1,400;1,700&family=IBM+Plex+Mono:ital,wght@0,400;0,700;1,400;1,700&display=swap"
      document.head.appendChild(link)
    }
  }, [enFont])

  return (
    <Dialog>
      <DialogTrigger disableButtonEnhancement>
        <Button appearance="outline" icon={<SettingsRegular />} />
      </DialogTrigger>
      <DialogSurface>
        <DialogBody>
          <DialogTitle>字体设置</DialogTitle>
          <div style={{ "--preview-font": fontStyle } as React.CSSProperties} className={styles.preview}>
            <Text size={200}>
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
              <Text weight="bold">Bold:</Text>
              <br />
              <Text weight="bold">abcdefghijklmnopqrstuvwxyz</Text>
              <br />
              <Text weight="bold">ABCDEFGHIJKLMNOPQRSTUVWXYZ</Text>
              <br />
              <Text weight="bold">0123456789 (){}[]</Text>
              <br />
              <Text weight="bold">+ - * / = .,;:!? #&$%@|^</Text>
              <br />
              <br />
              {"<!-- != := === >= >- >=> |-> -> <$>"}
              <br />
              {"</> #[ |||> |= ~@"}
              <br />
              <br />
              敏捷的棕色狐狸跨过懒狗
              <br />
              你好世界，欢迎使用 Aegis 系统
            </Text>
          </div>

          <div style={{ display: "flex", flexDirection: "column", gap: tokens.spacingVerticalM }}>
            <div>
              <Text size={200} weight="semibold">中文字体</Text>
              <Select value={cnFont} onChange={(_, d) => onCnFontChange(d.value)} style={{ width: "100%", marginTop: tokens.spacingVerticalXS }}>
                {CN_FONTS.map((f) => (
                  <option key={f.value} value={f.value} style={{ fontFamily: CN_FONT_MAP[f.value] }}>
                    {f.label}
                  </option>
                ))}
              </Select>
            </div>
            <div>
              <Text size={200} weight="semibold">英文字体</Text>
              <Select value={enFont} onChange={(_, d) => onEnFontChange(d.value)} style={{ width: "100%", marginTop: tokens.spacingVerticalXS }}>
                {EN_FONTS.map((f) => (
                  <option key={f.value} value={f.value} style={{ fontFamily: EN_FONT_MAP[f.value] }}>
                    {f.label}
                  </option>
                ))}
              </Select>
            </div>
          </div>
        </DialogBody>
        <DialogActions>
          <DialogTrigger disableButtonEnhancement>
            <Button appearance="secondary">关闭</Button>
          </DialogTrigger>
        </DialogActions>
      </DialogSurface>
    </Dialog>
  )
}

export { FontSettings, CN_FONTS, EN_FONTS }