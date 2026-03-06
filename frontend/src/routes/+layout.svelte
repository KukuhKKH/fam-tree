<script lang="ts">
  import "../app.css";
  import favicon from "$lib/assets/favicon.svg";
  import { Toaster } from "svelte-sonner";

  let { children } = $props();
</script>

<svelte:head>
  <title>Silsilah — Platform Silsilah Keluarga Mandiri</title>
  <meta
    name="description"
    content="Silsilah adalah platform modern untuk membangun, mengelola, dan berbagi pohon silsilah keluarga secara mandiri dan aman."
  />
  <link rel="icon" href={favicon} />

  <!-- Open Graph / Facebook -->
  <meta property="og:type" content="website" />
  <meta
    property="og:title"
    content="Silsilah — Platform Silsilah Keluarga Mandiri"
  />
  <meta
    property="og:description"
    content="Bangun dan bagikan silsilah keluarga Anda dengan mudah dan elegan menggunakan Silsilah."
  />
  <meta property="og:site_name" content="Silsilah" />
  <meta property="og:image" content="/og-image.png" />

  <!-- Twitter -->
  <meta property="twitter:card" content="summary_large_image" />
  <meta
    property="twitter:title"
    content="Silsilah — Platform Silsilah Keluarga Mandiri"
  />
  <meta
    property="twitter:description"
    content="Bangun dan bagikan silsilah keluarga Anda dengan mudah dan elegan menggunakan Silsilah."
  />
  <meta property="twitter:image" content="/og-image.png" />

  <script>
    // Anti-flash script: Jalankan sebelum Svelte merender apapun
    (function () {
      try {
        const theme = localStorage.getItem("theme");
        const supportDarkMode = window.matchMedia(
          "(prefers-color-scheme: dark)",
        ).matches;
        if (theme === "dark" || (!theme && supportDarkMode)) {
          document.documentElement.classList.add("dark");
        } else {
          document.documentElement.classList.remove("dark");
        }
      } catch (e) {}
    })();
  </script>
</svelte:head>

<Toaster position="top-right" richColors closeButton />
{@render children()}
