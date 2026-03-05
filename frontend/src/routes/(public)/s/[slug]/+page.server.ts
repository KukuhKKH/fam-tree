import type { PageServerLoad } from './$types';
import { env } from '$env/dynamic/private';
const INTERNAL_BACKEND_URL = env.INTERNAL_BACKEND_URL;
import { error } from '@sveltejs/kit';

export const load: PageServerLoad = async ({ params, fetch }) => {
    const { slug } = params;

    try {
        const [familyRes, personsRes, treeRes] = await Promise.all([
            fetch(`${INTERNAL_BACKEND_URL}/public/families/${slug}`),
            fetch(`${INTERNAL_BACKEND_URL}/public/families/${slug}/persons`),
            fetch(`${INTERNAL_BACKEND_URL}/public/families/${slug}/tree`)
        ]);

        if (!familyRes.ok) {
            const err = await familyRes.json().catch(() => ({}));
            const msg = err.messages?.[0] || 'Keluarga tidak ditemukan';
            if (msg.includes('not publicly')) {
                throw error(403, 'Keluarga ini bersifat privat dan tidak dapat diakses publik.');
            }
            throw error(404, msg);
        }

        const familyData = await familyRes.json();
        const personsData = personsRes.ok ? await personsRes.json() : { data: [] };
        const treeData = treeRes.ok ? await treeRes.json() : { data: { nodes: [], edges: [], root_id: 0 } };

        return {
            family: familyData.data,
            persons: personsData.data || [],
            treeData: treeData.data || { nodes: [], edges: [], root_id: 0 },
            slug
        };
    } catch (err: any) {
        if (err.status) throw err; // Re-throw SvelteKit errors
        console.error(`Error loading public family [${slug}]:`, err.message);
        throw error(404, 'Keluarga tidak ditemukan atau bersifat privat.');
    }
};
