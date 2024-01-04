<template>
  <div>
    <Table :columns="relatedApiTableColumns" :data="relatedApiTableData" :no-data-text="$t('general.noData')" :no-filtered-data-text="$t('general.noFilterData')"></Table>
  </div>
</template>

<script lang="ts">
import { Vue, Component, Prop, Watch } from 'vue-property-decorator'
import { CreateElement } from 'vue'

import API from "../../ts/api";
import _ from "lodash";

@Component

export default class depApiList extends Vue {
  @Prop() isRelatedApi: boolean
  @Prop() isSending: boolean = false
  @Prop() relatedApiListData: Req.RelatedApiListModel[]
  @Prop() allDataFile: string[]

  relatedApiTableData: any[] = []
  relatedApiTableColumns: any[] = [
    {
      title: '关联数据',
      key: 'dataFile',
      render: (h: CreateElement, params: any) => {
        return h('Select', {
              props: {
                value: params.row.dataFile,
                transfer: true,
                filterable: true,
                clearable: true,
                data: this.getOptions['dataFile'],
                'filter-method': this.onFilterMethod,
                // 'on-change': this.onAddInValueType(params.row.dataFile),
              },
              on: {
                'on-change': (value: string) => {
                  this.onChange(params.index, 'dataFile', value)
                }
              }
            },
            this.allDataFile.map((item) => {
              return h('Option', {
                props: {value: item}
              }, item)
            })
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

  dataFileOptions: string[] = []

  mounted() {
    this.syncTableData();
  }

  created() {
    this.onLocale();
  }

  @Watch('$i18n.locale')
  onLocale() {
    //注意 tableColumns 的处理位置需要与上面初始化定义一致
    this.relatedApiTableColumns[0].title = this.$t('relatedApiList.dataFile')
  }

  @Watch('relatedApiListData')
  onListDataChange() {
    this.syncTableData()
  }

  onFilterMethod(value: string, option: string) {
    return option.toLowerCase().indexOf(value.toLowerCase()) != -1
  }

  onAdd() {
    this.relatedApiListData.push({
      dataFile: ''
    })
    this.$emit('onChange')
    this.syncTableData()
  }

  onDel(index: number) {
    //删除指定的成员
    this.relatedApiListData.splice(index, 1)
    this.$emit('onChange')

    this.syncTableData()
  }

  onUp(index: number) {
    //上移指定的成员
    this.relatedApiListData.sort()
    this.$emit('onChange')

    this.syncTableData()
  }

  onDown(index: number) {
    //上移指定的成员
    this.relatedApiListData.sort()
    this.$emit('onChange')

    this.syncTableData()
  }

  onChange(index: number, name: string, value: string) {
    if (name === 'dataFile') {
      this.relatedApiListData[index].dataFile = value
    }
    this.$emit('onChange')
  }

  getOptions(type: 'dataFile') {
    return  this.dataFileOptions
  }

  syncTableData() {
    //全量复制
    this.relatedApiTableData = this.relatedApiListData.map((m): any => {
      return {
        dataFile: m.dataFile,
      }
    })
  }

}
</script>

