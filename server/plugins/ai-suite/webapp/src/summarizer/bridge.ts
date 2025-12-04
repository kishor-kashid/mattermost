import type {SummaryType} from '../services/summarizerApi';

export type SummaryTarget = {
    type: SummaryType;
    channelId?: string;
    rootPostId?: string;
    postId?: string;
    timeRange?: string;
    since?: number;
    until?: number;
    label?: string;
};

type Listener = (target: SummaryTarget | null) => void;

export class SummaryBridge {
    private listeners: Set<Listener> = new Set();
    private target: SummaryTarget | null = null;

    public setTarget(target: SummaryTarget | null) {
        this.target = target;
        this.listeners.forEach((listener) => listener(target));
    }

    public subscribe(listener: Listener) {
        this.listeners.add(listener);
        listener(this.target);
        return () => {
            this.listeners.delete(listener);
        };
    }

    public getTarget() {
        return this.target;
    }
}


