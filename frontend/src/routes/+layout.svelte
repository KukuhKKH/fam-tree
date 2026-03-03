<script lang="ts">
  import "../app.css";
  import favicon from "$lib/assets/favicon.svg";

  let { children } = $props();
</script>

<svelte:head>
  <link rel="icon" href={favicon} />
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

{@render children()}
