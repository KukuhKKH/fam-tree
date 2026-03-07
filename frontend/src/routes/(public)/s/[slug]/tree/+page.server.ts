import type { PageServerLoad } from './$types';
import { env } from '$env/dynamic/private';
const INTERNAL_BACKEND_URL = env.INTERNAL_BACKEND_URL;
import { error } from '@sveltejs/kit';

export const load: PageServerLoad = async ({ params, fetch }) => {
    const { slug } = params;

    try {
        const familyRes = await fetch(`${INTERNAL_BACKEND_URL}/public/families/${slug}`);
        if (!familyRes.ok) {
            throw error(404, 'Keluarga tidak ditemukan atau bersifat privat.');
        }

        const familyData = await familyRes.json();
        const treeDataPromise = fetch(`${INTERNAL_BACKEND_URL}/public/families/${slug}/tree`)
            .then(res => res.ok ? res.json() : { data: { nodes: [], edges: [], root_id: 0 } })
            .then(json => json.data || { nodes: [], edges: [], root_id: 0 });

        return {
            family: familyData.data,
            streamed: {
                treeData: treeDataPromise
            },
            slug
        };
    } catch (err: any) {
        if (err.status) throw err;
        console.error(`Error loading public tree [${slug}]:`, err.message);
        throw error(404, 'Data pohon silsilah tidak ditemukan.');
    }
};
