import type { PageServerLoad } from './$types';
import { INTERNAL_BACKEND_URL } from '$env/static/private';
import { error } from '@sveltejs/kit';

export const load: PageServerLoad = async ({ params, fetch }) => {
    const { slug } = params;

    try {
        const [familyRes, treeRes] = await Promise.all([
            fetch(`${INTERNAL_BACKEND_URL}/public/families/${slug}`),
            fetch(`${INTERNAL_BACKEND_URL}/public/families/${slug}/tree`)
        ]);

        if (!familyRes.ok) {
            throw error(404, 'Keluarga tidak ditemukan atau bersifat privat.');
        }

        const familyData = await familyRes.json();
        const treeData = treeRes.ok ? await treeRes.json() : { data: { nodes: [], edges: [], root_id: 0 } };

        return {
            family: familyData.data,
            treeData: treeData.data || { nodes: [], edges: [], root_id: 0 },
            slug
        };
    } catch (err: any) {
        if (err.status) throw err;
        console.error(`Error loading public tree [${slug}]:`, err.message);
        throw error(404, 'Data pohon silsilah tidak ditemukan.');
    }
};
