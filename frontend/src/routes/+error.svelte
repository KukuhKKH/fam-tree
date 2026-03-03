<script lang="ts">
  import { page } from "$app/state";
  import {
    TreeDeciduous,
    ArrowLeft,
    Home,
    RefreshCcw,
    LifeBuoy,
  } from "lucide-svelte";
  import { fade, fly, scale } from "svelte/transition";
  import { backOut } from "svelte/easing";

  const status = $derived(page.status || 404);
  const message = $derived(
    page.error?.message || "Terjadi kesalahan yang tidak terduga.",
  );

  const errorDetails: Record<
    number | string,
    { title: string; desc: string; iconColor: string }
  > = {
    400: {
      title: "Bad Request",
      desc: "Permintaan tidak valid.",
      iconColor: "text-amber-500",
    },
    401: {
      title: "Unauthorized",
      desc: "Sesi Anda berakhir, silakan login kembali.",
      iconColor: "text-rose-500",
    },
    403: {
      title: "Forbidden",
      desc: "Area terlarang! Anda tidak punya akses ke sini.",
      iconColor: "text-red-600",
    },
    404: {
      title: "Jejak Terputus",
      desc: "Halaman yang Anda cari tidak ditemukan dalam silsilah kami.",
      iconColor: "text-indigo-500",
    },
    500: {
      title: "Server Lelah",
      desc: "Server kami sedang istirahat sejenak. Coba lagi nanti.",
      iconColor: "text-red-500",
    },
    default: {
      title: "Oops!",
      desc: "Sesuatu yang salah telah terjadi.",
      iconColor: "text-primary",
    },
  };

  const detail = $derived(errorDetails[status] || errorDetails.default);
</script>

<div
  class="min-h-screen relative flex flex-col items-center justify-center bg-zinc-50 dark:bg-zinc-950 p-6 overflow-hidden transition-colors duration-500"
>
  <div
    class="absolute top-1/4 -left-20 w-72 h-72 bg-primary/10 rounded-full blur-[120px] pointer-events-none"
  ></div>
  <div
    class="absolute bottom-1/4 -right-20 w-72 h-72 bg-indigo-500/10 rounded-full blur-[120px] pointer-events-none"
  ></div>

  <div class="relative z-10 max-w-2xl w-full flex flex-col items-center">
    <div
      in:scale={{ duration: 600, easing: backOut }}
      class="relative mb-8 text-center"
    >
      <h1
        class="text-[12rem] md:text-[16rem] font-black leading-none tracking-tighter text-zinc-900/5 dark:text-white/5 select-none"
      >
        {status}
      </h1>

      <div
        in:fly={{ y: 20, delay: 300 }}
        class="absolute inset-0 flex items-center justify-center"
      >
        <div
          class="p-6 bg-white/40 dark:bg-zinc-900/40 backdrop-blur-xl rounded-3xl border border-white/20 dark:border-zinc-800 shadow-2xl shadow-primary/10 transition-transform hover:scale-110 duration-500"
        >
          <TreeDeciduous
            size={80}
            class="{detail.iconColor} transition-colors"
          />
        </div>
      </div>
    </div>

    <div in:fade={{ delay: 500 }} class="text-center space-y-4 max-w-md">
      <h2
        class="text-4xl font-extrabold text-zinc-900 dark:text-zinc-100 tracking-tight"
      >
        {detail.title}
      </h2>
      <p class="text-lg text-zinc-500 dark:text-zinc-400 font-medium">
        {message !== detail.title ? message : detail.desc}
      </p>
    </div>

    <div
      in:fly={{ y: 40, delay: 700 }}
      class="grid grid-cols-1 sm:grid-cols-2 gap-4 mt-12 w-full max-w-lg"
    >
      <button
        onclick={() => history.back()}
        class="flex items-center gap-4 p-4 rounded-2xl bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 hover:border-primary/50 transition-all group text-left cursor-pointer"
      >
        <div
          class="p-3 bg-zinc-100 dark:bg-zinc-800 rounded-xl group-hover:bg-primary/10 group-hover:text-primary transition-colors"
        >
          <ArrowLeft size={20} />
        </div>
        <div>
          <div class="font-bold text-zinc-900 dark:text-zinc-100">Kembali</div>
          <div class="text-xs text-zinc-500">Ke halaman sebelumnya</div>
        </div>
      </button>

      <a
        href="/dashboard"
        class="flex items-center gap-4 p-4 rounded-2xl bg-primary text-primary-foreground shadow-lg shadow-primary/20 hover:opacity-90 transition-all group text-left cursor-pointer"
      >
        <div
          class="p-3 bg-white/20 rounded-xl group-hover:scale-110 transition-transform"
        >
          <Home size={20} />
        </div>
        <div>
          <div class="font-bold">Dashboard</div>
          <div class="text-xs opacity-80">Kembali ke beranda utama</div>
        </div>
      </a>
    </div>

    <div
      in:fade={{ delay: 1000 }}
      class="mt-16 flex flex-col items-center gap-4"
    >
      <div class="flex items-center gap-6 text-zinc-400 dark:text-zinc-600">
        <a
          href="https://wa.me/6283845257534"
          target="_blank"
          rel="noopener noreferrer"
          class="flex items-center gap-2 hover:text-primary transition-colors cursor-pointer"
        >
          <LifeBuoy size={16} /> <span>Pusat Bantuan</span>
        </a>
        <button
          onclick={() => window.location.reload()}
          class="flex items-center gap-2 hover:text-primary transition-colors cursor-pointer"
        >
          <RefreshCcw size={16} /> <span>Muat Ulang</span>
        </button>
      </div>

      {#if status >= 500}
        <div
          class="px-4 py-2 rounded-full bg-zinc-100 dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 text-[10px] font-mono text-zinc-500"
        >
          LOG_ID: {Math.random().toString(36).substring(7).toUpperCase()}
        </div>
      {/if}
    </div>
  </div>
</div>
