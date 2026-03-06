import { env } from '$env/dynamic/public';

interface FetchOptions extends RequestInit {
	params?: Record<string, string>;
	fetch?: typeof fetch;
	cookie?: string;
}

export async function apiFetch<T = any>(path: string, options: FetchOptions = {}): Promise<T> {
	const baseUrl = env.PUBLIC_BACKEND_URL || 'http://localhost:8080';
	
	let url = `${baseUrl}${path.startsWith('/') ? '' : '/'}${path}`;
	
	if (options.params) {
		const searchParams = new URLSearchParams(options.params);
		url += `?${searchParams.toString()}`;
	}

	// Use provided fetch (e.g. from SvelteKit load) or global fetch
	const fetchFn = options.fetch || fetch;

	const headers = new Headers(options.headers as HeadersInit);
	if (!headers.has('Content-Type')) {
		headers.set('Content-Type', 'application/json');
	}

	if (options.cookie) {
		headers.set('Cookie', options.cookie);
	}

	const response = await fetchFn(url, {
		...options,
		headers: headers,
		credentials: 'include',
	});

	if (!response.ok) {
		const errorData = await response.json().catch(() => ({}));
		// Try to find error message in various common fields
		const msg = errorData.messages?.[0] || errorData.error || errorData.message || response.statusText;
		throw new Error(msg);
	}

	const result = await response.json();
	// Handle cases where the backend doesn't wrap in "data" or already returns it
	return (result && typeof result === 'object' && 'data' in result) ? result.data : result;
}
