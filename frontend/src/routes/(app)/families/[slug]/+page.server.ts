import type { PageServerLoad } from './$types';
import { apiFetch } from '$lib/api';
import { error } from '@sveltejs/kit';

export const load: PageServerLoad = async ({ params, fetch }) => {
    const { slug } = params;

    try {
        const [family, members] = await Promise.all([
            apiFetch(`/families/${slug}`, { fetch }),
            apiFetch(`/families/${slug}/members`, { fetch })
        ]);

        return {
            family,
            members,
            slug
        };
    } catch (err: any) {
        console.error(`Error loading family [${slug}]:`, err.message);
        
        const msg = err.message?.toLowerCase() || '';
        
        if (msg.includes('not found')) {
            throw error(404, 'Keluarga tidak ditemukan');
        }
        
        if (msg.includes('unauthorized') || msg.includes('login')) {
            throw error(401, 'Silakan login untuk mengakses keluarga ini');
        }
        
        if (msg.includes('forbidden') || msg.includes('insufficient role') || msg.includes('access')) {
            throw error(403, 'Anda tidak memiliki izin untuk mengakses data ini');
        }

        throw error(500, err.message || 'Terjadi kesalahan internal');
    }
};
