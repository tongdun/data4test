const message: i18nModel = {
    tabTitle: {
        newTab: "新标签",
        contentTab: "请求内容",
    },
    buttonTitle: {
        sceneSave: "场景保存",
        dataSave: "数据保存",
        apiSave: "接口保存",
        dataSend: "运行",
    },
    general: {
        noData: '暂无数据',
        noFilterData: '暂无筛选结果'
    },
    list: {
        isDisable: false,
        name: '参数名',
        valueType: '参数类型',
        isMust: '必填',
        testValue: '测试值',
        egValue: '示例值',
        desc: '描述',
        key: '键',
        value: '值'
    },

    actionList: {
        type: '动作类型',
        value: '值'
    },

    assertList: {
        type: '断言类型',
        source: '源数据定义',
        value: '值',
        assertResult: '断言结果',
    },
    otherConfigList: {
        version: "数据版本",
        apiId: "接口ID",
        isParallel: "是否并发",
        isUseEnvConfig: "是否使用公共环境"
    },
    envList: {
        product: '关联产品',
    },
    relatedApiList: {
        dataFile: '数据文件',
    },
    app: {
        modal: {
            title: '操作确认',
            content: '已退出，是否关闭本页面？',
            okText: '确定',
            cancelText: '取消'
        }
    },
    api: {
        appTips: '所属应用',
        moduleTips: '所属模块',
        defDescTips: '接口描述',
        dataDescTips: '数据描述',
        prefixTips: '公共前缀',
        methodTips: '请求方法',
        urlTips: '路由地址',
        prototypeTips: '请求协议',
        hostIpTips: '请求IP',
        envTips: '选择环境'
    },

    scene: {
        nameTips: '场景名称',
        typeTips: '场景类型',
        runNumTips: '执行次数'
    },

    advanced: {
        versionTips: "数据版本",
        apiIdTips: "接口ID",
        isParallelTips: "是否并发",
        isUseEnvConfigTips: "是否使用公共环境"
    },

    title: {
        reqDataTitle: "请求数据",
        respDataTitle: "响应数据",
        respResultTitle: "测试结果",
        advancedTitle: "高级选项",
        outputTitle: "输出参数"
    },

    // 报告相关
    req: {
        report: {
            title: "性能报告",
            summary: {
                total: "总耗时",
                slowest: "最慢",
                fastest: "最快",
                average: "平均",
                rps: "每秒请求数"
            },
            table: {
                type: "类型",
                average: "平均",
                slowest: "最慢",
                fastest: "最快",
                dns: "DNS查询",
                conn: "连接",
                req: "请求",
                delay: "延迟",
                res: "响应"
            },
            chart: {
                title: "耗时分布",
                serieName: "耗时"
            }
        },
        pre: {
            title: "预处理",
            code: "代码",
            report: "报告",
            body: {
                pretty: "格式化",
                raw: "原始",
                language: "语言"
            },
            listTable: {
                duration: "耗时",
                hasError: "是否有错"
            },
            cookie: {
                name: "名称",
                value: "值",
                domain: "域名",
                path: "路径",
                expires: "过期时间"
            }
        }
    },

    // 表头相关
    headerList: {
        name: "名称",
        valueType: "值类型",
        isDisable: "是否禁用",
        assertResult: "断言结果"
    },

    // go-admin 兼容
    goman: {
        general: {
            noFilterData: "暂无筛选结果"
        }
    }
}

export default message
