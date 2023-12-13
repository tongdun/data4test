<template>
  <div>
    <Card style="width:100%;">
      <Collapse v-model="defaultSelect" simple>
        <Panel name="1">
          {{$t('title.respResultTitle')}}
          <pre slot="content" v-html="requestScene.testResult"></pre>
          <pre slot="content" v-html="requestScene.failReason"></pre>
          <pre slot="content" v-html="requestScene.lastDataFile"></pre>
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
export default class SceneResponse extends Vue {
  @Prop() isAdvanced: boolean
  @Prop() resList: Req.ReqSceneRespModel[]
  @Prop() requestScene: Req.ReqSceneRespModel
  @Prop() requestSceneHistory: Req.ReqSceneRespModel

  defaultSelect: string = "1"
  created() {
    this.onLocale()
  }

  @Watch('$i18n.locale')
  onLocale() {
  }

  mounted() {
  }

}
</script>
