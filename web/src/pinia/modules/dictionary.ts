import { findSysDictionary } from '@/api/sysDictionary'
import { getDictionaryTreeListByType } from '@/api/sysDictionaryDetail'

import { defineStore } from 'pinia'
import { ref } from 'vue'

export interface DictItem {
    label: string
    value: string | number
    extend?: string
    children?: DictItem[]
}

export const useDictionaryStore = defineStore('dictionary', () => {
    const dictionaryMap = ref<Record<string, DictItem[]>>({})

    const setDictionaryMap = (dictionaryRes: Record<string, DictItem[]>) => {
        dictionaryMap.value = { ...dictionaryMap.value, ...dictionaryRes }
    }

    const filterTreeByDepth = (items: DictItem[], currentDepth: number, targetDepth: number): DictItem[] => {
        if (targetDepth === 0) {
            return items
        }

        if (currentDepth >= targetDepth) {
            return items.map((item) => ({
                label: item.label,
                value: item.value,
                extend: item.extend
            }))
        }

        return items.map((item) => ({
            label: item.label,
            value: item.value,
            extend: item.extend,
            children: item.children
                ? filterTreeByDepth(item.children, currentDepth + 1, targetDepth)
                : undefined
        }))
    }

    const flattenTree = (items: DictItem[]): DictItem[] => {
        const result: DictItem[] = []

        const traverse = (nodes: DictItem[]) => {
            nodes.forEach((item) => {
                result.push({
                    label: item.label,
                    value: item.value,
                    extend: item.extend
                })

                if (item.children && item.children.length > 0) {
                    traverse(item.children)
                }
            })
        }

        traverse(items)
        return result
    }

    const normalizeTreeData = (items: DictItem[]): DictItem[] => {
        return items.map((item) => ({
            label: item.label,
            value: item.value,
            extend: item.extend,
            children:
                item.children && item.children.length > 0
                    ? normalizeTreeData(item.children)
                    : undefined
        }))
    }

    const findNodeByValue = (
        items: DictItem[],
        targetValue: string | number,
        currentDepth: number = 1,
        maxDepth: number = 0
    ): DictItem[] | null => {
        for (const item of items) {
            if (item.value === targetValue) {
                if (maxDepth === 0) {
                    return item.children ? normalizeTreeData(item.children) : []
                }
                if (item.children && item.children.length > 0) {
                    return filterTreeByDepth(item.children, 1, maxDepth)
                }
                return []
            }

            if (
                item.children &&
                item.children.length > 0 &&
                (maxDepth === 0 || currentDepth < maxDepth)
            ) {
                const result = findNodeByValue(
                    item.children,
                    targetValue,
                    currentDepth + 1,
                    maxDepth
                )
                if (result !== null) {
                    return result
                }
            }
        }
        return null
    }

    const getDictionary = async (type: string, depth: number = 0, value: string | number | null = null) => {
        if (value !== null) {
            const cacheKey = `${type}_value_${value}_depth_${depth}`

            if (
                dictionaryMap.value[cacheKey] &&
                dictionaryMap.value[cacheKey].length
            ) {
                return dictionaryMap.value[cacheKey]
            }

            try {
                const treeRes = await getDictionaryTreeListByType({ type })
                if (
                    treeRes.code === 0 &&
                    treeRes.data &&
                    treeRes.data.list &&
                    treeRes.data.list.length > 0
                ) {
                    const targetNodeChildren = findNodeByValue(
                        treeRes.data.list,
                        value,
                        1,
                        depth
                    )

                    if (targetNodeChildren !== null) {
                        let resultData: DictItem[]
                        if (depth === 0) {
                            resultData = targetNodeChildren
                        } else {
                            resultData = flattenTree(targetNodeChildren)
                        }

                        const dictionaryRes: Record<string, DictItem[]> = {}
                        dictionaryRes[cacheKey] = resultData
                        setDictionaryMap(dictionaryRes)
                        return dictionaryMap.value[cacheKey]
                    } else {
                        return []
                    }
                }
            } catch (error) {
                console.error('根据value获取字典数据失败:', error)
                return []
            }
        }

        const cacheKey = depth === 0 ? `${type}_tree` : `${type}_depth_${depth}`

        if (dictionaryMap.value[cacheKey] && dictionaryMap.value[cacheKey].length) {
            return dictionaryMap.value[cacheKey]
        } else {
            try {
                const treeRes = await getDictionaryTreeListByType({ type })
                if (
                    treeRes.code === 0 &&
                    treeRes.data &&
                    treeRes.data.list &&
                    treeRes.data.list.length > 0
                ) {
                    const treeData = treeRes.data.list

                    let resultData: DictItem[]
                    if (depth === 0) {
                        resultData = normalizeTreeData(treeData)
                    } else {
                        const filteredData = filterTreeByDepth(treeData, 1, depth)
                        resultData = flattenTree(filteredData)
                    }

                    const dictionaryRes: Record<string, DictItem[]> = {}
                    dictionaryRes[cacheKey] = resultData
                    setDictionaryMap(dictionaryRes)
                    return dictionaryMap.value[cacheKey]
                } else {
                    const res = await findSysDictionary({ type })
                    if (res.code === 0) {
                        const dictionaryRes: Record<string, DictItem[]> = {}
                        const dict: DictItem[] = []
                        res.data.resysDictionary.sysDictionaryDetails &&
                            res.data.resysDictionary.sysDictionaryDetails.forEach((item: any) => {
                                dict.push({
                                    label: item.label,
                                    value: item.value,
                                    extend: item.extend
                                })
                            })
                        dictionaryRes[cacheKey] = dict
                        setDictionaryMap(dictionaryRes)
                        return dictionaryMap.value[cacheKey]
                    }
                }
            } catch (error) {
                console.error('获取字典数据失败:', error)
                const res = await findSysDictionary({ type })
                if (res.code === 0) {
                    const dictionaryRes: Record<string, DictItem[]> = {}
                    const dict: DictItem[] = []
                    res.data.resysDictionary.sysDictionaryDetails &&
                        res.data.resysDictionary.sysDictionaryDetails.forEach((item: any) => {
                            dict.push({
                                label: item.label,
                                value: item.value,
                                extend: item.extend
                            })
                        })
                    dictionaryRes[cacheKey] = dict
                    setDictionaryMap(dictionaryRes)
                    return dictionaryMap.value[cacheKey]
                }
            }
        }
    }

    return {
        dictionaryMap,
        setDictionaryMap,
        getDictionary
    }
})
