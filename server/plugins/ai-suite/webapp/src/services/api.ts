const DEFAULT_BASE_PATH = `/plugins/com.mattermost.ai-suite/api/v1`;

export type ApiError = {
    message: string;
    status: number;
    details?: Record<string, unknown>;
};

export class APIClient {
    private basePath: string;

    constructor(basePath: string = DEFAULT_BASE_PATH) {
        this.basePath = basePath;
    }

    public get<T>(path: string): Promise<T> {
        return this.request<T>(path, {method: 'GET'});
    }

    public post<T>(path: string, body?: unknown): Promise<T> {
        return this.request<T>(path, {
            method: 'POST',
            body: body ? JSON.stringify(body) : undefined,
            headers: {'Content-Type': 'application/json'},
        });
    }

    private async request<T>(path: string, init: RequestInit): Promise<T> {
        const url = `${this.basePath}${path}`;
        const options: RequestInit = {
            ...init,
            credentials: 'same-origin',
            headers: {
                ...init.headers,
                'X-Requested-With': 'XMLHttpRequest',
                'X-CSRF-Token': (window as any).csrf_token ?? '',
            },
        };

        const response = await fetch(url, options);
        const contentType = response.headers.get('Content-Type') ?? '';
        const isJSON = contentType.includes('application/json');
        const payload = isJSON ? await response.json() : await response.text();

        if (!response.ok) {
            const error: ApiError = {
                status: response.status,
                message: isJSON && payload?.error ? payload.error : response.statusText,
                details: isJSON ? payload : undefined,
            };
            throw error;
        }

        return payload as T;
    }
}

export const apiClient = new APIClient();

