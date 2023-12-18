declare namespace Req {
    interface TabModel {
        id: string
        isShow: boolean
        label: string
    }

    interface RunListModel {
        isDisable: boolean
        name: string
        valueType: string
        isMust: string
        testValue: string
        egValue: string
        desc: string
    }

    interface DataListModel {
        isDisable: boolean
        name: string
        testValue: string
    }

    interface AssertListModel {
        type: string
        source: string
        value: string
        assertResult: string
    }

    interface ActionListModel {
        type: string
        value: string
    }

    interface HeaderListModel {
        isDisable: boolean
        name: string
        valueType: string
        isMust: string
        testValue: string
        egValue: string
        desc: string
    }

    interface OtherConfigListModel {
        version: string
        apiId: string
        isParallel: string
        isUseEnvConfig: string
    }

    interface EnvListModel {
        product: string
    }

    interface ShareModel {
        inValueTypeOptions: string[]
    }

    interface RelatedApiListModel {
        dataFile: string
    }

    interface HistoryApiListModel {
        dataFile: string
    }

    interface DepDataListModel {
        value: string
        label: string
    }

    interface DefListModel {
        name: string
        valueType: string
        isMust: string
        egValue: string
        desc: string
    }

    interface ApiDefModel {
        app: string
        module: string
        apiDesc: string
        prefix: string
        path: string
    }

    interface ApiDefSaveModel {
        app: string
        module: string
        apiDesc: string
        prefix: string
        method: string
        path: string
        bodyMode: string
        bodyStr: string
        pathVars: DefListModel[]
        queryVars: DefListModel[]
        bodyVars: DefListModel[]
        headerVars: DefListModel[]
        respVars: DefListModel[]
    }

    interface ApiRunSaveModel {
        app: string
        module: string
        apiDesc: string
        dataDesc: string
        prefix: string
        method: string
        prototype: string
        hostIp: string
        path: string
        product: string
        bodyMode: string
        bodyStr: string
        pathVars: RunListModel[]
        queryVars: RunListModel[]
        bodyVars: RunListModel[]
        headerVars: RunListModel[]
        respVars: RunListModel[]
        actions: ActionListModel[]
        asserts: AssertListModel[]
        preApis: RelatedApiListModel[]
        postApis: RelatedApiListModel[]
        otherConfigs: OtherConfigListModel[]
    }

    interface ApiHistorySaveModel {
        app: string
        module: string
        apiDesc: string
        dataDesc: string
        prefix: string
        method: string
        prototype: string
        host: string
        path: string
        product: string
        fileName: string
        bodyMode: string
        bodyStr: string
        pathVars: RunListModel[]
        queryVars: RunListModel[]
        bodyVars: RunListModel[]
        headerVars: RunListModel[]
        respVars: RunListModel[]
        actions: ActionListModel[]
        asserts: AssertListModel[]
        preApis: RelatedApiListModel[]
        postApis: RelatedApiListModel[]
        otherConfigs: OtherConfigListModel[]
        output: string
        // response: string
        // request: string
        // url: string
        // header: string
        // testResult: string
        // failReason: string
    }

    interface DataRunSaveModel {
        app: string
        module: string
        apiDesc: string
        dataDesc: string
        prefix: string
        method: string
        prototype: string
        hostIp: string
        path: string
        product: string
        bodyMode: string
        pathVars: RunListModel[]
        queryVars: RunListModel[]
        bodyVars: RunListModel[]
        headerVars: RunListModel[]
        respVars: RunListModel[]
        actions: ActionListModel[]
        asserts: AssertListModel[]
        preApis: RelatedApiListModel[]
        postApis: RelatedApiListModel[]
        otherConfigs: OtherConfigListModel[]
    }

    interface ApiDataModel {
        app: string
        module: string
        apiDesc: string
        prefix: string
        path: string
        dataDesc: string
    }

    interface ReqDataRespModel {
        response: string
        request: string
        url: string
        header: string
        testResult: string
        failReason: string
        output: string
    }

    interface AdvancedModel {
        version: string,
        apiId: string,
        isParallel: string,
        isUseEnvConfig: string
    }

    interface ReqSceneRespModel {
        lastDataFile: string
        testResult: string
        failReason: string
    }

    interface RequestModel {
        n: number
        c: number
        timeout: number
        method: string
        url: string
        headers: ParamModel
        body: string
    }

    interface SceneModel {
        product: string
        name: string
        dataList: RelatedApiListModel[]
        type: string,
        runNum: number
    }

    interface SceneHistoryModel {
        product: string
        name: string
        dataList: RelatedApiListModel[]
        type: string
        runNum: number
        result: string
        lastFile: string
        failReason: string
    }

    interface ResponseModel {
        error?: string
        status: string
        statusCode: number
        contentLength: number
        proto: string
        headers: ParamModel
        cookies: CookieModel[]
        body: string
        code: string
        duration: ResponseDuration
    }

    interface ResponseDuration {
        dns: string
        conn: string
        req: string
        res: string
        delay: string
        finish: string
    }

    interface ParamModel {
        [key: string]: string[]
    }

    interface HeaderModel {
        key: string
        value: string
    }

    interface CookieModel {
        Name: string
        Value: string

        Path: string    // optional
        Domain: string    // optional
        Expires: string // optional
        RawExpires: string    // for reading cookies only

        // MaxAge=0 means no 'Max-Age' attribute specified.
        // MaxAge<0 means delete cookie now, equivalently 'Max-Age: 0'
        // MaxAge>0 means Max-Age attribute present and given in seconds
        MaxAge: number
        Secure: boolean
        HttpOnly: boolean
        Raw: string
        Unparsed: string[] // Raw text of unparsed attribute-value pairs
    }
}