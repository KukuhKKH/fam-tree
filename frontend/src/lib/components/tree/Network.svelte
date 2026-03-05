<script lang="ts">
  import { onMount, onDestroy } from "svelte";
  import { Network } from "vis-network";
  import { DataSet } from "vis-data";
  import TreeControls from "./parts/TreeControls.svelte";
  import TreeLegend from "./parts/TreeLegend.svelte";
  import TreeSidebar from "./parts/TreeSidebar.svelte";

  let { treeData, familySlug, readOnly = false } = $props();

  let container: HTMLDivElement;
  let network = $state<Network | null>(null);
  let selectedPersonId = $state<number | null>(null);
  let searchQuery = $state("");
  let showLegend = $state(false);

  // Status indicators for formatting
  const getGenderColor = (gender: string) =>
    gender === "male" ? "blue" : "rose";
  const getGenderLabel = (gender: string) =>
    gender === "male" ? "Laki-laki" : "Perempuan";

  const relationshipLabels: Record<string, string> = {
    spouse: "Pasangan",
    sibling: "Saudara",
    parent_child: "Orang Tua - Anak",
    "parent-child": "Orang Tua - Anak",
  };

  const selectedPerson = $derived(
    treeData.nodes.find((n: any) => n.id === selectedPersonId),
  );

  const relatives = $derived(
    selectedPersonId
      ? treeData.edges
          .filter(
            (e: any) =>
              e.from === selectedPersonId || e.to === selectedPersonId,
          )
          .map((e: any) => {
            const isFrom = e.from === selectedPersonId;
            const otherId = isFrom ? e.to : e.from;
            const otherPerson = treeData.nodes.find(
              (n: any) => n.id === otherId,
            );
            let label = "";

            if (
              e.relationship_type === "parent_child" ||
              e.relationship_type === "parent-child"
            ) {
              if (isFrom) {
                label = "Anak";
              } else {
                label = otherPerson.gender === "male" ? "Ayah" : "Ibu";
              }
            } else if (e.relationship_type === "spouse") {
              label = otherPerson.gender === "male" ? "Suami" : "Istri";
            } else if (e.relationship_type === "sibling") {
              label = "Saudara";
            } else {
              label =
                relationshipLabels[e.relationship_type] || e.relationship_type;
            }

            return {
              ...e,
              person: otherPerson,
              displayLabel: label,
            };
          })
      : [],
  );

  onMount(() => {
    if (!container) return;

    const parentChildEdges = treeData.edges.filter(
      (e: any) =>
        e.relationship_type === "parent_child" ||
        e.relationship_type === "parent-child",
    );
    const spouseEdges = treeData.edges.filter(
      (e: any) => e.relationship_type === "spouse",
    );

    // Map: parentID -> [childIDs]
    const childrenOf: Record<number, number[]> = {};
    for (const e of parentChildEdges) {
      if (!childrenOf[e.from]) childrenOf[e.from] = [];
      childrenOf[e.from].push(e.to);
    }

    // Map: personID -> spouseID (bidirectional)
    const spouseOf: Record<number, number> = {};
    for (const e of spouseEdges) {
      spouseOf[e.from] = e.to;
      spouseOf[e.to] = e.from;
    }

    // BFS from root to assign levels
    const levels: Record<number, number> = {};
    const rootID = treeData.root_id || treeData.nodes[0]?.id;

    if (rootID) {
      const queue: number[] = [rootID];
      levels[rootID] = 0;

      if (spouseOf[rootID] !== undefined) {
        levels[spouseOf[rootID]] = 0;
        queue.push(spouseOf[rootID]);
      }

      const visited = new Set<number>(queue);

      while (queue.length > 0) {
        const current = queue.shift()!;
        const currentLevel = levels[current];

        const children = childrenOf[current] || [];
        for (const childId of children) {
          if (visited.has(childId)) continue;

          levels[childId] = currentLevel + 1;
          visited.add(childId);
          queue.push(childId);

          if (
            spouseOf[childId] !== undefined &&
            !visited.has(spouseOf[childId])
          ) {
            levels[spouseOf[childId]] = currentLevel + 1;
            visited.add(spouseOf[childId]);
            queue.push(spouseOf[childId]);
          }
        }
      }

      for (const n of treeData.nodes) {
        if (levels[n.id] === undefined) {
          levels[n.id] = 0;
        }
      }
    }

    // ─── Manual Position Calculation
    const SPOUSE_GAP = 140; // tight gap between husband & wife
    const FAMILY_GAP = 400; // wide gap between family units
    const LEVEL_HEIGHT = 240; // vertical distance between generations

    const positions: Record<number, { x: number; y: number }> = {};
    const maxLevel = Math.max(...Object.values(levels), 0);

    // Helper: find parent couple midpoint for a node
    function getParentMidX(nodeId: number): number {
      for (const e of parentChildEdges) {
        if (e.to === nodeId) {
          const px = positions[e.from];
          if (px) {
            const sid = spouseOf[e.from];
            if (sid !== undefined && positions[sid]) {
              return (px.x + positions[sid].x) / 2;
            }
            return px.x;
          }
        }
      }
      return 0;
    }

    for (let lv = 0; lv <= maxLevel; lv++) {
      const nodesAtLevel = treeData.nodes
        .filter((n: any) => levels[n.id] === lv)
        .map((n: any) => n.id);

      // Group nodes into couples and singles
      const placed = new Set<number>();
      const groups: number[][] = [];

      for (const nodeId of nodesAtLevel) {
        if (placed.has(nodeId)) continue;
        placed.add(nodeId);

        const sid = spouseOf[nodeId];
        if (sid !== undefined && levels[sid] === lv && !placed.has(sid)) {
          groups.push([nodeId, sid]);
          placed.add(sid);
        } else {
          groups.push([nodeId]);
        }
      }

      // Sort groups by parent position (children appear under parents)
      if (lv > 0) {
        groups.sort((a, b) => getParentMidX(a[0]) - getParentMidX(b[0]));
      }

      // Assign x positions
      let currentX = 0;
      for (const group of groups) {
        if (group.length === 2) {
          positions[group[0]] = {
            x: currentX,
            y: lv * LEVEL_HEIGHT,
          };
          positions[group[1]] = {
            x: currentX + SPOUSE_GAP,
            y: lv * LEVEL_HEIGHT,
          };
          currentX += SPOUSE_GAP + FAMILY_GAP;
        } else {
          positions[group[0]] = {
            x: currentX,
            y: lv * LEVEL_HEIGHT,
          };
          currentX += FAMILY_GAP;
        }
      }

      // Center the row around x=0
      const allX = nodesAtLevel.map((id: number) => positions[id]?.x ?? 0);
      if (allX.length > 0) {
        const midX = (Math.min(...allX) + Math.max(...allX)) / 2;
        for (const id of nodesAtLevel) {
          if (positions[id]) positions[id].x -= midX;
        }
      }
    }

    // Second pass: center child groups under parent couple midpoints
    for (let lv = 1; lv <= maxLevel; lv++) {
      const nodesAtLevel = treeData.nodes
        .filter((n: any) => levels[n.id] === lv)
        .map((n: any) => n.id);

      // Group by parent couple
      const parentGroups: Map<string, number[]> = new Map();
      for (const nodeId of nodesAtLevel) {
        const pmx = getParentMidX(nodeId);
        const key = pmx.toFixed(0);
        if (!parentGroups.has(key)) parentGroups.set(key, []);
        parentGroups.get(key)!.push(nodeId);
      }

      for (const [keyStr, childIds] of parentGroups) {
        const parentMidX = parseFloat(keyStr);
        const childXs = childIds.map((id) => positions[id]?.x ?? 0);
        const childMidX = (Math.min(...childXs) + Math.max(...childXs)) / 2;
        const offset = parentMidX - childMidX;

        for (const id of childIds) {
          if (positions[id]) positions[id].x += offset;
        }
      }
    }

    // ─── Transform data for vis-network
    const nodes = new DataSet(
      treeData.nodes.map((n: any) => ({
        id: n.id,
        x: positions[n.id]?.x ?? 0,
        y: positions[n.id]?.y ?? 0,
        label: n.nickname || n.full_name.split(" ")[0],
        title: `
          <div class="px-3 py-2 space-y-1">
            <div class="font-bold text-sm leading-tight text-zinc-900">${n.full_name}</div>
            <div class="text-[10px] font-black uppercase tracking-wider text-zinc-500">${getGenderLabel(n.gender)}</div>
          </div>
        `,
        shape: "circularImage",
        image:
          n.photo_url ||
          `https://ui-avatars.com/api/?name=${encodeURIComponent(n.full_name)}&background=${n.gender === "male" ? "0ea5e9" : "f43f5e"}&color=fff&bold=true&size=128`,
        borderWidth: 4,
        size: 35,
        color: {
          border: n.gender === "male" ? "#0ea5e9" : "#f43f5e",
          background: "#ffffff",
          highlight: {
            border: "#fbbf24",
            background: "#ffffff",
          },
        },
        font: {
          color: "#4b5563",
          size: 14,
          face: "Inter, sans-serif",
          base: "bold",
          strokeWidth: 4,
          strokeColor: "#ffffff",
        },
      })),
    );

    // Deduplicate parent_child edges
    const seenChildEdges = new Set<number>();
    const dedupedEdges = treeData.edges.filter((e: any) => {
      const isParentChild =
        e.relationship_type === "parent_child" ||
        e.relationship_type === "parent-child";

      if (isParentChild) {
        if (seenChildEdges.has(e.to)) return false;
        seenChildEdges.add(e.to);
      }
      return true;
    });

    const edges = new DataSet(
      dedupedEdges.map((e: any) => {
        const type = e.relationship_type;
        const isParentChild =
          type === "parent_child" || type === "parent-child";

        return {
          from: e.from,
          to: e.to,
          label: isParentChild ? "" : relationshipLabels[type] || type,
          arrows: isParentChild ? "to" : "",
          color: {
            color: isParentChild ? "#94a3b8" : "#fbbf24",
            highlight: "#fbbf24",
          },
          dashes: type === "spouse" ? true : false,
          width: type === "spouse" ? 2.5 : 1.5,
          font: {
            size: 11,
            face: "Inter, sans-serif",
            color: "#d97706",
            align: "top",
            strokeWidth: 3,
            strokeColor: "#ffffff",
          },
        };
      }),
    );

    const data = { nodes, edges };

    const options = {
      nodes: {
        shapeProperties: {
          useBorderWithImage: true,
        },
      },
      edges: {
        smooth: {
          type: "cubicBezier",
          forceDirection: "vertical",
          roundness: 0.5,
        },
      },
      layout: {
        // No hierarchical — we use manual x,y positions
      },
      physics: {
        enabled: false,
      },
      interaction: {
        hover: true,
        tooltipDelay: 100,
        zoomView: true,
      },
    };

    network = new Network(container, data, options);

    network.on("selectNode", (params: any) => {
      selectedPersonId = params.nodes[0];
      network?.focus(params.nodes[0], {
        scale: 1.1,
        animation: {
          duration: 1000,
          easingFunction: "easeInOutQuart",
        },
      });
    });

    network.on("deselectNode", () => {
      selectedPersonId = null;
    });

    // Initial Focus
    network.once("afterDrawing", () => {
      if (treeData.root_id) {
        network?.focus(treeData.root_id, {
          scale: 0.9,
          animation: { duration: 1000, easingFunction: "easeInOutQuad" },
        });
      } else {
        network?.fit({
          animation: { duration: 1000, easingFunction: "easeInOutQuad" },
        });
      }
    });
  });

  onDestroy(() => {
    network?.destroy();
  });

  function handleSearch() {
    if (!network || !searchQuery) return;
    const person = treeData.nodes.find(
      (n: any) =>
        n.full_name.toLowerCase().includes(searchQuery.toLowerCase()) ||
        n.nickname?.toLowerCase().includes(searchQuery.toLowerCase()),
    );
    if (person) {
      network.focus(person.id, {
        scale: 1.2,
        animation: { duration: 1200, easingFunction: "easeInOutQuart" },
      });
      network.selectNodes([person.id]);
      selectedPersonId = person.id;
    }
  }

  function resetView() {
    network?.fit({
      animation: { duration: 1000, easingFunction: "easeInOutQuad" },
    });
  }
