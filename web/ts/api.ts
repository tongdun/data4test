import _ from 'lodash'
import Axios, { AxiosRequestConfig } from 'axios'

interface APIResult<T> {
    code: number
    msg: string
    data?: T
}

interface AppListResult<T> {
    code: number
    msg: string
    data: string[]
}

interface APIStatic {
    /**
     * get请求
     */
    get<T>(url: string): Promise<APIResult<T>>
    /**
     * post请求
     */
    post<T>(url: string, data?: URLSearchParams): Promise<APIResult<T>>
}

interface AppStatic {
    /**
     * get请求
     */
    get<T>(url: string): Promise<AppListResult<T>>
    /**
     * post请求
     */
    post<T>(url: string, data?: URLSearchParams): Promise<AppListResult<T>>
}

async function request<T>(method: 'get' | 'post', url: string, data?: URLSearchParams) {
    let result: APIResult<T> = { code: -1, msg: '' }

    let conf: AxiosRequestConfig = {
        method: method,
        // url: '/admin' + url
        url: url
    }
    if (data) {
        conf.data = data
    }

    try {
        let res = await Axios.request(conf)
        result.code = res.data.code
        result.msg = res.data.msg
        result.data = res.data.data
    } catch (err) {
        result.msg = err.message
    }

    return result
}

async function defRequest<T>(method: 'get' | 'post', url: string, data: string) {
    let result: APIResult<T> = { code: -1, msg: '' }

    let conf: AxiosRequestConfig = {
        method: method,
        url: '/api' + url
    }
    if (data) {
        conf.data = data
    }

    try {
        let res = await Axios.request(conf)
        result.code = res.data.code
        result.msg = res.data.msg
        result.data = res.data.data
    } catch (err) {
        result.msg = err.message
    }

    return result
}

const API: APIStatic = {
    get: <T>(url: string) => {
        return request<T>('get', url)
    },
    post: <T>(url: string, data?: URLSearchParams) => {
        return request<T>('post', url, data)
    }
}


export default API
