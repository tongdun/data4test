<template>
  <div>
    <Table :columns="defTableColumns" :data="defTableData" :no-data-text="$t('general.noData')" :no-filtered-data-text="$t('general.noFilterData')"></Table>
  </div>
</template>

<script lang="ts">
import { Vue, Component, Prop, Watch } from 'vue-property-decorator'
import { CreateElement } from 'vue'
import _ from 'lodash'
import moment from 'moment'

@Component
export default class DefList extends Vue {
  @Prop() isHeaders: boolean
  @Prop() isSending: boolean = false
  @Prop() contentType: string
  @Prop() defListData: Req.DefListModel[]

  defTableData: any[] = []

  defTableColumns: any[] = [
    {
      title: '',
      key: 'name',
      width: 150,
      render: (h: CreateElement, params: any) => {
        return h('AutoComplete', {
          props: {
            value: params.row.name,
            transfer: true,
            data: this.getOptions('name'),
            'filter-method': this.onFilterMethod
          },
          on: {
            'on-change': (value: string) => {
              this.onChange(params.index, 'name', value)
            }
          }
        })
      }
    },
    {
      title: '',
      width: 120,
      key: 'valueType',
      render: (h: CreateElement, params: any) => {
        return h('Select', {
              props: {
                value: params.row.valueType,
                transfer: true,
                data: this.getOptions('valueType'),
                'filter-method': this.onFilterMethod
              },
              on: {
                'on-change': (value: string) => {
                  this.onChange(params.index, 'valueType', value)
                }
              }
            },
            [
              h('Option', {
                props: {
                  value: 'string'
                }
              }, 'string'),
              h('Option', {
                props : {
                  value: 'bool'
                }
              }, 'bool'),
              h('Option', {
                props : {
                  value: 'integer'
                }
              }, 'integer'),
              h('Option', {
                props : {
                  value: 'object'
                }
              }, 'object'),
              h('Option', {
                props : {
                  value: 'array'
                }
              }, 'array'),
              h('Option', {
                props : {
                  value: 'float64'
                }
              }, 'float32'),
              h('Option', {
                props : {
                  value: 'float64'
                }
              }, 'float32')
            ])
      }
    },
    {
      title: '',
      key: 'isMust',
      width: 100,
      render: (h: CreateElement, params: any) => {
        return h('Select', {
              props: {
                value: params.row.isMust,
                transfer: true,
                data: this.getOptions("isMust"),
                'filter-method': this.onFilterMethod
              },
              on: {
                'on-change': (value: string) => {
                  this.onChange(params.index, 'isMust', value)
                }
              },
            },
            [
              h('Option', {
                props: {
                  value: 'yes'
                }
              }, 'yes'),
              h('Option', {
                props : {
                  value: 'no'
                }
              }, 'no')
            ])
      }
    },
    {
      title: '',
      key: 'egValue',
      render: (h: CreateElement, params: any) => {
        return h('AutoComplete', {
          props: {
            value: params.row.egValue,
            transfer: true,
            data: this.getOptions("egValue"),
            'filter-method': this.onFilterMethod
          },
          on: {
            'on-change': (value: string) => {
              this.onChange(params.index, 'egValue', value)
            }
          }
        })
      }
    },
    {
      title: '',
      key: 'desc',
      render: (h: CreateElement, params: any) => {
        return h('AutoComplete', {
          props: {
            value: params.row.desc,
            transfer: true,
            data: this.getOptions("desc"),
            'filter-method': this.onFilterMethod
          },
          on: {
            'on-change': (value: string) => {
              this.onChange(params.index, 'desc', value)
            }
          }
        })
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
  bodyTypeOptions: string[] = []
  bodyValueTypeOptions: string[] = []
  nameOptions: string[] = []
  valueTypeOptions: string[] = []
  isMustOptions: string[] = []
  egValueOptions: string[] = []
  descOptions: string[] = []

  created() {
    this.onLocale()
  }

  @Watch('$i18n.locale')
  onLocale() {
    //注意 tableColumns 的处理位置需要与上面初始化定义一致
    this.defTableColumns[0].title = this.$t('list.name')
    this.defTableColumns[1].title = this.$t('list.valueType')
    this.defTableColumns[2].title = this.$t('list.isMust')
    this.defTableColumns[3].title = this.$t('list.egValue')
    this.defTableColumns[4].title = this.$t('list.desc')
  }

  mounted() {
    this.syncDefTableData()
    // this.onDefContentTypeChange()
    this.onDefListDataChange()

    if (this.isHeaders) {
      this.bodyTypeOptions = ['Content-Type', 'User-Agent']
      this.bodyValueTypeOptions = [
        'text/plain',
        'text/html',
        'application/json',
        'application/javascript',
        'application/xml',
        'application/x-www-form-urlencoded'
      ]
    }
  }

  @Watch('defListData')
  onDefListDataChange() {
    this.syncDefTableData()
  }

  @Watch('contentType')
  onDefContentTypeChange() {
    let isFound = false
    if (!isFound) {
      this.defListData.push({
        name: '',
        valueType: 'string',
        isMust: 'no',
        egValue: '',
        desc: '',
      })
    }

    this.syncDefTableData()
  }

  onFilterMethod(value: string, option: string) {
    return option.toLowerCase().indexOf(value.toLowerCase()) != -1
  }

  onAdd() {
    this.defListData.push({
      name: '',
      valueType: 'string',
      isMust: 'no',
      egValue: '',
      desc: '',
    })
    this.$emit('onChange')
    this.syncDefTableData()
  }

  onDel(index: number) {
    //删除指定的成员
    this.defListData.splice(index, 1)
    this.$emit('onChange')

    this.syncDefTableData()
  }

  onChange(index: number, name: string, value: string) {
    if (name === 'name') {
      this.defListData[index].name = value
    } else if (name === 'valueType') {
      this.defListData[index].valueType = value
    }  else if (name === 'isMust') {
      this.defListData[index].isMust = value
    }  else if (name === 'egValue') {
      this.defListData[index].egValue = value
    } else if (name === 'desc') {
      this.defListData[index].desc = value
    }
    this.$emit('onChange')
  }

  getOptions(type: string) {
    if (type === 'name') {
      return this.nameOptions
    } else if (type === 'valueType') {
      return this.valueTypeOptions
    } else if (type === 'isMust') {
      return this.isMustOptions
    } else if (type === 'egValue') {
      return this.egValueOptions
    } else if (type === 'desc') {
      return this.descOptions
    }
  }

  syncDefTableData() {
    //全量复制
    this.defTableData = this.defListData.map((n): any => {
      return {
        name: n.name,
        valueType: n.valueType,
        isMust: n.isMust,
        egValue: n.egValue,
        desc: n.desc,
      }
    })
  }
}
</script>

