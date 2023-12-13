<template>
    <div>
        <Card style="width:100%;">
            <div slot="title">{{$t('req.report.title')}}</div>
            <div>
                <Row style="height:240px;">
                    <Col span="8">
                    <div>{{$t('req.report.summary.total')}}：{{total}}&nbsp;s</div>
                    <div>{{$t('req.report.summary.slowest')}}：{{slowestFinish}}&nbsp;ms</div>
                    <div>{{$t('req.report.summary.fastest')}}：{{fastestFinish}}&nbsp;ms</div>
                    <div>{{$t('req.report.summary.average')}}：{{avgFinish}}&nbsp;ms</div>
                    <div>{{$t('req.report.summary.rps')}}：{{rps}}</div>
                    </Col>
                    <Col span="16">
                    <div>
                        <Table stripe border size="small" :columns="durationColumns" :data="durations" :no-data-text="$t('general.noData')" :no-filtered-data-text="$t('general.noFilterData')"></Table>
                    </div>
                    </Col>
                    <!-- <Col span="10">
                        <div :id="chartID" style="height:240px;"></div>
                    </Col> -->
                </Row>
            </div>
        </Card>
    </div>
</template>

<script lang="ts">
import { Vue, Component, Prop, Watch } from 'vue-property-decorator'
import echarts from 'echarts'
import moment from 'moment'

@Component
export default class report extends Vue {
    @Prop() start: number
    @Prop() resList: Req.ResponseModel[]
    @Prop() tabID: string

    durationColumns = <any[]>[
        {
            title: '',
            key: 'type'
        },
        {
            title: '',
            key: 'average'
        },
        {
            title: '',
            key: 'slowest'
        },
        {
            title: '',
            key: 'fastest'
        }
    ]
    durations: any[] = []

    ec: echarts.ECharts
    ecOption: echarts.EChartOption

    total = 0
    rps = 0
    avgDNS = 0
    avgConn = 0
    avgReq = 0
    avgDelay = 0
    avgRes = 0
    avgFinish = 0

    fastestDNS = 0
    fastestConn = 0
    fastestReq = 0
    fastestDelay = 0
    fastestRes = 0
    fastestFinish = 0

    slowestDNS = 0
    slowestConn = 0
    slowestReq = 0
    slowestDelay = 0
    slowestRes = 0
    slowestFinish = 0

    created() {
        this.onLocale()
    }

    @Watch('$i18n.locale')
    onLocale() {
        this.durationColumns[0].title = this.$t(
            'req.report.table.type'
        ).toString()
        this.durationColumns[1].title =
            this.$t('req.report.table.average').toString() + '(ms)'
        this.durationColumns[2].title =
            this.$t('req.report.table.slowest').toString() + '(ms)'
        this.durationColumns[3].title =
            this.$t('req.report.table.fastest').toString() + '(ms)'

        this.onDurationShow()
        // this.onCharts();
    }

    mounted() {
        // this.ec = echarts.init(<HTMLDivElement>document.getElementById(
        //   this.chartID
        // ));
    }

    get chartID() {
        return 'mychart' + this.tabID
    }

    @Watch('resList')
    onResListChange() {
        this.onCount()

        this.onDurationShow()
        // this.onCharts();
    }

    onDurationShow() {
        this.durations = [
            {
                type: this.$t('req.report.table.dns'),
                fastest: this.fastestDNS,
                slowest: this.slowestDNS,
                average: this.avgDNS
            },
            {
                type: this.$t('req.report.table.conn'),
                fastest: this.fastestConn,
                slowest: this.slowestConn,
                average: this.avgConn
            },
            {
                type: this.$t('req.report.table.req'),
                fastest: this.fastestReq,
                slowest: this.slowestReq,
                average: this.avgReq
            },
            {
                type: this.$t('req.report.table.delay'),
                fastest: this.fastestDelay,
                slowest: this.slowestDelay,
                average: this.avgDelay
            },
            {
                type: this.$t('req.report.table.res'),
                fastest: this.fastestRes,
                slowest: this.slowestRes,
                average: this.avgRes
            }
        ]
    }

