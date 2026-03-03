<script lang="ts">
  import { env } from "$env/dynamic/public";
  import { Button } from "$lib/components/ui/button/index.js";
  import * as Card from "$lib/components/ui/card/index.js";
  import ThemeToggle from "$lib/components/ThemeToggle.svelte";
  import {
    TreeDeciduous,
    Users,
    Share2,
    ShieldCheck,
    Heart,
  } from "lucide-svelte";

  const loginUrl = `${env.PUBLIC_BACKEND_URL || "http://localhost:8080"}/auth/login`;

  const features = [
    {
      icon: TreeDeciduous,
      label: "Visualisasi Pohon",
      desc: "Lihat silsilah keluarga dalam diagram interaktif yang cantik.",
    },
    {
      icon: Share2,
      label: "Kolaborasi",
      desc: "Undang keluarga lain untuk melengkapi data bersama-sama.",
    },
    {
      icon: ShieldCheck,
      label: "Aman & Privat",
      desc: "Data keluarga Anda dilindungi dengan standar keamanan tinggi.",
    },
  ];
</script>

<div
  class="min-h-screen flex flex-col md:flex-row bg-zinc-50 dark:bg-zinc-950 transition-colors duration-300"
>
  <!-- Left Side: Hero Section (Desktop) -->
  <div
    class="hidden md:flex flex-1 relative overflow-hidden bg-zinc-900 p-12 text-white items-center justify-center"
  >
    <!-- Abstract Background Pattern -->
    <div
      class="absolute inset-0 opacity-20 bg-[radial-gradient(circle_at_center,_var(--tw-gradient-stops))] from-primary/40 via-transparent to-transparent"
    ></div>
    <div
      class="absolute -bottom-24 -left-24 w-96 h-96 bg-primary/20 rounded-full blur-3xl opacity-50"
    ></div>
    <div
      class="absolute -top-24 -right-24 w-64 h-64 bg-primary/10 rounded-full blur-3xl opacity-50"
    ></div>

    <div class="relative z-10 max-w-lg space-y-8 text-center lg:text-left">
      <div
        class="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-white/20 backdrop-blur-sm text-sm font-medium border border-white/30"
      >
        <Heart size={14} fill="currentColor" />
        <span>Warisan Digital Keluarga</span>
      </div>
      <h1 class="text-5xl lg:text-7xl font-bold tracking-tight leading-tight">
        Abadikan Akar <br />
        <span class="text-white/80 italic">Keluarga Anda.</span>
      </h1>
      <p class="text-lg text-white/80 leading-relaxed font-medium">
        Platform modern untuk mencatat, menghubungkan, dan membagikan kisah
        silsilah keluarga Anda lintas generasi dengan aman.
      </p>

      <div class="grid grid-cols-1 gap-4 pt-4">
        {#each features as feature}
          <div class="flex items-center gap-4 group">
            <div
              class="p-3 bg-white/10 rounded-xl backdrop-blur-md group-hover:bg-primary group-hover:text-primary-foreground transition-all duration-300"
            >
              <feature.icon size={24} />
            </div>
            <div class="text-left">
              <h4 class="font-bold text-lg">{feature.label}</h4>
              <p class="text-white/70 text-sm">{feature.desc}</p>
            </div>
          </div>
        {/each}
      </div>
    </div>
  </div>

  <!-- Right Side: Login Section -->
  <div class="flex-1 flex flex-col relative">
    <!-- Floating Theme Toggle -->
    <div class="absolute top-4 right-4 z-50">
      <ThemeToggle />
    </div>

    <div class="flex-1 flex items-center justify-center p-6 md:p-12">
      <div class="w-full max-w-md space-y-8">
        <div class="text-center md:hidden">
          <div
            class="inline-flex items-center gap-2 text-primary font-bold text-3xl mb-2"
          >
            <TreeDeciduous size={32} />
            <span>Silsilah</span>
          </div>
          <p class="text-zinc-500 dark:text-zinc-400">
            Abadikan akar keluarga Anda hari ini.
          </p>
        </div>

        <Card.Root
          class="border-zinc-200 dark:border-zinc-800 shadow-xl bg-white/80 dark:bg-zinc-900/80 backdrop-blur-lg"
        >
          <Card.Header class="space-y-2 text-center">
            <Card.Title class="text-3xl font-bold tracking-tight"
              >Selamat Datang</Card.Title
            >
            <Card.Description class="text-zinc-500 dark:text-zinc-400">
              Gunakan BangLipai Secure Portal untuk masuk dengan aman.
            </Card.Description>
          </Card.Header>
          <Card.Content class="grid gap-6 pt-4">
            <Button
              href={loginUrl}
              class="w-full py-7 text-lg font-bold bg-primary hover:bg-primary/90 transition-all shadow-lg shadow-primary/20 flex items-center gap-3 rounded-xl"
            >
              <ShieldCheck size={24} />
              <span>Login ke Portal Aman</span>
            </Button>

            <div class="relative">
              <div class="absolute inset-0 flex items-center">
                <span
                  class="w-full border-t border-zinc-200 dark:border-zinc-800"
                ></span>
              </div>
              <div class="relative flex justify-center text-xs uppercase">
                <span class="bg-white dark:bg-zinc-900 px-2 text-zinc-500"
                  >Aman & Terenkripsi</span
                >
              </div>
            </div>

            <p
              class="text-center text-xs text-zinc-400 dark:text-zinc-500 px-6"
            >
              Akses akun silsilah Anda dilindungi oleh BangLipai Secure Portal.
              Kami tidak akan pernah membagikan data pribadi Anda ke pihak
              ketiga.
            </p>
          </Card.Content>
        </Card.Root>

        <p class="text-center text-sm text-zinc-500">
          Butuh bantuan akses?
          <a
            href="https://wa.me/6283845257534"
            class="text-primary font-semibold hover:underline">Hubungi Admin</a
          >
        </p>
      </div>
    </div>

    <footer class="p-6 text-center text-xs text-zinc-400 dark:text-zinc-600">
      &copy; {new Date().getFullYear()} Silsilah Keluarga by Kukuh. All rights reserved.
    </footer>
  </div>
</div>
