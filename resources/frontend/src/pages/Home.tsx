import { Text, makeStyles, tokens } from "@fluentui/react-components"
import { useDocumentTitle } from "@/hooks/useDocumentTitle"

const useStyles = makeStyles({
  root: {
    flex: 1,
    display: "flex",
    alignItems: "center",
    justifyContent: "center",
    padding: tokens.spacingHorizontalXXL,
  },
})

function Home() {
  const styles = useStyles()
  useDocumentTitle("首页")

  return (
    <div className={styles.root}>
      <Text style={{ color: tokens.colorNeutralForeground3 }}>欢迎使用 Aegis</Text>
    </div>
  )
}

export default Home