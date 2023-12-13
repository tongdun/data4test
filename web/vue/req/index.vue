<template>
    <div>
        <Tabs v-model="tabID" type="card" closable @on-tab-remove="onTabDel" :animated="false" style="height:100%;" name="first">
            <Tab-pane v-for="tab in tabs" v-show="tab.isShow" :name="tab.id" :key="tab.id" :label="tab.label" tab="first">
                <ReqTab :tab="tab"></ReqTab>
            </Tab-pane>
            <ButtonGroup slot="extra">
                <Button type="primary" icon="plus-round" ghost @click="onTabAdd" size="large" class="layout-button-margin"></Button>
            </ButtonGroup>
        </Tabs>
        <BackTop></BackTop>
    </div>
</template>

<style>
.layout-button-margin {
  margin-top: 0px;
  height: 28px;
}
</style>

<script lang="ts">
import { Component } from 'vue-property-decorator'
import _ from 'lodash'
import moment from 'moment'
import ReqTab from './tab.vue'
import Vue from 'vue'

@Component({
    components: {
        ReqTab
    }
})
export default class Index extends Vue {
    tabs: Req.TabModel[] = []
    tabID: string = ''

    mounted() {
        this.onTabAdd()
    }

    /**
     * 获取指定的tab对象
     */
    getTab(id: string): Req.TabModel | undefined {
        return _.find(this.tabs, function(v) {
            return v.id === id
        })
    }

     handleTabsAdd () {
            this.tabs = this.tabs.concat()
  }
    /**
     * 添加tab
     */
    onTabAdd() {
        let id =
            'tab' +
            moment()
                .valueOf()
                .toString()
        let tab: Req.TabModel = {
            id: id,
            isShow: true,
            label: ''
        }
        this.tabs.push(tab)
        this.tabID = id
    }

    /**
     * 移除tab
     */
    onTabDel(id: string) {
        _.remove(this.tabs, tab => {
            let has = tab.id === id
            if (has && this.tabs.length == 1) {
                this.onTabAdd()
            }

            return has
        })
    }
}
</script>

