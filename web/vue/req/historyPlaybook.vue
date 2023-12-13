<template>
  <div>
    <Table border :columns="historyApiTableColumns" :data="historyApiTableData" :no-data-text="$t('general.noData')" :no-filtered-data-text="$t('general.noFilterData')"></Table>
  </div>
</template>

<script lang="ts">
import { Vue, Component, Prop, Watch } from 'vue-property-decorator'
import { CreateElement } from 'vue'

import API from "../../ts/api";
import _ from "lodash";

@Component

export default class historyApiList extends Vue {
  @Prop() isRelatedApi: boolean
  @Prop() isSending: boolean = false
  @Prop() historyApiListData: Req.HistoryApiListModel[]

  historyApiTableData: any[] = []
  historyApiTableColumns: any[] = [
    {
      title: '关联数据',
      key: 'historyDataFile',
      width: 800,
      render: (h: CreateElement, params: any) => {
        return h('AutoComplete', {
              props: {
                value: params.row.historyDataFile,
                transfer: true,
                data: this.getOptions['historyDataFile'],
                'filter-method': this.onFilterMethod
              },
              on: {
                'on-change': (value: string) => {
                  this.onChange(params.index, 'historyDataFile', value)
                }
              }
            }
        )
      }
    },
    {
      width: 80,
      renderHeader: (h: CreateElement, params: any) => {
        return h('Button', {
          props: {
            disabled: this.isSending,
            size: 'small',
            type: 'primary',
            icon: 'plus-round'
          },
          on: {
            click: () => {
              this.onAdd()
            }
          }
        })
      },
      render: (h: CreateElement, params: any) => {
        return h('Button', {
          props: {
            disabled: this.isSending,
            size: 'small',
            icon: 'close-round'
          },
          on: {
            click: () => {
              this.onDel(params.index)
            }
          }
        })
      }
    }
  ]

   historyDataFileOptions: string[] = []

  mounted() {
    this.syncTableData();
  }

  created() {
    this.onLocale();
  }

  @Watch('$i18n.locale')
  onLocale() {
    //注意 tableColumns 的处理位置需要与上面初始化定义一致
    this.historyApiTableColumns[0].title = this.$t('relatedApiList.dataFile')
  }

  @Watch('historyApiListData')
  onListDataChange() {
    this.syncTableData()
  }

  onFilterMethod(value: string, option: string) {
    return option.toLowerCase().indexOf(value.toLowerCase()) != -1
  }

  onAdd() {
    this.historyApiListData.push({
      dataFile: ''
    })
    this.$emit('onChange')
    this.syncTableData()
  }

  onDel(index: number) {
    //删除指定的成员
    this.historyApiListData.splice(index, 1)
    this.$emit('onChange')

    this.syncTableData()
  }

  onUp(index: number) {
    //上移指定的成员
    this.historyApiListData.sort()
    this.$emit('onChange')

    this.syncTableData()
  }

  onDown(index: number) {
    //上移指定的成员
    this.historyApiListData.sort()
    this.$emit('onChange')

    this.syncTableData()
  }

  onChange(index: number, name: string, value: string) {
    if (name === 'historyDataFile') {
      this.historyApiListData[index].dataFile = value
    }
    this.$emit('onChange')
  }

  getOptions(type: 'historyDataFile') {
    return  this.historyDataFileOptions
  }

  syncTableData() {
    //全量复制
    this.historyApiTableData = this.historyApiListData.map((m): any => {
      return {
        historyDataFile: m.dataFile,
      }
    })
  }

}
</script>

