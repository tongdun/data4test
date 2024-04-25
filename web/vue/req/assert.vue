<template>
  <div>
    <Table :columns="assertTableColumns" :data="assertTableData" :no-data-text="$t('general.noData')" :no-filtered-data-text="$t('general.noFilterData')"></Table>
  </div>
</template>

<script lang="ts">
import { Vue, Component, Prop, Watch } from 'vue-property-decorator'
import { CreateElement } from 'vue'
import _ from 'lodash'
import API from "../../ts/api";
// import moment from 'moment'

@Component
export default class AssertList extends Vue {
  @Prop() isAsserts: boolean
  @Prop() isSending: boolean = false
  @Prop() contentType: string
  @Prop() assertListData: Req.AssertListModel[]

  assertTableData: any[] = []
  assertType: string[] = ["re", "output", "output_re","=", "!=", ">", ">=", "<", "<=", "in", "!in", "sum", "avg", "count", "not_in", "equal", "not_equal", "contain", "not_contain", "null", "not_null", "regexp"]
  typeOptions: string[] = []
  sourceOptions: string[] = []
  valueOptions: string[] = []
  assertResultOptions: string[] = []

  assertTableColumns: any[] = [
    {
      title: '',
      key: 'type',
      width: 140,
      render: (h: CreateElement, params: any) => {
        return h('Select', {
              props: {
                value: params.row.type,
                transfer: true,
                filterable: true,
                data: this.getOptions('type')
                // 'filter-method': this.onFilterMethod
              },
              on: {
                'on-change': (value: string) => {
                  this.onChange(params.index, 'type', value)
                }
              }
            },
            this.assertType.map((item) => {
              return h('Option', {
                props: {value: item}
              }, item)
            })
        )
      }
    },
    {
      title: '',
      // width: 450,
      key: 'source',
      render: (h: CreateElement, params: any) => {
        return h('AutoComplete', {
          props: {
            value: params.row.source,
            transfer: true,
            data: this.getOptions('source')
            // 'filter-method': this.onFilterMethod
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
      key: 'value',
      width: 250,
      render: (h: CreateElement, params: any) => {
        return h('AutoComplete', {
              props: {
                value: params.row.value,
                transfer: true,
                filterable: true,
                'filter-method': this.onFilterMethod,
                // data: this.getOptions("value"),
                'allow-create': true,
                'on-create': this.onAddAssertTemplateValue(params.row.value)
              },
              on: {
                'on-change': (value: string) => {
                  this.onChange(params.index, 'value', value)
                }
              }
            },
            this.valueOptions.map((item) => {
              return h('Option', {
                props: {label: item, value: item}
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



  created() {
    this.onLocale()
  }

  @Watch('$i18n.locale')
  onLocale() {
    //注意 tableColumns 的处理位置需要与上面初始化定义一致
    this.assertTableColumns[0].title = this.$t('assertList.type')
    this.assertTableColumns[1].title = this.$t('assertList.source')
    this.assertTableColumns[2].title = this.$t('assertList.value')
    this.assertTableColumns[3].title = this.$t('assertList.assertResult')

  }

  mounted() {
    this.syncTableData()
    this.getAssertTemplateList()
  }

  @Watch('assertListData')
  onListDataChange() {
    this.syncTableData()
  }

  onAddAssertTemplateValue(val) {
    if (this.valueOptions.indexOf(val)<0) {
      this.valueOptions.push(val)
    }
  }

  onFilterMethod(value: string, option: string) {
    return option.toLowerCase().indexOf(value.toLowerCase()) != -1
  }

  onAdd() {
    this.assertListData.push({
      isDisable: false,
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

  onSelect(selection: Req.AssertListModel[]) {
    _.forEach(this.assertListData, p => {
      p.isDisable = true

      let s = _.find(selection, function(s) {
        return p.value === s.value
      })
      if (s) {
        p.isDisable = false
      }
    })
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

  async getAssertTemplateList() {
    let result = await API.get<Req.ResponseModel[]>('/assertTemplateList')
    if (result.data) {
      this.valueOptions = result.data.map(item => (item['name']));
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

