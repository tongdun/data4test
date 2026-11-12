const message: i18nModel = {
    tabTitle: {
        newTab: "New Tab",
        contentTab: "Request Content",
    },
    buttonTitle: {
        sceneSave: "Save Playbook",
        dataSave: "Save Data",
        apiSave: "Save API",
        dataSend: "Run",
    },
    general: {
        noData: 'No Data',
        noFilterData: 'No Results'
    },
    list: {
        isDisable: false,
        name: 'Name',
        valueType: 'Type',
        isMust: 'Required',
        testValue: 'Test Value',
        egValue: 'Example',
        desc: 'Desc',
        key: 'Key',
        value: 'Value'
    },

    actionList: {
        type: 'Action Type',
        value: 'Value'
    },

    assertList: {
        type: 'Assert Type',
        source: 'Source',
        value: 'Value',
        assertResult: 'Result',
    },
    otherConfigList: {
        version: "Data Version",
        apiId: "API ID",
        isParallel: "Concurrent",
        isUseEnvConfig: "Use Env Config"
    },
    envList: {
        product: 'Product',
    },
    relatedApiList: {
        dataFile: 'Data File',
    },
    app: {
        modal: {
            title: 'Confirm',
            content: 'Already logged out. Close this page?',
            okText: 'OK',
            cancelText: 'Cancel'
        }
    },
    api: {
        appTips: 'Application',
        moduleTips: 'Module',
        defDescTips: 'API Desc',
        dataDescTips: 'Data Desc',
        prefixTips: 'Prefix',
        methodTips: 'Method',
        urlTips: 'URL',
        prototypeTips: 'Protocol',
        hostIpTips: 'Host IP',
        envTips: 'Select Environment'
    },

    scene: {
        nameTips: 'Playbook Name',
        typeTips: 'Playbook Type',
        runNumTips: 'Run Count'
    },

    advanced: {
        versionTips: "Data Version",
        apiIdTips: "API ID",
        isParallelTips: "Concurrent",
        isUseEnvConfigTips: "Use Env Config"
    },

    title: {
        reqDataTitle: "Request Data",
        respDataTitle: "Response Data",
        respResultTitle: "Test Result",
        advancedTitle: "Advanced",
        outputTitle: "Output"
    },

    // Report related
    req: {
        report: {
            title: "Performance Report",
            summary: {
                total: "Total",
                slowest: "Slowest",
                fastest: "Fastest",
                average: "Average",
                rps: "Requests/s"
            },
            table: {
                type: "Type",
                average: "Average",
                slowest: "Slowest",
                fastest: "Fastest",
                dns: "DNS Lookup",
                conn: "Connect",
                req: "Request",
                delay: "Delay",
                res: "Response"
            },
            chart: {
                title: "Duration Distribution",
                serieName: "Duration"
            }
        },
        pre: {
            title: "Preprocess",
            code: "Code",
            report: "Report",
            body: {
                pretty: "Pretty",
                raw: "Raw",
                language: "Language"
            },
            listTable: {
                duration: "Duration",
                hasError: "Has Error"
            },
            cookie: {
                name: "Name",
                value: "Value",
                domain: "Domain",
                path: "Path",
                expires: "Expires"
            }
        }
    },

    // Header related
    headerList: {
        name: "Name",
        valueType: "Type",
        isDisable: "Disabled",
        assertResult: "Result"
    },

    // go-admin compatibility
    goman: {
        general: {
            noFilterData: "No Results"
        }
    }
}

export default message
