<template>
  <div>
    <Card style="width:100%;">
      <Collapse v-model="defaultSelect" simple>
        <Panel name="1">
          {{$t('title.respResultTitle')}}
          <h2 align="center" slot="content">{{ requestData.testResult }}</h2>
          <p slot="content">{{ requestData.failReason }}</p>
        </Panel>
        <Panel name="2">
          {{$t('title.outputTitle')}}
          <pre slot="content" v-html="requestData.output"></pre>
        </Panel>
        <Panel name="2">
          {{$t('title.reqDataTitle')}}
          <p slot="content">Url: {{ requestData.url }}</p>
          <p slot="content">Header: </p>
          <pre slot="content" v-html="requestData.header"></pre>
          <p slot="content">Request: </p>
          <pre slot="content" v-html="requestData.request"></pre>
        </Panel>
      </Collapse>
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
export default class Request extends Vue {
  @Prop() isAdvanced: boolean
  @Prop() resList: Req.ResponseModel[]
  @Prop() requestData: Req.ReqDataRespModel

  defaultSelect: string = "1"
  res: Req.ResponseModel | null = null
  bodyMode: 'pretty' | 'raw' = 'pretty'

  created() {
    this.onLocale()
  }

  @Watch('$i18n.locale')
  onLocale() {
  }

  mounted() {
    // this.getPrettyResp()
  }

  onResChange(res: Req.ResponseModel, oldRes: Req.ResponseModel) {
    this.res = res

  }

  @Watch('getPrettyBody')
  onResListChange() {
    this.bodyMode = 'pretty'
    return
  }

  get getBody(): string {
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
    } else {
      body = JSON.stringify(JSON.parse(body), undefined, 4)
      language = ['json']
    }

    let pretty = hljs.highlightAuto(body, language)

    body =
        '<ul class="hljs-ul"><li>' +
        _.replace(pretty.value, /\n/g, '\n</li><li>') +
        '</li></ul>'
    console.log("body")
    console.log(body)
    return body
  }
}
</script>
