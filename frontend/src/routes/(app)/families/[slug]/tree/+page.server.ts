import type { PageServerLoad } from './$types';
import { apiFetch } from '$lib/api';
import { error } from '@sveltejs/kit';

export const load: PageServerLoad = async ({ params, fetch }) => {
    const { slug } = params;

    try {
        const [family, treeData] = await Promise.all([
            apiFetch(`/families/${slug}`, { fetch }),
            apiFetch(`/families/${slug}/tree`, { fetch })
        ]);

        return {
            family,
            treeData,
            slug
        };
    } catch (err: any) {
        console.error(`Error loading tree for family [${slug}]:`, err.message);
        throw error(404, 'Data pohon silsilah tidak ditemukan');
    }
};
