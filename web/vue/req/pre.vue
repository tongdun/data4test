<template>
    <div>
        <Card v-show="isAdvanced" :bordered="false" :padding="0" style="width:100%;">
            <Table stripe border highlight-row :columns="resColumns" :data="pageData" @on-current-change="onResChange" :no-data-text="$t('general.noData')" :no-filtered-data-text="$t('general.noFilterData')"></Table>
            <div style="margin-top:15px;text-align:right">
                <Page class="myPage" :total="resList.length" :current="1" :page-size="pageSize" :page-size-opts="[5,10,20,30,40,50]" @on-change="onPageChange" @on-page-size-change="onPageSizeChange" placement="top" show-total show-sizer>
                </Page>
            </div>
            <div>&nbsp;</div>
        </Card>
        <Card style="width:100%;">
            <div slot="title">{{$t('req.pre.title')}}</div>
            <div slot="extra" v-html="status"></div>
            <div>
                <Tabs v-show="res" :animated="false" v-model="resMode">
                    <Tab-pane label="Body">
                        <RadioGroup v-model="bodyMode" type="button" size="large">
                            <Radio label="pretty">{{$t('req.pre.body.pretty')}}</Radio>
                            <Radio label="raw">{{$t('req.pre.body.raw')}}</Radio>
                        </RadioGroup>
                        <span>&nbsp;&nbsp;&nbsp;&nbsp;{{$t('req.pre.body.language')}}：{{ bodyLanguage }}</span>
                        <pre v-show="bodyMode==='raw'">{{ bodyRaw }}</pre>
                        <pre v-show="bodyMode==='pretty'" v-html="body"></pre>
                    </Tab-pane>
                    <Tab-pane label="Headers">
                        <Table stripe border :columns="headerColumns" :data="headers" :no-data-text="$t('general.noData')" :no-filtered-data-text="$t('general.noFilterData')"></Table>
                    </Tab-pane>
                    <Tab-pane label="Cookies">
                        <Table stripe border :columns="cookieColumns" :data="cookies" :no-data-text="$t('general.noData')" :no-filtered-data-text="$t('general.noFilterData')"></Table>
                    </Tab-pane>
                    <Tab-pane :label="$t('req.pre.code')">
                        <pre>{{ code }}</pre>
                    </Tab-pane>
                    <Tab-pane :label="$t('req.pre.report')">
                        <Table stripe border :columns="reportDurationColumns" :data="durations" :no-data-text="$t('general.noData')" :no-filtered-data-text="$t('general.noFilterData')"></Table>
                    </Tab-pane>
                </Tabs>
            </div>
        </Card>
    </div>
</template>

<style>
.hljs-ul {
    list-style: decimal;
    background-color: #fff;
    margin: 0px 0px 0 40px !important;
    padding: 0px;
}

.hljs-ul li {
    list-style: decimal-leading-zero;
    border-left: 1px solid #ddd !important;
    background: #fff;
    padding: 5px !important;
    margin: 0 !important;
    line-height: 14px;
    word-break: break-all;
    word-wrap: break-word;
}

.hljs-ul li:nth-of-type(even) {
    background-color: #f8f8f9;
    color: inherit;
}
</style>

<script lang="ts">
import { Vue, Component, Prop, Watch } from 'vue-property-decorator'
import { CreateElement } from 'vue/types/vue'
import _ from 'lodash'
import hljs from 'highlight.js'

@Component
export default class Pre extends Vue {
    @Prop() isAdvanced: boolean
    @Prop() resList: Req.ResponseModel[]

    res: Req.ResponseModel | null = null
    pageData: Req.ResponseModel[] = []
    pageSize = 5
    pageCount = 0

    resColumns = <any[]>[
        {
            type: 'index',
            width: 60,
            align: 'center'
        },
        {
            title: '',
            key: 'error'
        },
        {
            title: '',
            key: 'contentLength',
            sortable: true
        },
        {
            title: '',
            render: (h: CreateElement, p: any) => {
                return h('span', p.row.duration.finish)
            }
        }
    ]
    headerColumns = <any[]>[
        {
            title: '',
            key: 'key'
        },
        {
            title: '',
            key: 'value'
        }
    ]
    cookieColumns = <any[]>[
        {
            title: '',
            key: 'Name'
        },
        {
            title: '',
            key: 'Value'
        },
        {
            title: '',
            key: 'Path'
        },
        {
            title: '',
            key: 'Domain'
        },
        {
            title: '',
            key: 'Expires'
        }
    ]
    reportDurationColumns = <any[]>[
        {
            title: '',
            key: 'dns'
        },
        {
            title: '',
            key: 'conn'
        },
        {
            title: '',
            key: 'req'
        },
        {
            title: '',
            key: 'delay'
        },
        {
            title: '',
            key: 'res'
        },
        {
            title: '',
            key: 'finish'
        }
    ]

    bodyMode: 'pretty' | 'raw' = 'pretty'
    bodyLanguage: string = ''
    isPretty: boolean = true
    isRaw: boolean = false

    resMode: number = 0

    created() {
        this.onLocale()
    }