</script>

<div
  class="relative w-full h-[calc(100vh-12rem)] min-h-[450px] sm:min-h-[600px] rounded-[2rem] sm:rounded-[3.5rem] bg-zinc-50 dark:bg-zinc-950/20 border border-zinc-200 dark:border-zinc-800 shadow-inner overflow-hidden"
>
  <!-- The Graph Container -->
  <div bind:this={container} class="w-full h-full"></div>

  <!-- Search & Controls Overlay -->
  <TreeControls
    bind:searchQuery
    bind:showLegend
    onSearch={handleSearch}
    onReset={resetView}
  />

  <!-- Legend Overlay -->
  <TreeLegend bind:showLegend />

  <!-- Selection Info Sidebar -->
  <TreeSidebar
    bind:selectedPersonId
    {selectedPerson}
    {relatives}
    {familySlug}
    {network}
    {getGenderColor}
    {getGenderLabel}
    {readOnly}
  />
</div>

<style>
  :global(.vis-network) {
    outline: none;
  }
  :global(.vis-network:focus) {
    outline: none;
  }

  :global(.vis-tooltip) {
    background: rgba(255, 255, 255, 0.8) !important;
    backdrop-filter: blur(20px) !important;
    border-radius: 20px !important;
    padding: 0 !important;
    border: 1px solid rgba(0, 0, 0, 0.05) !important;
    box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.15) !important;
    overflow: hidden !important;
    pointer-events: none !important;
  }
</style>
