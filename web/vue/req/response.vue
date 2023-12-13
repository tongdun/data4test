<template>
  <div>
    <Card style="width:100%;">
      <Collapse v-model="defaultSelect" simple>
        <Panel name="1">
          {{$t('title.respDataTitle')}}
          <pre slot="content" v-html="requestData.response"></pre>
        </Panel>
      </Collapse>
<!--      <div slot="title">{{$t('title.respDataTitle')}}</div>-->
<!--      <pre v-html="requestData.response"></pre>-->
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
export default class Response extends Vue {
  @Prop() isAdvanced: boolean
  @Prop() resList: Req.ResponseModel[]
  @Prop() requestData: Req.ReqDataRespModel

  res: Req.ResponseModel | null = null
  bodyMode: 'pretty' | 'raw' = 'pretty'
  defaultSelect: string = "1"
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
}
</script>
