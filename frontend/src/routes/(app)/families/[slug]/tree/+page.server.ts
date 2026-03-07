import type { PageServerLoad } from './$types';
import { apiFetch } from '$lib/api';
import { error } from '@sveltejs/kit';

export const load: PageServerLoad = async ({ params, fetch, cookies }) => {
    const { slug } = params;
    const sessionId = cookies.get('session_id');
    const cookieString = sessionId ? `session_id=${sessionId}` : '';

    try {
        const family = await apiFetch(`/families/${slug}`, { fetch, cookie: cookieString });
        const treeDataPromise = apiFetch(`/families/${slug}/tree`, { fetch, cookie: cookieString });

        return {
            family,
            streamed: {
                treeData: treeDataPromise
            },
            slug
        };
    } catch (err: any) {
        console.error(`Error loading tree for family [${slug}]:`, err.message);
        throw error(404, 'Data keluarga tidak ditemukan');
    }
};