    onCount() {
        if (this.resList.length == 0) {
            return
        }

        //总耗时，单位：秒
        this.total = this.roundNumber((moment().valueOf() - this.start) / 1000)

        // req/sec 每秒请求籹
        this.rps = this.roundNumber(this.resList.length / this.total)

        let resCount = 0

        this.avgDNS = 0
        this.avgConn = 0
        this.avgReq = 0
        this.avgDelay = 0
        this.avgRes = 0
        this.avgFinish = 0

        this.fastestDNS = 0
        this.fastestConn = 0
        this.fastestReq = 0
        this.fastestDelay = 0
        this.fastestRes = 0
        this.fastestFinish = 0

        this.slowestDNS = 0
        this.slowestConn = 0
        this.slowestReq = 0
        this.slowestDelay = 0
        this.slowestRes = 0
        this.slowestFinish = 0

        for (let res of this.resList) {
            if (!res.error) {
                resCount++
            }

            let dnsDuration = this.formatDuration(res.duration.dns),
                connDuration = this.formatDuration(res.duration.conn),
                reqDuration = this.formatDuration(res.duration.req),
                delayDuration = this.formatDuration(res.duration.delay),
                resDuration = this.formatDuration(res.duration.res),
                finishDuration = this.formatDuration(res.duration.finish)

            if (dnsDuration < this.fastestDNS) {
                this.fastestDNS = dnsDuration
            }
            if (dnsDuration > this.slowestDNS) {
                this.slowestDNS = dnsDuration
            }

            if (connDuration < this.fastestConn) {
                this.fastestConn = connDuration
            }
            if (connDuration > this.slowestConn) {
                this.slowestConn = connDuration
            }

            if (reqDuration < this.fastestReq) {
                this.fastestReq = reqDuration
            }
            if (reqDuration > this.slowestReq) {
                this.slowestReq = reqDuration
            }

            if (delayDuration < this.fastestDelay) {
                this.fastestDelay = delayDuration
            }
            if (delayDuration > this.slowestReq) {
                this.slowestReq = delayDuration
            }

            if (resDuration < this.fastestRes) {
                this.fastestRes = resDuration
            }
            if (resDuration > this.slowestRes) {
                this.slowestRes = resDuration
            }

            if (finishDuration < this.fastestFinish) {
                this.fastestFinish = finishDuration
            }
            if (finishDuration > this.slowestFinish) {
                this.slowestFinish = finishDuration
            }

            this.avgDNS += dnsDuration
            this.avgConn += connDuration
            this.avgReq += reqDuration
            this.avgDelay += delayDuration
            this.avgRes += resDuration
            this.avgFinish += finishDuration
        }

        this.avgDNS = this.roundNumber(this.avgDNS / resCount)
        this.avgConn = this.roundNumber(this.avgConn / resCount)
        this.avgReq = this.roundNumber(this.avgReq / resCount)
        this.avgDelay = this.roundNumber(this.avgDelay / resCount)
        this.avgRes = this.roundNumber(this.avgRes / resCount)
        this.avgFinish = this.roundNumber(this.avgFinish / resCount)
    }

    roundNumber(n: number, f?: number): number {
        if (!f || f < 0) {
            f = 2
        }

        f = f * 10
        return Math.round(n * f) / f
    }

    formatDuration(str: string): number {
        let x = 1

        if (str.indexOf('ns') != -1) {
            str = str.replace('ns', '')
            x = 1000000
        } else if (str.indexOf('µs') != -1) {
            str = str.replace('µs', '')
            x = 1000
        } else if (str.indexOf('ms') != -1) {
            str = str.replace('ms', '')
        } else if (str.indexOf('s') != -1) {
            str = str.replace('s', '')
            x = 0.001
        }

        return parseFloat(str) / x
    }

    onCharts() {
        if (this.resList.length == 0) {
            return
        }

        this.ecOption = {
            title: {
                text:
                    this.$t('req.report.chart.title').toString() + '(ms)'
            },
            tooltip: {
                trigger: 'item',
                formatter: '{a} <br/>{b} : {c} ({d}%)'
            },
            legend: {
                orient: 'vertical',
                left: 'left',
                data: [
                    this.$t('req.report.table.dns').toString(),
                    this.$t('req.report.table.conn').toString(),
                    this.$t('req.report.table.req').toString(),
                    this.$t('req.report.table.delay').toString(),
                    this.$t('req.report.table.res').toString()
                ]
            },
            series: [
                {
                    name: this.$t(
                        'req.report.chart.serieName'
                    ).toString(),
                    type: 'pie',
                    radius: '55%',
                    center: ['50%', '60%'],
                    selectedMode: 'single',
                    data: [
                        {
                            value: this.avgDNS,
                            name: this.$t(
                                'req.report.table.dns'
                            ).toString()
                        },
                        {
                            value: this.avgConn,
                            name: this.$t(
                                'req.report.table.conn'
                            ).toString()
                        },
                        {
                            value: this.avgReq,
                            name: this.$t(
                                'req.report.table.req'
                            ).toString()
                        },
                        {
                            value: this.avgDelay,
                            name: this.$t(
                                'req.report.table.delay'
                            ).toString()
                        },
                        {
                            value: this.avgRes,
                            name: this.$t(
                                'req.report.table.res'
                            ).toString()
                        }
                    ],
                    itemStyle: {
                        emphasis: {
                            shadowBlur: 10,
                            shadowOffsetX: 0,
                            shadowColor: 'rgba(0, 0, 0, 0.5)'
                        }
                    }
                }
            ]
        }

        let opt = <any>this.ecOption
        opt.title.x = 'center'
        this.ec.setOption(opt)
        this.ec.resize()
    }

    getDurations(key: string): number[] {
        let data: number[] = []

        for (var res of this.resList) {
            data.push(parseInt(res.duration[key]))
        }

        return data
    }
}
</script>
