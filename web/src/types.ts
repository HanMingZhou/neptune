export enum InstanceStatus {
    RUNNING = 'Running',
    STOPPED = 'Stopped',
    STARTING = 'Starting',
    ERROR = 'Error',
    DELETING = 'Deleting'
}

export interface NotebookInstance {
    id: string;
    name: string;
    status: InstanceStatus;
    config: string;
    gpuResource: string;
    createdAt: string;
    region: string;
}

export interface TrainingJob {
    id: string;
    name: string;
    status: 'Running' | 'Succeeded' | 'Failed' | 'Queued';
    framework: string;
    gpuCount: number;
    gpuType: string;
    duration: string;
}

export interface ServiceMetric {
    qps: number;
    latency: number;
    replicas: string;
}

export interface InferenceService {
    id: string;
    name: string;
    status: 'Running' | 'Deploying' | 'Error';
    version: string;
    metrics: ServiceMetric;
}
