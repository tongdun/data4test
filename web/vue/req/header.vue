<template>
  <div>
    <Table :columns="headerTableColumns" :data="headerTableData" :no-data-text="$t('general.noData')" :no-filtered-data-text="$t('general.noFilterData')" border></Table>
  </div>
</template>

<script lang="ts">
import { Vue, Component, Prop, Watch } from 'vue-property-decorator'
import { CreateElement } from 'vue'
import _ from 'lodash'
// import moment from 'moment'

@Component
export default class HeaderList extends Vue {
  @Prop() isHeader: boolean
  @Prop() isSending: boolean = false
  @Prop() contentType: string
  @Prop() headerListData: Req.HeaderListModel[]

  headerTableData: any[] = []
  headerType: string[] = ["Content-Type", "Cookie", "Referer", "X-Cf-Random"]

  headerTableColumns: any[] = [
    {
      title: '',
      key: 'name',
      width: 150,
      render: (h: CreateElement, params: any) => {
        return h('Select', {
              props: {
                value: params.row.type,
                transfer: true,
                filterable: true,
                data: this.getOptions('type'),
                'filter-method': this.onFilterMethod
              },
              on: {
                'on-change': (value: string) => {
                  this.onChange(params.index, 'name', value)
                }
              }
            },
            this.headerType.map((item) => {
              return h('Option', {
                props: {value: item}
              }, item)
            })
        )
      }
    },
    {
      title: '',
      width: 250,
      // key: 'source',
      render: (h: CreateElement, params: any) => {
        return h('AutoComplete', {
          props: {
            value: params.row.source,
            transfer: true,
            data: this.getOptions('source'),
            'filter-method': this.onFilterMethod
          },
          on: {
            'on-change': (value: string) => {
              this.onChange(params.index, 'source', value)
            }
          }
        })
      }
    },
    {
      title: '',
      // key: 'value',
      width: 250,
      render: (h: CreateElement, params: any) => {
        return h('AutoComplete', {
              props: {
                value: params.row.value,
                transfer: true,
                data: this.getOptions("value"),
                'filter-method': this.onFilterMethod
              },
              on: {
                'on-change': (value: string) => {
                  this.onChange(params.index, 'value', value)
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

  typeOptions: string[] = []
  sourceOptions: string[] = []
  valueOptions: string[] = []
  assertResultOptions: string[] = []

  created() {
    this.onLocale()
  }

  @Watch('$i18n.locale')
  onLocale() {
    //注意 tableColumns 的处理位置需要与上面初始化定义一致
    this.headerTableColumns[0].title = this.$t('headerList.isDisable')
    this.headerTableColumns[1].title = this.$t('headerList.name')
    this.headerTableColumns[2].title = this.$t('headerList.valueType')
    this.headerTableColumns[3].title = this.$t('headerList.assertResult')

  }

  mounted() {
    this.syncTableData()
    // this.onRunListDataChange()
  }

  @Watch('assertListData')
  onListDataChange() {
    this.syncTableData()
  }

  onFilterMethod(value: string, option: string) {
    return option.toLowerCase().indexOf(value.toLowerCase()) != -1
  }

  onAdd() {
    this.assertListData.push({
      type: '',
      source: '',
      value: '',
      assertResult: ''
    })
    this.$emit('onChange')
    this.syncTableData()
  }

  onDel(index: number) {
    //删除指定的成员
    this.assertListData.splice(index, 1)
    this.$emit('onChange')

    this.syncTableData()
  }

  onChange(index: number, name: string, value: string) {
    if (name === 'type') {
      this.assertListData[index].type = value
    } else if (name === 'source') {
      this.assertListData[index].source = value
    }  else if (name === 'value') {
      this.assertListData[index].value = value
    } else if (name === 'assertResult') {
      this.assertListData[index].assertResult = value
    }
    this.$emit('onChange')
  }

  getOptions(type: 'type' | 'source' | 'value' | 'assertResult') {
    if (type==='type') {
      return this.typeOptions
    } else if (type==='source') {
      return this.sourceOptions
    } else if (type==='value') {
      return this.valueOptions
    }  else if (type==='assertResult') {
      return this.assertResultOptions
    }

  }

  syncTableData() {
    //全量复制
    this.assertTableData = this.assertListData.map((m): any => {
      return {
        type: m.type,
        source: m.source,
        value: m.value,
        assertResult: m.assertResult
      }
    })
  }
}
</script>

