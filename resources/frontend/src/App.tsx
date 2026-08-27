import { BrowserRouter, Routes, Route } from "react-router-dom"
import DashboardLayout from "./pages/Dashboard"
import Home from "./pages/Home"
import Login from "./pages/Login"
import GitHubCallback from "./pages/GitHubCallback"
import User from "./pages/User"
import Agent from "./pages/Agent"
import AgentDetail from "./pages/AgentDetail"

function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/login" element={<Login />} />
        <Route path="/login/github" element={<GitHubCallback />} />
        <Route element={<DashboardLayout />}>
          <Route path="/" element={<Home />} />
          <Route path="/user" element={<User />} />
          <Route path="/agent" element={<Agent />} />
          <Route path="/agent/:id" element={<AgentDetail />} />
        </Route>
      </Routes>
    </BrowserRouter>
  )
}

export default App
