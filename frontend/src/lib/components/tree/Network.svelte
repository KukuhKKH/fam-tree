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
            const type = e.relationship_type?.replace("_", "-");
            const isFrom = e.from === selectedPersonId;
            const otherId = isFrom ? e.to : e.from;
            const otherPerson = treeData.nodes.find(
              (n: any) => n.id === otherId,
            );

            let label = "";

            if (type === "parent-child") {
              if (isFrom) {
                label = "Anak";
              } else {
                label = otherPerson?.gender === "male" ? "Ayah" : "Ibu";
              }
            } else if (type === "spouse") {
              label = otherPerson?.gender === "male" ? "Suami" : "Istri";
            } else if (type === "sibling") {
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

    network?.destroy();
    network = null;

    const normalizeType = (t: string) => t.replace("_", "-");

    type PersonNode = {
      id: number;
      full_name: string;
      nickname?: string;
      gender: string;
      photo_url?: string;
      is_alive: boolean;
      birth_date?: string;
      death_date?: string;
    };

    type FamilyUnit = {
      id: string;
      parents: number[];
      children: number[];
      level: number;
    };

    const people = treeData.nodes as PersonNode[];
    const edges = treeData.edges as any[];

    const personMap = new Map<number, PersonNode>();
    people.forEach((p) => personMap.set(p.id, p));

    const parentsOf = new Map<number, number[]>();
    const childrenOf = new Map<number, number[]>();
    const spouseSet = new Map<number, Set<number>>();
    const spousePairs = new Set<string>();

    function addMapArray(
      map: Map<number, number[]>,
      key: number,
      value: number,
    ) {
      if (!map.has(key)) map.set(key, []);
      const arr = map.get(key)!;
      if (!arr.includes(value)) arr.push(value);
    }

    function addSpouse(a: number, b: number) {
      if (!spouseSet.has(a)) spouseSet.set(a, new Set());
      if (!spouseSet.has(b)) spouseSet.set(b, new Set());

      spouseSet.get(a)!.add(b);
      spouseSet.get(b)!.add(a);

      const key = a < b ? `${a}-${b}` : `${b}-${a}`;
      spousePairs.add(key);
    }

    edges.forEach((e: any) => {
      const type = normalizeType(e.relationship_type || "");

      if (type === "parent-child") {
        addMapArray(childrenOf, e.from, e.to);
        addMapArray(parentsOf, e.to, e.from);
      } else if (type === "spouse") {
        addSpouse(e.from, e.to);
      }
    });

    const familyUnits = new Map<string, FamilyUnit>();
    const personParentFamily = new Map<number, string>();

    function makeFamilyId(parents: number[]) {
      return `family-${[...parents].sort((a, b) => a - b).join("-")}`;
    }

    function ensureFamily(parents: number[]) {
      const sortedParents = [...parents].sort((a, b) => a - b);
      const familyId = makeFamilyId(sortedParents);

      if (!familyUnits.has(familyId)) {
        familyUnits.set(familyId, {
          id: familyId,
          parents: sortedParents,
          children: [],
          level: 0,
        });
      }

      return familyId;
    }

    people.forEach((child) => {
      const rawParents = [...(parentsOf.get(child.id) || [])].sort(
        (a, b) => a - b,
      );

      if (rawParents.length === 0) return;

      let finalParents = rawParents;

      if (rawParents.length === 1) {
        const p1 = rawParents[0];
        const spouses = spouseSet.get(p1);
        if (spouses && spouses.size > 0) {
          const s1 = [...spouses][0];
          finalParents = [p1, s1].sort((a, b) => a - b);
        }
      } else if (rawParents.length > 2) {
        finalParents = rawParents.slice(0, 2);
      }

      const familyId = ensureFamily(finalParents);
      const fam = familyUnits.get(familyId)!;

      if (!fam.children.includes(child.id)) {
        fam.children.push(child.id);
      }

      personParentFamily.set(child.id, familyId);
    });

    spousePairs.forEach((pairKey) => {
      const parents = pairKey.split("-").map(Number);
      ensureFamily(parents);
    });

    const usedPeople = new Set<number>();
    familyUnits.forEach((f) => {
      f.parents.forEach((id) => usedPeople.add(id));
      f.children.forEach((id) => usedPeople.add(id));
    });

    people.forEach((p) => {
      if (!usedPeople.has(p.id)) {
        ensureFamily([p.id]);
      }
    });

    const personParentFamilies = new Map<number, string[]>();
    familyUnits.forEach((f) => {
      f.parents.forEach((p) => {
        if (!personParentFamilies.has(p)) personParentFamilies.set(p, []);
        const arr = personParentFamilies.get(p)!;
        if (!arr.includes(f.id)) arr.push(f.id);
      });
    });

    const personLevel = new Map<number, number>();
    const rootPeople = people
      .filter((p) => !parentsOf.get(p.id)?.length)
      .map((p) => p.id);

    const queue: Array<{ personId: number; level: number }> = rootPeople.map(
      (id) => ({
        personId: id,
        level: 0,
      }),
    );

    while (queue.length > 0) {
      const { personId, level } = queue.shift()!;
      const current = personLevel.get(personId);

      if (current !== undefined && current >= level) continue;

      personLevel.set(personId, level);

      const spouses = spouseSet.get(personId);
      if (spouses) {
        spouses.forEach((sid) => {
          queue.push({ personId: sid, level });
        });
      }

      const kids = childrenOf.get(personId) || [];
      kids.forEach((cid) => {
        queue.push({ personId: cid, level: level + 1 });
      });
    }

    people.forEach((p) => {
      if (!personLevel.has(p.id)) personLevel.set(p.id, 0);
    });

    familyUnits.forEach((f) => {
      const lv = Math.max(...f.parents.map((p) => personLevel.get(p) ?? 0));
      f.level = lv;
    });

    const familyChildren = new Map<string, string[]>();
    const familyParents = new Map<string, string[]>();

    function addFamilyEdge(parentFamilyId: string, childFamilyId: string) {
      if (parentFamilyId === childFamilyId) return;

      if (!familyChildren.has(parentFamilyId))
        familyChildren.set(parentFamilyId, []);
      if (!familyParents.has(childFamilyId))
        familyParents.set(childFamilyId, []);

      const fc = familyChildren.get(parentFamilyId)!;
      const fp = familyParents.get(childFamilyId)!;

      if (!fc.includes(childFamilyId)) fc.push(childFamilyId);
      if (!fp.includes(parentFamilyId)) fp.push(parentFamilyId);
    }

    familyUnits.forEach((f) => {
      f.children.forEach((childId) => {
        const ownFamilyIds = personParentFamilies.get(childId) || [];
        ownFamilyIds.forEach((ownId) => {
          addFamilyEdge(f.id, ownId);
        });
      });
    });

    const SPOUSE_GAP = 160;
    const NODE_WIDTH = 120;
    const CHILD_SPACING = 80;
    const LEVEL_HEIGHT = 320;
    const FAMILY_NODE_OFFSET_Y = 120;

    const subtreeWidth = new Map<string, number>();

    function getFamilySelfWidth(f: FamilyUnit) {
      return f.parents.length >= 2 ? SPOUSE_GAP + NODE_WIDTH : NODE_WIDTH;
    }

    function getChildBlockWidth(childId: number): number {
      const ownFamilyIds = personParentFamilies.get(childId) || [];
      if (ownFamilyIds.length > 0) {
        let total = 0;
        ownFamilyIds.forEach((fid) => (total += getSubtreeWidth(fid)));
        total += (ownFamilyIds.length - 1) * CHILD_SPACING;
        return total;
      }
      return NODE_WIDTH;
    }

    function getSubtreeWidth(familyId: string): number {
      if (subtreeWidth.has(familyId)) return subtreeWidth.get(familyId)!;

      const f = familyUnits.get(familyId);
      if (!f) return 0;
      const ownWidth = getFamilySelfWidth(f);

      if (f.children.length === 0) {
        subtreeWidth.set(familyId, ownWidth);
        return ownWidth;
      }

      let childrenWidth = 0;
      f.children.forEach((childId) => {
        childrenWidth += getChildBlockWidth(childId);
      });
      childrenWidth += (f.children.length - 1) * CHILD_SPACING;

      const finalWidth = Math.max(ownWidth, childrenWidth);
      subtreeWidth.set(familyId, finalWidth);
      return finalWidth;
    }

    const positions: Record<number, { x: number; y: number }> = {};
    const familyNodePositions: Record<string, { x: number; y: number }> = {};
    const placedFamily = new Set<string>();

    function positionFamily(familyId: string, xStart: number) {
      if (placedFamily.has(familyId)) return;
      placedFamily.add(familyId);

      const f = familyUnits.get(familyId)!;
      const width = getSubtreeWidth(familyId);

      const parentY = f.level * LEVEL_HEIGHT;
      const familyY = parentY + FAMILY_NODE_OFFSET_Y;
      const childY = parentY + LEVEL_HEIGHT;
      const centerX = xStart + width / 2;

      // Position parents
      if (f.parents.length >= 2) {
        positions[f.parents[0]] = { x: centerX - SPOUSE_GAP / 2, y: parentY };
        positions[f.parents[1]] = { x: centerX + SPOUSE_GAP / 2, y: parentY };
      } else {
        positions[f.parents[0]] = { x: centerX, y: parentY };
      }

      familyNodePositions[familyId] = { x: centerX, y: familyY };

      if (f.children.length === 0) return;

      // Position children
      let totalChildrenWidth = 0;
      f.children.forEach((childId) => {
        totalChildrenWidth += getChildBlockWidth(childId);
      });
      totalChildrenWidth += (f.children.length - 1) * CHILD_SPACING;

      let childX = centerX - totalChildrenWidth / 2;

      f.children.forEach((childId) => {
        const ownFamilyIds = personParentFamilies.get(childId) || [];
        const blockWidth = getChildBlockWidth(childId);

        if (ownFamilyIds.length > 0) {
          let currX = childX;
          ownFamilyIds.forEach((fid) => {
            positionFamily(fid, currX);
            currX += getSubtreeWidth(fid) + CHILD_SPACING;
          });

          // Position the child person relative to their families
          if (!positions[childId]) {
            const firstFam = familyUnits.get(ownFamilyIds[0])!;
            const famCenterX = childX + getSubtreeWidth(ownFamilyIds[0]) / 2;
            positions[childId] =
              firstFam.parents.length >= 2
                ? { x: famCenterX - SPOUSE_GAP / 2, y: childY }
                : { x: famCenterX, y: childY };
          }
        } else {
          positions[childId] = { x: childX + blockWidth / 2, y: childY };
        }
        childX += blockWidth + CHILD_SPACING;
      });
    }

    const rootFamilies = [...familyUnits.values()]
      .filter((f) => !familyParents.get(f.id)?.length)
      .sort((a, b) => a.level - b.level || a.parents[0] - b.parents[0]);

    let currentX = 0;
    rootFamilies.forEach((f) => {
      positionFamily(f.id, currentX);
      currentX += getSubtreeWidth(f.id) + CHILD_SPACING * 10;
    });

    people.forEach((p) => {
      if (!positions[p.id]) {
        positions[p.id] = {
          x: currentX,
          y: ((personLevel.get(p.id) ?? 0) + 2) * LEVEL_HEIGHT,
        };
        currentX += NODE_WIDTH + CHILD_SPACING;
      }
    });

    const visNodes: any[] = people.map((n: PersonNode) => {
      const isAlive = n.is_alive ?? true;
      const themeColor = n.gender === "male" ? "#0ea5e9" : "#f43f5e";
      const borderColor = isAlive ? themeColor : "#94a3b8";

      const getYear = (dateStr?: string) => {
        if (!dateStr) return "";
        try {
          return dateStr.split("-")[0];
        } catch (e) {
          return "";
        }
      };

      const birthYear = getYear(n.birth_date);
      const deathYear = !isAlive ? getYear(n.death_date) : "Sekarang";
      const yearStr = birthYear ? `\n(${birthYear} - ${deathYear})` : "";

      return {
        id: n.id,
        x: positions[n.id].x,
        y: positions[n.id].y,
        label:
          (n.nickname || n.full_name.split(" ")[0]) + (yearStr ? yearStr : ""),
        title: `
          <div class="px-3 py-3 space-y-2 min-w-[160px] bg-white dark:bg-zinc-900">
            <div class="font-black text-sm tracking-tight text-zinc-900 dark:text-white leading-tight">${n.full_name}</div>
            <div class="flex flex-wrap items-center gap-1.5">
              <span class="text-[9px] font-black uppercase tracking-widest px-2 py-0.5 rounded-md ${n.gender === "male" ? "bg-blue-50 text-blue-600" : "bg-rose-50 text-rose-600"}">
                ${getGenderLabel(n.gender)}
              </span>
              ${
                !isAlive
                  ? `
                <span class="text-[9px] font-black uppercase tracking-widest px-2 py-0.5 rounded-md bg-zinc-100 dark:bg-zinc-800 text-zinc-500">
                  Wafat ${deathYear ? `(${deathYear})` : ""}
                </span>
              `
                  : `
                <span class="text-[9px] font-black uppercase tracking-widest px-2 py-0.5 rounded-md bg-emerald-50 text-emerald-600">
                  Hidup
                </span>
              `
              }
            </div>
            ${
              birthYear
                ? `
              <div class="text-[10px] font-bold text-zinc-400 mt-1">
                ${birthYear} — ${isAlive ? "Sekarang" : deathYear || "?"}
              </div>
            `
                : ""
            }
          </div>
        `,
        shape: "circularImage",
        image:
          n.photo_url ||
          `https://ui-avatars.com/api/?name=${encodeURIComponent(n.full_name)}&background=${(isAlive ? themeColor : "#cbd5e1").replace("#", "")}&color=fff&bold=true&size=128`,
        borderWidth: isAlive ? (n.photo_url ? 6 : 4) : 4,
        size: 45,
        opacity: isAlive ? 1 : 0.8,
        color: {
          border: borderColor,
          background: "#ffffff",
          highlight: { border: "#fbbf24", background: "#ffffff" },
        },
        font: {
          color: isAlive ? "#09090b" : "#71717a",
          size: isAlive ? 15 : 13,
          face: "Inter, sans-serif",
          strokeWidth: 5,
          strokeColor: "#ffffff",
          multi: true,
          bold: true,
        },
        fixed: true,
        physics: false,
      };
    });

    Object.entries(familyNodePositions).forEach(([fid, pos]) => {
      visNodes.push({
        id: fid,
        x: pos.x,
        y: pos.y,
        size: 4,
        shape: "dot",
        color: {
          background: "rgba(148,163,184,0.45)",
          border: "rgba(148,163,184,0.45)",
          highlight: "rgba(148,163,184,0.45)",
          hover: "rgba(148,163,184,0.45)",
        },
        label: "",
        physics: false,
        fixed: true,
        font: { size: 1, color: "rgba(0,0,0,0)" },
      });
    });

    const visEdges: any[] = [];

    familyUnits.forEach((f) => {
      if (f.parents.length >= 2) {
        visEdges.push({
          from: f.parents[0],
          to: f.parents[1],
          dashes: true,
          width: 3,
          color: { color: "#fbbf24", highlight: "#fbbf24" },
          smooth: false,
        });
      }

      f.parents.forEach((parentId) => {
        visEdges.push({
          from: parentId,
          to: f.id,
          arrows: "",
          width: 1.5,
          color: { color: "#94a3b8", highlight: "#94a3b8" },
          smooth: {
            enabled: true,
            type: "cubicBezier",
            forceDirection: "vertical",
            roundness: 0.5,
          },
        });
      });

      f.children.forEach((childId) => {
        visEdges.push({
          from: f.id,
          to: childId,
          arrows: "",
          width: 1.5,
          color: { color: "#94a3b8", highlight: "#94a3b8" },
          smooth: {
            enabled: true,
            type: "cubicBezier",
            forceDirection: "vertical",
            roundness: 0.5,
          },
        });
      });
    });

    const data = {
      nodes: new DataSet(visNodes),
      edges: new DataSet(visEdges),
    };

    const options = {
      nodes: {
        shapeProperties: { useBorderWithImage: true },
      },
      edges: {
        smooth: {
          enabled: true,
          type: "cubicBezier",
          forceDirection: "vertical",
          roundness: 0.5,
        },
      },
      layout: {
        improvedLayout: false,
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
      const id = params.nodes[0];

      if (typeof id === "string" && id.startsWith("family-")) {
        network?.unselectAll();
        return;
      }

      selectedPersonId = id;
      network?.focus(id, {
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

    network.once("afterDrawing", () => {
      if (treeData.root_id && positions[treeData.root_id]) {
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
  <div bind:this={container} class="w-full h-full"></div>

  <TreeControls
    bind:searchQuery
    bind:showLegend
    onSearch={handleSearch}
    onReset={resetView}
  />

  <TreeLegend bind:showLegend />

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
