import { ReactNode } from "react";
import {
  Sidebar,
  SidebarContent,
  SidebarHeader,
  SidebarMenu,
  SidebarMenuItem,
  SidebarMenuButton,
  SidebarProvider,
  SidebarTrigger,
} from "@/components/ui/sidebar";
import { Button } from "@/components/ui/button";
import {
  Database,
  Activity,
  Settings,
  LogOut,
  RefreshCw,
  Download,
  Server,
  MessageSquare,
} from "lucide-react";
import { useNavigate, useLocation } from "react-router-dom";

interface LayoutProps {
  children: ReactNode;
}

export function Layout({ children }: LayoutProps) {
  const navigate = useNavigate();
  const location = useLocation();

  const menuItems = [
    { icon: Activity, label: "Dashboard", path: "/dashboard" },
    { icon: Database, label: "Redis Monitor", path: "/redis" },
    { icon: Database, label: "Database Monitor", path: "/database" },
    { icon: MessageSquare, label: "Kafka Monitor", path: "/kafka" },
    { icon: Settings, label: "Configuration", path: "/config" },
  ];

  const handleLogout = () => {
    // TODO: Implement logout logic
    navigate("/login");
  };

  const handleRefresh = () => {
    window.location.reload();
  };

  const handleExport = () => {
    // TODO: Implement export functionality
    console.log("Export data");
  };

  return (
    <div className="min-h-screen bg-background dark">
      <SidebarProvider>
        <Sidebar className="border-r border-border">
          <SidebarHeader className="p-4">
            <div className="flex items-center gap-2">
              <Server className="h-6 w-6 text-primary" />
              <span className="font-semibold text-lg">
                Redis & Kafka Monitor
              </span>
            </div>
          </SidebarHeader>
          <SidebarContent>
            <SidebarMenu>
              {menuItems.map((item) => (
                <SidebarMenuItem key={item.path}>
                  <SidebarMenuButton
                    onClick={() => navigate(item.path)}
                    isActive={location.pathname === item.path}
                    className="w-full justify-start"
                  >
                    <item.icon className="h-4 w-4" />
                    <span>{item.label}</span>
                  </SidebarMenuButton>
                </SidebarMenuItem>
              ))}
            </SidebarMenu>
          </SidebarContent>
          <div className="p-4 border-t border-border">
            <Button
              variant="ghost"
              onClick={handleLogout}
              className="w-full justify-start text-muted-foreground hover:text-foreground"
            >
              <LogOut className="h-4 w-4 mr-2" />
              Logout
            </Button>
          </div>
        </Sidebar>

        <main className="flex-1">
          <header className="border-b border-border bg-card/50 backdrop-blur-sm">
            <div className="flex items-center justify-between px-6 py-4">
              <div className="flex items-center gap-4">
                <SidebarTrigger />
                <div className="text-sm text-muted-foreground">
                  Last updated: {new Date().toLocaleTimeString()}
                </div>
              </div>
              <div className="flex items-center gap-2">
                <Button
                  variant="outline"
                  size="sm"
                  onClick={handleRefresh}
                  className="gap-2"
                >
                  <RefreshCw className="h-4 w-4" />
                  Refresh
                </Button>
                <Button
                  variant="outline"
                  size="sm"
                  onClick={handleExport}
                  className="gap-2"
                >
                  <Download className="h-4 w-4" />
                  Export CSV
                </Button>
              </div>
            </div>
          </header>

          <div className="p-6">{children}</div>
        </main>
      </SidebarProvider>
    </div>
  );
}
