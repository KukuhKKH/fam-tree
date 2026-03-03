<script lang="ts">
  import { page } from "$app/state";
  import { env } from "$env/dynamic/public";
  import { Button } from "$lib/components/ui/button/index.js";
  import ThemeToggle from "$lib/components/ThemeToggle.svelte";
  import {
    LayoutDashboard,
    Users,
    TreeDeciduous,
    LogOut,
    Menu,
  } from "lucide-svelte";

  let { children } = $props();
  let isSidebarOpen = $state(true);

  const navItems = [
    { name: "Dashboard", href: "/dashboard", icon: LayoutDashboard },
    { name: "Keluarga Saya", href: "/families", icon: Users },
    { name: "Silsilah", href: "/tree", icon: TreeDeciduous },
  ];

  async function logout() {
    try {
      const response = await fetch(
        `${env.PUBLIC_BACKEND_URL || "http://localhost:8080"}/auth/logout`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          credentials: "include",
        },
      );

      if (response.ok) {
        const { data } = await response.json();
        if (data && data.logout_url) {
          window.location.href = data.logout_url;
        } else {
          window.location.href = "/login";
        }
      }
    } catch (error) {
      console.error("Logout error:", error);
      window.location.href = "/login";
    }
  }
</script>

<div class="flex h-screen bg-zinc-50 dark:bg-zinc-950">
  <!-- Sidebar -->
  <aside
    class="border-r border-zinc-200 dark:border-zinc-800 bg-white dark:bg-zinc-900 transition-all duration-300 ease-in-out flex flex-col
		{isSidebarOpen ? 'w-64' : 'w-20'}"
  >
    <div
      class="h-16 flex items-center px-6 border-b border-zinc-200 dark:border-zinc-800"
    >
      <span
        class="font-bold text-xl tracking-tight text-primary {isSidebarOpen
          ? 'block'
          : 'hidden'}">Silsilah</span
      >
      <span
        class="font-bold text-xl tracking-tight text-primary {!isSidebarOpen
          ? 'block'
          : 'hidden'}">S</span
      >
    </div>

    <nav class="flex-1 p-4 space-y-2">
      {#each navItems as item}
        <a
          href={item.href}
          class="flex items-center gap-4 px-3 py-2 rounded-md transition-colors
					{page.url.pathname === item.href
            ? 'bg-primary/10 text-primary font-medium'
            : 'text-zinc-600 dark:text-zinc-400 hover:bg-zinc-100 dark:hover:bg-zinc-800'}"
        >
          <item.icon size={22} />
          <span class={isSidebarOpen ? "block" : "hidden"}>{item.name}</span>
        </a>
      {/each}
    </nav>

    <div class="p-4 border-t border-zinc-200 dark:border-zinc-800">
      <Button
        variant="ghost"
        class="w-full justify-start gap-4 text-zinc-600 dark:text-zinc-400"
        onclick={logout}
      >
        <LogOut size={22} />
        <span class={isSidebarOpen ? "block" : "hidden"}>Logout</span>
      </Button>
    </div>
  </aside>

  <!-- Main Content -->
  <main class="flex-1 flex flex-col overflow-hidden">
    <!-- Header -->
    <header
      class="h-16 border-b border-zinc-200 dark:border-zinc-800 bg-white dark:bg-zinc-900 flex items-center justify-between px-6"
    >
      <Button
        variant="ghost"
        size="icon"
        onclick={() => (isSidebarOpen = !isSidebarOpen)}
      >
        <Menu size={20} />
      </Button>

      <div class="flex items-center gap-4">
        <ThemeToggle />
        <span class="text-sm font-medium text-zinc-700 dark:text-zinc-300">
          {page.data.user?.name || "User"}
        </span>
        <div
          class="h-8 w-8 rounded-full bg-primary/20 flex items-center justify-center text-primary font-bold"
        >
          {(page.data.user?.name || "U")[0].toUpperCase()}
        </div>
      </div>
    </header>

    <!-- Page Content -->
    <div class="flex-1 overflow-y-auto p-8">
      {@render children()}
    </div>
  </main>
</div>
