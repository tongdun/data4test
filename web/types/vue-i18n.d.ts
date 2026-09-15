/// <reference types='vue-i18n' />

declare interface i18nModel {
    tabTitle: {
        newTab: string
        contentTab: string
    },
    buttonTitle: {
        sceneSave: string
        dataSave: string
        apiSave: string
        dataSend: string
    }
    general: {
        noData: string
        noFilterData: string
    }
    list: {
        isDisable: boolean
        name: string
        valueType: string
        isMust: string
        testValue: string
        egValue: string
        desc: string
        key: string
        value: string
    }

    actionList: {
        type: string
        value: string
    }

    assertList: {
        type: string
        source: string
        value: string
        assertResult: string
    }
    otherConfigList: {
        version: string
        apiId: string
        isParallel: string
        isUseEnvConfig: string
    }
    envList: {
        product: string
    },
    relatedApiList: {
        dataFile: string
    },
    app: {
        modal: {
            title: string
            content: string
            okText: string
            cancelText: string
        }
    }
    api: {
        appTips: string
        moduleTips: string
        defDescTips: string
        dataDescTips: string
        prefixTips: string
        methodTips: string
        urlTips: string
        prototypeTips: string
        hostIpTips: string
        envTips: string
    }

    advanced: {
        versionTips: string,
        apiIdTips: string,
        isParallelTips: string,
        isUseEnvConfigTips: string
    }

    scene: {
        nameTips: string
        typeTips: string
        runNumTips: string
    }
    title: {
        reqDataTitle: string
        respDataTitle: string
        respResultTitle: string
        advancedTitle: string
        outputTitle: string
    }
    req: {
        report: {
            title: string
            summary: {
                total: string
                slowest: string
                fastest: string
                average: string
                rps: string
            }
            table: {
                type: string
                average: string
                slowest: string
                fastest: string
                dns: string
                conn: string
                req: string
                delay: string
                res: string
            }
            chart: {
                title: string
                serieName: string
            }
        }
        pre: {
            title: string
            code: string
            report: string
            body: {
                pretty: string
                raw: string
                language: string
            }
            listTable: {
                duration: string
                hasError: string
            }
            cookie: {
                name: string
                value: string
                domain: string
                path: string
                expires: string
            }
        }
    }
    headerList: {
        name: string
        valueType: string
        isDisable: string
        assertResult: string
    }
    goman: {
        general: {
            noFilterData: string
        }
    }
}
