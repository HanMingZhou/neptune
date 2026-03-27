export interface NavItem {
    key: string;
    icon?: string;
    path?: string;
    children?: NavItem[];
}

export interface NavGroup {
    title: string;
    items: NavItem[];
}

export const NAVIGATION_GROUPS: NavGroup[] = [
    {
        title: 'compute',
        items: [
            { key: 'dashboard', icon: 'dashboard', path: '/' },
            { key: 'notebooks', icon: 'terminal', path: '/notebooks' },
            { key: 'training', icon: 'model_training', path: '/training' },
            { key: 'inference', icon: 'rocket_launch', path: '/inference' },
        ]
    },
    {
        title: 'resources',
        items: [
            { key: 'sshkeys', icon: 'vpn_key', path: '/sshkeys' },
            { key: 'images', icon: 'photo_library', path: '/images' },
            { key: 'storage', icon: 'storage', path: '/storage' },
            {
                key: 'order',
                icon: 'payments',
                children: [
                    { key: 'transactions', icon: 'receipt', path: '/order/transactions' },
                    { key: 'usage', icon: 'data_usage', path: '/order/usage' },
                    { key: 'invoice', icon: 'description', path: '/order/invoice' },
                ]
            },
        ]
    },
    {
        title: 'management',
        items: [
            {
                key: 'account',
                icon: 'person',
                children: [
                    { key: 'security', icon: 'shield', path: '/account/security' },
                    { key: 'accessRecords', icon: 'history', path: '/account/records' },
                    { key: 'settings', icon: 'settings', path: '/account/settings' },
                ]
            }
        ]
    },
    {
        title: 'admin',
        items: [
            { key: 'products', icon: 'inventory_2', path: '/admin/products' },
            {
                key: 'roles',
                icon: 'admin_panel_settings',
                children: [
                    { key: 'roles', icon: 'badge', path: '/admin/roles' },
                    { key: 'menus', icon: 'menu_open', path: '/admin/menus' },
                ]
            },
            { key: 'apis', icon: 'api', path: '/admin/apis' },
            { key: 'users', icon: 'group', path: '/admin/users' },
            { key: 'operations', icon: 'receipt_long', path: '/admin/operations' },
            { key: 'clusterManage', icon: 'cloud', path: '/admin/clusters' },
            { key: 'nodeManage', icon: 'hub', path: '/admin/nodes' },
        ]
    }
];

export const STATUS_COLORS: Record<string, string> = {
    Running: 'text-emerald-500 bg-emerald-500/10',
    Succeeded: 'text-emerald-500 bg-emerald-500/10',
    Stopped: 'text-slate-400 bg-slate-400/10',
    '已停止': 'text-slate-400 bg-slate-400/10',
    Error: 'text-red-500 bg-red-500/10',
    Failed: 'text-red-500 bg-red-500/10',
    Starting: 'text-blue-500 bg-blue-500/10',
    Queued: 'text-amber-500 bg-amber-500/10',
    Deploying: 'text-primary bg-primary/10',
};
