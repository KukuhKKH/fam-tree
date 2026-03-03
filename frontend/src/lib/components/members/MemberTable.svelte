<script lang="ts">
  import * as Table from "$lib/components/ui/table/index.js";
  import * as DropdownMenu from "$lib/components/ui/dropdown-menu/index.js";
  import { Button } from "$lib/components/ui/button/index.js";
  import { Badge } from "$lib/components/ui/badge/index.js";
  import {
    MoreHorizontal,
    UserMinus,
    Shield,
    ShieldCheck,
    User,
    Eye,
    Loader2,
  } from "lucide-svelte";
  import { apiFetch } from "$lib/api";
  import { toast } from "svelte-sonner";
  import { invalidateAll } from "$app/navigation";

  let { members = [], familySlug, myRole } = $props();

  let loadingId = $state<number | null>(null);

  const roles = [
    { value: "owner", label: "Owner", icon: ShieldCheck, color: "bg-red-500" },
    { value: "admin", label: "Admin", icon: Shield, color: "bg-orange-500" },
    { value: "editor", label: "Editor", icon: User, color: "bg-blue-500" },
    { value: "viewer", label: "Viewer", icon: Eye, color: "bg-zinc-500" },
  ];

  async function updateRole(userId: number, newRole: string) {
    loadingId = userId;
    try {
      await apiFetch(`/families/${familySlug}/members/${userId}`, {
        method: "PATCH",
        body: JSON.stringify({ role: newRole }),
      });

      toast.success("Role anggota berhasil diperbarui");
      await invalidateAll();
    } catch (err: any) {
      toast.error(err.message || "Gagal memperbarui role");
    } finally {
      loadingId = null;
    }
  }

  async function removeMember(userId: number, userName: string) {
    if (
      !confirm(
        `Apakah Anda yakin ingin mengeluarkan ${userName} dari keluarga ini?`,
      )
    )
      return;

    loadingId = userId;
    try {
      await apiFetch(`/families/${familySlug}/members/${userId}`, {
        method: "DELETE",
      });

      toast.success(`${userName} berhasil dikeluarkan`);
      await invalidateAll();
    } catch (err: any) {
      toast.error(err.message || "Gagal mengeluarkan anggota");
    } finally {
      loadingId = null;
    }
  }

  const canManage = $derived(myRole === "owner" || myRole === "admin");

  function getRoleBadge(role: string) {
    return roles.find((r) => r.value === role) || roles[3];
  }
</script>

<div
  class="rounded-xl border border-zinc-200 dark:border-zinc-800 bg-white dark:bg-zinc-900 overflow-hidden shadow-sm"
>
  <Table.Root>
    <Table.Header class="bg-zinc-50 dark:bg-zinc-950">
      <Table.Row>
        <Table.Head class="w-[250px]">Nama</Table.Head>
        <Table.Head>Email</Table.Head>
        <Table.Head>Role</Table.Head>
        <Table.Head>Bergabung</Table.Head>
        <Table.Head class="text-right">Aksi</Table.Head>
      </Table.Row>
    </Table.Header>
    <Table.Body>
      {#each members as member (member.user_id)}
        <Table.Row class="group transition-colors">
          <Table.Cell class="font-medium">
            <div class="flex items-center gap-3">
              <div
                class="w-8 h-8 rounded-full bg-primary/10 text-primary flex items-center justify-center text-xs font-bold uppercase"
              >
                {member.user_name.substring(0, 2)}
              </div>
              <span>{member.user_name}</span>
            </div>
          </Table.Cell>
          <Table.Cell class="text-zinc-500">{member.user_email}</Table.Cell>
          <Table.Cell>
            <Badge variant="secondary" class="gap-1.5 font-medium capitalize">
              {@const r = getRoleBadge(member.role)}
              <r.icon
                size={12}
                class={r.value === "owner" ? "text-red-500" : "text-zinc-500"}
              />
              {member.role}
            </Badge>
          </Table.Cell>
          <Table.Cell class="text-zinc-400 text-sm">
            {new Date(member.joined_at).toLocaleDateString("id-ID", {
              day: "numeric",
              month: "short",
              year: "numeric",
            })}
          </Table.Cell>
          <Table.Cell class="text-right">
            {#if canManage && member.role !== "owner"}
              <DropdownMenu.Root>
                <DropdownMenu.Trigger>
                  {#snippet child({ props })}
                    <Button
                      {...props}
                      variant="ghost"
                      size="icon"
                      disabled={loadingId === member.user_id}
                      class="cursor-pointer"
                    >
                      {#if loadingId === member.user_id}
                        <Loader2 class="h-4 w-4 animate-spin" />
                      {:else}
                        <MoreHorizontal class="h-4 w-4" />
                      {/if}
                    </Button>
                  {/snippet}
                </DropdownMenu.Trigger>
                <DropdownMenu.Content align="end" class="w-48 backdrop-blur-md">
                  <DropdownMenu.Label>Ubah Role</DropdownMenu.Label>
                  <DropdownMenu.Separator />
                  {#each roles.filter((r) => r.value !== "owner") as r}
                    <DropdownMenu.Item
                      onclick={() => updateRole(member.user_id, r.value)}
                      class="gap-2 {member.role === r.value
                        ? 'bg-primary/5 text-primary'
                        : ''} cursor-pointer"
                    >
                      <r.icon size={14} />
                      <span>{r.label}</span>
                    </DropdownMenu.Item>
                  {/each}
                  <DropdownMenu.Separator />
                  <DropdownMenu.Item
                    onclick={() =>
                      removeMember(member.user_id, member.user_name)}
                    class="text-destructive gap-2 focus:text-destructive cursor-pointer"
                  >
                    <UserMinus size={14} />
                    <span>Keluarkan</span>
                  </DropdownMenu.Item>
                </DropdownMenu.Content>
              </DropdownMenu.Root>
            {:else if member.role === "owner"}
              <span class="text-xs text-zinc-400 italic pr-3 font-medium"
                >Owner</span
              >
            {:else}
              <span class="text-xs text-zinc-400 italic pr-3">No Access</span>
            {/if}
          </Table.Cell>
        </Table.Row>
      {/each}
      {#if members.length === 0}
        <Table.Row>
          <Table.Cell colspan={5} class="h-32 text-center text-zinc-500">
            Belum ada anggota keluarga di daftar ini.
          </Table.Cell>
        </Table.Row>
      {/if}
    </Table.Body>
  </Table.Root>
</div>
