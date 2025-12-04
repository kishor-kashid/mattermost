import {useCallback, useEffect, useMemo, useRef, useState} from 'react';

import {fetchSummary, SummarizeRequest, SummaryResponse} from '../services/summarizerApi';

type Options = {
    request: SummarizeRequest | null;
    enabled?: boolean;
};

export const useSummary = ({request, enabled = true}: Options) => {
    const cacheRef = useRef<Map<string, SummaryResponse>>(new Map());
    const [summary, setSummary] = useState<SummaryResponse | null>(null);
    const [isLoading, setIsLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);

    const cacheKey = useMemo(() => {
        if (!request) {
            return null;
        }
        const clone = {...request};
        delete clone.force;
        return JSON.stringify(clone);
    }, [request]);

    const load = useCallback(
        async (force?: boolean) => {
            if (!request || !cacheKey) {
                return;
            }

            if (!force && cacheRef.current.has(cacheKey)) {
                setSummary(cacheRef.current.get(cacheKey) ?? null);
                setError(null);
                return;
            }

            setIsLoading(true);
            setError(null);
            try {
                const response = await fetchSummary({...request, force});
                cacheRef.current.set(cacheKey, response);
                setSummary(response);
            } catch (err: unknown) {
                let message = 'Unable to load summary';
                if (err instanceof Error) {
                    message = err.message;
                } else if (err && typeof err === 'object' && 'message' in err) {
                    message = String((err as Record<string, unknown>).message);
                }
                setError(message);
            } finally {
                setIsLoading(false);
            }
        },
        [cacheKey, request],
    );

    useEffect(() => {
        if (!enabled || !request || !cacheKey) {
            return;
        }

        setSummary(null);
        setError(null);
        load(false);
    }, [cacheKey, enabled, load, request]);

    const refresh = useCallback(() => load(true), [load]);

    return {
        summary,
        isLoading,
        error,
        refresh,
    };
};