    @Watch('$i18n.locale')
    onLocale() {
        this.resColumns[1].title = this.$t('req.pre.listTable.hasError')
        this.resColumns[2].title = this.$t(
            'req.pre.listTable.contentLength'
        )
        this.resColumns[3].title = this.$t('req.pre.listTable.duration')

        this.headerColumns[0].title = this.$t('list.key')
        this.headerColumns[1].title = this.$t('list.value')

        this.cookieColumns[0].title = this.$t('req.pre.cookie.name')
        this.cookieColumns[1].title = this.$t('req.pre.cookie.value')
        this.cookieColumns[2].title = this.$t('req.pre.cookie.path')
        this.cookieColumns[3].title = this.$t('req.pre.cookie.domain')
        this.cookieColumns[4].title = this.$t('req.pre.cookie.expires')

        this.reportDurationColumns[0].title = this.$t(
            'req.report.table.dns'
        )
        this.reportDurationColumns[1].title = this.$t(
            'req.report.table.conn'
        )
        this.reportDurationColumns[2].title = this.$t(
            'req.report.table.req'
        )
        this.reportDurationColumns[3].title = this.$t(
            'req.report.table.delay'
        )
        this.reportDurationColumns[4].title = this.$t(
            'req.report.table.res'
        )
        this.reportDurationColumns[5].title = this.$t(
            'req.report.table.finish'
        )
    }

    mounted() {
        this.onResListChange()
    }

    onResChange(res: Req.ResponseModel, oldRes: Req.ResponseModel) {
        this.res = res
    }

    @Watch('resList')
    onResListChange() {
        if (this.resList.length == 0) {
            this.bodyMode = 'pretty'
            this.resMode = 0
            return
        }

        this.onPageSizeChange(this.pageSize)
    }

    onPageChange(pi: number) {
        if (pi > this.pageCount) {
            return
        }

        this.pageData = []

        let i = (pi - 1) * this.pageSize

        for (let j = 0; j < this.pageSize; j++) {
            let k = i + j
            if (k >= this.resList.length) {
                break
            }

            this.pageData.push(this.resList[k])
        }

        this.pageData[0]['_highlight'] = true
        this.onResChange(this.pageData[0], this.pageData[0])
    }

    onPageSizeChange(ps: number) {
        this.pageSize = ps
        this.pageCount = Math.ceil(this.resList.length / this.pageSize)
        this.onPageChange(1)
    }

    get status(): string {
        if (!this.res) {
            return ''
        }

        let status = `${this.$t(
            'req.pre.proto'
        )}：<span style="color:#3399ff;">${this.res.proto}</span>`
        status += `&nbsp;&nbsp;&nbsp;&nbsp;${this.$t(
            'req.pre.status'
        )}：<span style="color:#3399ff;">${this.res.status}</span>`
        status += `&nbsp;&nbsp;&nbsp;&nbsp;${this.$t(
            'req.report.table.finish'
        )}：<span style="color:#3399ff;">${this.res.duration.finish}</span>`

        return status
    }

    get bodyRaw(): string {
        if (!this.res || !this.res.body) {
            return ''
        }

        return this.res.body
    }

    get body(): string {
        if (!this.res || !this.res.body) {
            return ''
        }

        let body = this.res.body
        let contentType = 'text/plain'

        _.forEach(this.res.headers, (v, k) => {
            k = k.toLowerCase()
            if (k == 'content-type') {
                contentType = v[0]
            }
            return
        })

        let language: string[] | undefined = undefined
        if (contentType.indexOf('application/json') != -1) {
            body = JSON.stringify(JSON.parse(body), undefined, 4)
            language = ['json']
        } else if (contentType.indexOf('application/xml') != -1) {
            language = ['xml']
        } else if (contentType.indexOf('text/html') != -1) {
            language = ['html']
        }

        let pretty = hljs.highlightAuto(body, language)
        /* this.bodyLanguage = pretty.language */
        if (!language) {
            this.bodyLanguage += '(auto)'
        }
        body =
            '<ul class="hljs-ul"><li>' +
            _.replace(pretty.value, /\n/g, '\n</li><li>') +
            '</li></ul>'

        return body
    }

    get headers(): Req.HeaderModel[] {
        if (!this.res) {
            return []
        }

        let headers: Req.HeaderModel[] = []

        _.forIn(this.res.headers, (va, k) => {
            // headers += k + '：' + _.join(v, ';') + '\r\n'
            _.forEach(va, v => {
                headers.push({
                    key: k,
                    value: v
                })
            })
        })

        return headers
    }

    get cookies(): Req.CookieModel[] {
        return this.res ? this.res.cookies : []

        // if (!this.res) {
        //     return ''
        // }

        // let cookies: string = ''

        // _.forEach(this.res.Cookies, (v) => {
        //     cookies += JSON.stringify(v) + '\r\n'
        // })

        // return cookies
    }

    get durations(): Req.ResponseDuration[] {
        let durations: Req.ResponseDuration[] = []

        if (this.res) {
            durations.push(this.res.duration)
        }

        return durations
    }

    get code(): string {
        return this.res ? this.res.code + '[body ignore]' : ''
    }
}
</script>
