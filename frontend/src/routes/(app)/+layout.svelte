<script lang="ts">
  import { page } from "$app/state";
  import { env } from "$env/dynamic/public";
  import { Button } from "$lib/components/ui/button/index.js";
  import ThemeToggle from "$lib/components/ThemeToggle.svelte";
  import {
    LayoutDashboard,
    TreeDeciduous,
    LogOut,
    Menu,
    Home,
    Share2,
    ChevronRight,
    Globe,
    Lock,
    Link,
    UserCircle,
  } from "lucide-svelte";

  let { children } = $props();
  let isSidebarOpen = $state(true);

  // Detect family context
  const familySlug = $derived(() => {
    const match = page.url.pathname.match(/^\/families\/([^/]+)/);
    return match ? match[1] : null;
  });

  const isInFamilyContext = $derived(!!familySlug());

  // Family data from page.data (available when in family context)
  const familyName = $derived(
    page.data?.family?.name || page.data?.slug || familySlug() || "",
  );
  const familyVisibility = $derived(page.data?.family?.visibility || "private");

  const visibilityIcons = {
    public: Globe,
    private: Lock,
    link_only: Link,
  };

  const VisibilityIcon = $derived(
    visibilityIcons[familyVisibility as keyof typeof visibilityIcons] || Lock,
  );

  // Global nav items
  const globalNavItems = [
    { name: "Dashboard", href: "/dashboard", icon: LayoutDashboard },
  ];

  // Family-context nav items (dynamic based on slug)
  const familyNavItems = $derived(
    familySlug()
      ? [
          {
            name: "Ikhtisar",
            href: `/families/${familySlug()}`,
            icon: Home,
            exact: true,
          },
          {
            name: "Pohon Keluarga",
            href: `/families/${familySlug()}/tree`,
            icon: TreeDeciduous,
            exact: false,
          },
        ]
      : [],
  );

  function isActive(href: string, exact = false): boolean {
    if (exact) return page.url.pathname === href;
    return page.url.pathname.startsWith(href);
  }

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
		{isSidebarOpen ? 'w-72' : 'w-20'}"
  >
    <div
      class="h-16 flex items-center px-6 border-b border-zinc-200 dark:border-zinc-800"
    >
      <a
        href="/dashboard"
        class="flex items-center gap-2 hover:opacity-80 transition-opacity"
      >
        <TreeDeciduous size={22} class="text-primary flex-shrink-0" />
        <span
          class="font-bold text-xl tracking-tight text-primary {isSidebarOpen
            ? 'block'
            : 'hidden'}"
        >
          Silsilah
        </span>
      </a>
    </div>

    <nav class="flex-1 p-4 space-y-1 overflow-y-auto">
      <!-- Global Navigation -->
      {#each globalNavItems as item}
        <a
          href={item.href}
          class="flex items-center gap-3 px-3 py-2.5 rounded-xl transition-all duration-200
					{isActive(item.href)
            ? 'bg-primary/10 text-primary font-bold shadow-sm'
            : 'text-zinc-500 dark:text-zinc-400 hover:bg-zinc-100 dark:hover:bg-zinc-800 hover:text-zinc-900 dark:hover:text-zinc-100'}"
        >
          <item.icon size={20} class="flex-shrink-0" />
          <span class="{isSidebarOpen ? 'block' : 'hidden'} text-sm"
            >{item.name}</span
          >
        </a>
      {/each}

      <!-- Family Context Section -->
      {#if isInFamilyContext && familySlug()}
        <div class="pt-4 mt-3">
          <!-- Family Context Header -->
          {#if isSidebarOpen}
            <div class="px-3 mb-3">
              <div
                class="p-3 rounded-xl bg-gradient-to-br from-primary/5 via-primary/3 to-transparent border border-primary/10"
              >
                <div class="flex items-center gap-2 mb-1">
                  <div
                    class="w-6 h-6 rounded-lg bg-primary/15 flex items-center justify-center"
                  >
                    <UserCircle size={14} class="text-primary" />
                  </div>
                  <span
                    class="text-[10px] font-bold uppercase tracking-widest text-zinc-400"
                  >
                    Keluarga Aktif
                  </span>
                </div>
                <p
                  class="text-sm font-bold text-zinc-900 dark:text-white truncate leading-tight"
                >
                  {familyName}
                </p>
                <div
                  class="flex items-center gap-1 mt-1.5 text-[10px] text-zinc-400"
                >
                  <VisibilityIcon size={9} />
                  <span class="capitalize">
                    {familyVisibility === "link_only"
                      ? "Link Only"
                      : familyVisibility}
                  </span>
                </div>
              </div>
            </div>
          {:else}
            <div class="px-3 mb-3">
              <div
                class="w-full aspect-square rounded-xl bg-primary/10 flex items-center justify-center"
                title={familyName}
              >
                <span class="text-primary font-black text-sm">
                  {familyName[0]?.toUpperCase() || "F"}
                </span>
              </div>
            </div>
          {/if}

          <!-- Divider with label -->
          {#if isSidebarOpen}
            <div class="px-3 mb-2">
              <span
                class="text-[10px] font-bold uppercase tracking-widest text-zinc-400"
              >
                Menu Keluarga
              </span>
            </div>
          {:else}
            <div
              class="mx-3 mb-2 border-t border-zinc-200 dark:border-zinc-700"
            ></div>
          {/if}

          <!-- Family nav items -->
          {#each familyNavItems as item}
            <a
              href={item.href}
              class="flex items-center gap-3 px-3 py-2.5 rounded-xl transition-all duration-200
							{isActive(item.href, item.exact)
                ? 'bg-primary/10 text-primary font-bold shadow-sm'
                : 'text-zinc-500 dark:text-zinc-400 hover:bg-zinc-100 dark:hover:bg-zinc-800 hover:text-zinc-900 dark:hover:text-zinc-100'}"
            >
              <item.icon size={20} class="flex-shrink-0" />
              <span class="{isSidebarOpen ? 'block' : 'hidden'} text-sm"
                >{item.name}</span
              >
            </a>
          {/each}

          <!-- Share link (public families only) -->
          {#if familyVisibility !== "private"}
            <a
              href={`/s/${familySlug()}`}
              target="_blank"
              class="flex items-center gap-3 px-3 py-2.5 rounded-xl transition-all duration-200 text-zinc-500 dark:text-zinc-400 hover:bg-emerald-50 dark:hover:bg-emerald-900/20 hover:text-emerald-600 dark:hover:text-emerald-400 mt-1"
            >
              <Share2 size={20} class="flex-shrink-0" />
              <span class="{isSidebarOpen ? 'block' : 'hidden'} text-sm"
                >Lihat Halaman Publik</span
              >
            </a>
          {/if}
        </div>
      {/if}
    </nav>

    <div class="p-4 border-t border-zinc-200 dark:border-zinc-800">
      <Button
        variant="ghost"
        class="w-full justify-start gap-3 text-zinc-500 dark:text-zinc-400 hover:text-red-500 rounded-xl"
        onclick={logout}
      >
        <LogOut size={20} />
        <span class="{isSidebarOpen ? 'block' : 'hidden'} text-sm">Logout</span>
      </Button>
    </div>
  </aside>

  <!-- Main Content -->
  <main class="flex-1 flex flex-col overflow-hidden">
    <!-- Header -->
    <header
      class="h-16 border-b border-zinc-200 dark:border-zinc-800 bg-white dark:bg-zinc-900 flex items-center justify-between px-6"
    >
      <div class="flex items-center gap-4">
        <Button
          variant="ghost"
          size="icon"
          class="rounded-xl"
          onclick={() => (isSidebarOpen = !isSidebarOpen)}
        >
          <Menu size={20} />
        </Button>

        <!-- Breadcrumb -->
        {#if isInFamilyContext}
          <nav class="hidden sm:flex items-center gap-1.5 text-sm">
            <a
              href="/dashboard"
              class="text-zinc-400 hover:text-zinc-600 dark:hover:text-zinc-300 transition-colors font-medium"
            >
              Dashboard
            </a>
            <ChevronRight size={14} class="text-zinc-300" />
            <span
              class="text-zinc-700 dark:text-zinc-200 font-bold truncate max-w-[200px]"
            >
              {familyName}
            </span>
          </nav>
        {/if}
      </div>

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
