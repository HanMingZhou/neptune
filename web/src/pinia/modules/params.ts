import { getSysParam } from '@/api/sysParams'
import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useParamsStore = defineStore('params', () => {
    const paramsMap = ref<Record<string, any>>({})

    const setParamsMap = (paramsRes: Record<string, any>) => {
        paramsMap.value = { ...paramsMap.value, ...paramsRes }
    }

    const getParams = async (key: string) => {
        if (paramsMap.value[key] && paramsMap.value[key].length) {
            return paramsMap.value[key]
        } else {
            const res = await getSysParam({ key })
            if (res.code === 0) {
                const paramsRes: Record<string, any> = {}
                paramsRes[key] = res.data.value
                setParamsMap(paramsRes)
                return paramsMap.value[key]
            }
        }
    }

    return {
        paramsMap,
        setParamsMap,
        getParams
    }
})
